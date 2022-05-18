// Author mogd 2022-05-12
//
// Update mogd 2022-05-17

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	// Client *http.Client
	INPUTFILE  string
	OUTPUTFILE string
	URL        = "http://mip.chinaz.com/?query="
)

func main() {
	flag.StringVar(&INPUTFILE, "infile", "D:/文件/2022-05-12/ip.txt", "input file")
	flag.StringVar(&OUTPUTFILE, "outfile", "D:/文件/2022-05-12/ip.csv", "output file for csv")
	flag.Parse()

	f, err := os.Open(INPUTFILE)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer f.Close()

	fw, err := os.OpenFile(OUTPUTFILE, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("文件打开失败", err)
	}
	defer fw.Close()

	write := bufio.NewWriter(fw)

	scanner := bufio.NewScanner(f)

	countLine := int(0)
	for scanner.Scan() {
		line := scanner.Text() // or
		//line := scanner.Bytes()
		if (countLine != 0) && (countLine%20 == 0) {
			log.Println(countLine)
			time.Sleep(20 * time.Second)
		}
		getRes := Get(line)
		r := regexp.MustCompile(`<td class="z-tc">\s*(.*?)\s*<br />`)
		address := r.FindAllStringSubmatch(getRes, -1)

		var addressResponse string
		for _, v := range address {
			// log.Println(v[1])
			addressResponse = fmt.Sprintf(line + "," + v[1])
		}

		// cloud := regexp.MustCompile(`(微软云)|(谷歌云)|(亚马逊云)|(华为云)|(阿里云)|(腾讯云)`)
		// judgeTmp := cloud.FindAllStringSubmatch(addressResponse, -1)
		// var judge string
		// for _, v := range judgeTmp {
		// 	judge = v[0]
		// }
		// if strings.Contains(addressResponse, judge) && len(judge) != 0 {
		// 	log.Printf("%v : %v: %v\n", line, judge, len(judge))
		// 	write.WriteString(fmt.Sprintln(addressResponse + "," + judge))
		// } else {
		// 	log.Printf("%v\n", line)
		// 	write.WriteString(fmt.Sprintln(addressResponse + ",other"))
		// }
		// 写入文件
		fmt.Println(addressResponse)
		write.WriteString(fmt.Sprintln(addressResponse))
		countLine += 1

	}
	write.Flush()

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: ip.txt, err: [%v]", err)
	}
}

// Get Get the physical address over IP
//
// param ip string
//
// return string address
func Get(ip string) string {
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(URL + ip)
	if err != nil {
		log.Println(err)
		return "NAT"
	}
	defer res.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := res.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			log.Println(ip)
			return "NAT"
		}
	}

	return result.String()
}
