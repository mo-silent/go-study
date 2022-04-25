package main

import (
	"net/http"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var sugarLogger *zap.SugaredLogger

// func main() {
// 	InitLogger()
// 	defer sugarLogger.Sync()
// 	simpleHttpGet("www.google.com")
// 	simpleHttpGet("http://www.google.com")
// }

// func InitLogger() {
// 	logger, _ := zap.NewProduction()
// 	sugarLogger = logger.Sugar()
// }

// func simpleHttpGet(url string) {
// 	sugarLogger.Debugf("Trying to hit GET request for %s", url)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
// 	} else {
// 		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
// 		resp.Body.Close()
// 	}
// }

// customize logger
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}
