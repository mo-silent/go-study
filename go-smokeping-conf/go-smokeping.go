// Author mogd 2022-05-17
//
// Update mogd 2022-05-18
//
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type SmokePing struct {
	Targets string
	Menu    string
	Tittle  string
	Host    string
	Alerts  string
}

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
	go readCsv()
	for i := 0; i <= 6; i++ {
		Wg.Add(1)
		go readFileRes()
	}
	Wg.Wait()
	mergeFile()
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
			break
		}
		if err != nil {
			SugarLogger.Fatal(fmt.Sprintf("csv read fail: %v\n", err))
		}
		record[1], _ = simplifiedchinese.GBK.NewDecoder().String(record[1])
		record[2], _ = simplifiedchinese.GBK.NewDecoder().String(record[2])
		// fmt.Println(record[2])
		// recordTmp := [][]string{record}
		ResRecord <- record
	}
	defer close(ResRecord)
}

// readFileRes Loop through the CSV file return value
func readFileRes() {
	for {
		res, ok := <-ResRecord
		if !ok {
			defer Wg.Done()
			break
		}
		var smokeConf SmokePing
		cloudVendors := res[2]
		// fmt.Println(cloudVendors)
		smokeConf.Menu = res[2]
		smokeConf.Tittle = res[1]
		smokeConf.Host = res[0]
		smokeConf.Alerts = "someloss"
		switch cloudVendors {
		case "微软云":
			smokeConf.Targets = "+++ Azure_" + uuid.New().String()
			writeConf(cloudVendors, smokeConf)
		case "谷歌云":
			smokeConf.Targets = "+++ Google_" + uuid.New().String()
			writeConf(cloudVendors, smokeConf)
		case "亚马逊云":
			smokeConf.Targets = "+++ Aws_" + uuid.New().String()
			writeConf(cloudVendors, smokeConf)
		case "华为云":
			smokeConf.Targets = "+++ Huawei_" + uuid.New().String()
			writeConf(cloudVendors, smokeConf)
		case "阿里云":
			smokeConf.Targets = "+++ Alibaba_" + uuid.New().String()
			writeConf(cloudVendors, smokeConf)
		case "腾讯云":
			smokeConf.Targets = "+++ Tencent_" + uuid.New().String()
			writeConf(cloudVendors, smokeConf)
		default:
			SugarLogger.Info("Does not exist cloud vendors!")
		}
	}
}

// writeConf write smokeping conf file
func writeConf(filePrefix string, smokeConf SmokePing) {
	// fmt.Println(filePrefix)
	fileName := OutPath + Regional + "_" + filePrefix + "_tmp.txt"
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		SugarLogger.Fatal(fmt.Sprintf("open file error: %v\n", err))
	}
	defer f.Close()

	write := bufio.NewWriter(f)

	content := fmt.Sprintf("%v\nmenu = %v\ntitle = %v\nalerts = %v\nhost = %v\n", smokeConf.Targets, smokeConf.Menu, smokeConf.Tittle, smokeConf.Alerts, smokeConf.Host)
	write.WriteString(content)
	write.Flush()
}

// mergeFile Merge files that are split by cloud vendor
func mergeFile() {
	filePrefix := []string{"微软云", "谷歌云", "亚马逊云", "华为云", "阿里云", "腾讯云"}
	fw, err := os.OpenFile(OutPath+Regional+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		SugarLogger.Fatal(fmt.Sprintf("open file error: %v\n", err))
	}
	defer fw.Close()

	for i := range filePrefix {
		fr, err := os.OpenFile(OutPath+Regional+"_"+filePrefix[i]+"_tmp.txt", os.O_RDONLY, 0644)
		// defer
		if err != nil {
			SugarLogger.Error(fmt.Sprintf("open file %v error: %v\n", filePrefix[i], err))
			continue
		}
		firstSmoke := fmt.Sprintf("+ %v\nmenu = %v\ntitle = %v\n", Regional, filePrefix[i], filePrefix[i])
		fw.WriteString(firstSmoke)
		secondHost := "/" + Regional + "/" + Regional + "_" + filePrefix[i] + "/"
		secondSmoke := fmt.Sprintf("++ %v_%v\nmenu = %v\ntitle = %v\nhost = %v\n", Regional, filePrefix[i], filePrefix[i], filePrefix[i], secondHost)
		fw.WriteString(secondSmoke)
		data, _ := ioutil.ReadAll(fr)
		fw.Write(data)
		fr.Close()
	}
}
