// Author mogd 2022-05-18
//
// Update mogd 2022-05-18
//
// Description Fetch different cloud vendor latencies by pinging
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	InputFile   string
	OutPath     string
	LogFile     string
	Regional    string
	CORE        int
	CHANNEL     int
	SugarLogger *zap.SugaredLogger
	Stats       = make(chan *ping.Statistics, 1000)
	ResRecord   = make(chan string, 1000)
)

func main() {
	flag.StringVar(&InputFile, "infile", "./test.csv", "input csv file")
	flag.StringVar(&OutPath, "outpath", "./", "output conf file")
	flag.StringVar(&Regional, "region", "Asia", "continents")
	flag.StringVar(&LogFile, "log", "./test.log", "log file")
	flag.IntVar(&CORE, "core", 4, "cpu cores")
	flag.IntVar(&CHANNEL, "channel", 100, "cpu cores")
	flag.Parse()
	start := time.Now()
	fmt.Println(time.Now())
	InitLogger()
	runtime.GOMAXPROCS(CORE)
	go readCsv()
	go GetPing()
	writeResultCSV()
	fmt.Println(time.Since(start))
	SugarLogger.Info("Success!")
}

// InitLogger init logger
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	SugarLogger = logger.Sugar()
}

// getEncoder encoder
//
// return zapcore.Encoder
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

// getLogWriter set log writer config
//
// return zapcore.WriteSyncer
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   LogFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// readCsv read csv
func readCsv() {
	inCsv, err := os.OpenFile(InputFile, os.O_RDONLY, 0777)
	if err != nil {
		SugarLogger.Fatal(fmt.Sprintf("Can not open file %v, err: %v\n", InputFile, err))
	}
	defer inCsv.Close()
	r := csv.NewReader(inCsv)

	for {
		record, err := r.Read()
		if err == io.EOF {
			SugarLogger.Error("io eof, err:", zap.Error(err))
			// Wg.Done()
			break
		}
		if err != nil {
			SugarLogger.Fatal(fmt.Sprintf("csv read fail: %v\n", err))
		}
		// record[1], _ = simplifiedchinese.GBK.NewDecoder().String(record[1])
		// record[2], _ = simplifiedchinese.GBK.NewDecoder().String(record[2])
		SugarLogger.Info("read: ", record[0])
		ResRecord <- record[0]
	}
	defer close(ResRecord)
}

// GetPing GetPing the IP or AvgRtt by pinging
func GetPing() {
	var wg sync.WaitGroup
	defer close(Stats)
	for {
		ip, ok := <-ResRecord
		if !ok {
			break
		}
		wg.Add(1)
		go func() {
			// SugarLogger.Info("stat ping IP: ", s)
			pinger, err := ping.NewPinger(ip)
			pinger.SetPrivileged(true)
			if err != nil {
				SugarLogger.Error("New Pinger error: ", zap.Error(err))
				return
			}
			pinger.Count = 100
			pinger.Timeout = 30 * time.Second
			// pinger.PacketsSent = 100
			err = pinger.Run() // Blocks until finished.
			if err != nil {
				SugarLogger.Error("Pinger run error: ", zap.Error(err))
				return
			}
			stats := pinger.Statistics()
			// SugarLogger.Info("end ping IP: ", s)
			Stats <- stats
			wg.Done()
		}()
	}
	wg.Wait()
}

// writeResultCSV write end of result to csv file
func writeResultCSV() {
	fileName := OutPath + Regional + ".csv"
	fw, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		SugarLogger.Fatal(fmt.Sprintf("open file %v error: %v", fileName, err))
	}
	defer fw.Close()

	var wg sync.WaitGroup
	for {
		stats, ok := <-Stats
		if !ok {
			// Wg.Done()
			break
		}
		wg.Add(1)
		go func() {
			write := bufio.NewWriter(fw)
			SugarLogger.Info(fmt.Sprintf("IP = %v\n", stats.IPAddr.String()))
			SugarLogger.Info(fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt))
			SugarLogger.Info(fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss))
			writeData := fmt.Sprintf("%v,%v,%v%%\n", stats.IPAddr.String(), stats.AvgRtt.Milliseconds(), stats.PacketLoss)

			write.WriteString(writeData)
			write.Flush()
			wg.Done()
		}()
	}
	wg.Wait()
}
