// Author mogd 2022-05-13
// Update mogd 2022-05-13

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-ping/ping"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	INPUTFILE   string
	OUTPUTFILE  string
	LogFile     string
	OPT         string
	SugarLogger *zap.SugaredLogger
)

func main() {
	flag.StringVar(&INPUTFILE, "infile", "./tmp.txt", "input file")
	flag.StringVar(&OUTPUTFILE, "outfile", "./tmp.csv", "output file for csv")
	flag.StringVar(&LogFile, "log", "./tmp.log", "log file")
	flag.StringVar(&OPT, "opt", "ip", "rtt or ip")
	flag.Parse()

	InitLogger()
	defer SugarLogger.Sync()

	f, err := os.Open(INPUTFILE)
	if err != nil {
		SugarLogger.Fatal(fmt.Sprintf("open file %v error: %v", INPUTFILE, err))
	}
	defer f.Close()

	fw, err := os.OpenFile(OUTPUTFILE, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		SugarLogger.Fatal(fmt.Sprintf("open file %v error: %v", OUTPUTFILE, err))
	}
	defer fw.Close()
	write := bufio.NewWriter(fw)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		res := GetPing(s)
		SugarLogger.Info("write: ",
			zap.String("source: ", s),
			zap.String("target: ", res))

		// 写入文件
		write.WriteString(fmt.Sprintln(s + "," + res))

	}
	write.Flush()
	fmt.Println("write success!")

	if err := scanner.Err(); err != nil {
		SugarLogger.Error(
			"Cannot scanner text file!",
			zap.String("file: ", INPUTFILE),
			zap.Error(err))
	}

}

// GetPing GetPing the IP or AvgRtt by pinging
//
// param domain or ip string
//
// return string ip or NAT
func GetPing(s string) string {
	pinger, err := ping.NewPinger(s)
	pinger.SetPrivileged(true)
	if err != nil {
		SugarLogger.Error("New Pinger error: ", zap.Error(err))
		return "NAT"
	}
	pinger.Count = 1
	pinger.Timeout = 5 * time.Second
	pinger.PacketsSent = 5
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		SugarLogger.Error("Pinger run error: ", zap.Error(err))
		return "NAT"
	}
	stats := pinger.Statistics()
	switch OPT {
	case "ip":
		return stats.IPAddr.String()
	case "rtt":
		return stats.AvgRtt.String()
	default:
		return "NAT"
	}
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
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

// getLogWriter set log writer config
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
