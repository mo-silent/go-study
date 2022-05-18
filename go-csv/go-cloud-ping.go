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
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	InputFile   string
	OutPath     string
	LogFile     string
	Regional    string
	SugarLogger *zap.SugaredLogger
	Wg          sync.WaitGroup
	ResRecord   = make(chan []string, 6)
)

func main() {
	flag.StringVar(&InputFile, "infile", "./test.csv", "input csv file")
	flag.StringVar(&OutPath, "outpath", "./", "output conf file")
	flag.StringVar(&Regional, "region", "Asia", "continents")
	flag.StringVar(&LogFile, "log", "./test.log", "log file")
	flag.Parse()
	InitLogger()
	Wg.Add(2)
	go readCsv()
	go readFileRes()
	Wg.Wait()
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
			Wg.Done()
			break
		}
		if err != nil {
			SugarLogger.Fatal(fmt.Sprintf("csv read fail: %v\n", err))
		}
		record[1], _ = simplifiedchinese.GBK.NewDecoder().String(record[1])
		record[2], _ = simplifiedchinese.GBK.NewDecoder().String(record[2])
		// fmt.Println(record[2])
		SugarLogger.Info("read: ", record)
		// recordTmp := [][]string{record}
		ResRecord <- record
	}
	defer close(ResRecord)
}

// GetPing GetPing the IP or AvgRtt by pinging
//
// param domain or ip string
//
// return *ping.Statistics
func GetPing(s string) *ping.Statistics {
	SugarLogger.Info("stat ping IP: ", s)
	pinger, err := ping.NewPinger(s)
	pinger.SetPrivileged(true)
	if err != nil {
		SugarLogger.Error("New Pinger error: ", zap.Error(err))
		return nil
	}
	pinger.Count = 20
	pinger.Timeout = 5 * time.Second
	pinger.PacketsSent = 100
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		SugarLogger.Error("Pinger run error: ", zap.Error(err))
		return nil
	}
	stats := pinger.Statistics()
	SugarLogger.Info("end ping IP: ", s)
	return stats
}

// readFileRes Loop through the CSV file return ping value
func readFileRes() {
	OUTPUTFILE := OutPath + Regional + ".csv"
	fmt.Println(OUTPUTFILE)
	fw, err := os.OpenFile(OUTPUTFILE, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		SugarLogger.Fatal(fmt.Sprintf("open file %v error: %v", OUTPUTFILE, err))
	}
	defer fw.Close()
	write := bufio.NewWriter(fw)
	for {
		resCsv, ok := <-ResRecord
		if !ok {
			defer Wg.Done()
			break
		}
		res := GetPing(resCsv[0])
		if res != nil {
			SugarLogger.Info("write: ", zap.String("source: ", resCsv[0]), res)
			// 写入文件
			writeData := fmt.Sprintf("%v,%v,%v", resCsv[0], resCsv[1], resCsv[2])
			writeData = fmt.Sprintf("%v,%v,%v%%\n", writeData, res.AvgRtt.String(), res.PacketLoss)

			write.WriteString(writeData)
		}
	}
	write.Flush()
}
