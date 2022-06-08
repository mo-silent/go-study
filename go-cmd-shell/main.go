package main

import (
	"flag"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	WG          sync.WaitGroup
	SugarLogger *zap.SugaredLogger
	LogFile     string
	DIR         string
	IP          string = "127.0.0.1"
)

func main() {
	flag.StringVar(&DIR, "dir", "./", "input directory, example /root/zabbix-review-export-import/history/")
	flag.StringVar(&LogFile, "log", "./go-cmd-shell.log", "log file")
	flag.Parse()

	InitLogger()
	defer SugarLogger.Sync()

	files, _ := ioutil.ReadDir(DIR)
	syncCount := 0

	for _, f := range files {
		name := f.Name()
		SugarLogger.Infof("range file : ", name)
		WG.Add(1)
		go zabbixSender(name)
		syncCount += 1

		if syncCount%8 == 0 {
			WG.Wait()
		}
	}

	WG.Wait()
}

// zabbixSender Invoke the Linux command line to execute zabbix_sender, in order to import zabbix historical data
//
// Param string fileName
func zabbixSender(fileName string) {
	if strings.Contains(DIR, "zabbix-proxy-sz2") {
		IP = "39.96.193.87"
	}
	if strings.Contains(DIR, "zabbixproxysz1") {
		IP = "139.196.122.84"
	}
	if strings.Contains(DIR, "zabbixproxyhk1") {
		IP = "172.16.30.16"
	}
	cmdStr := "zabbix_sender  -z " + IP + " -p10051 -NT -vv -i '" + DIR + fileName +
		"' > /var/log/zabbix_sender_log/" + fileName + ".log"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	err := cmd.Run()
	SugarLogger.Infof("Command finished with error: ", err)
	// SugarLogger.Info("write: ", string(out[:]))
	WG.Done()
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
