package main

import (
	"bufio"
	"bytes"
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
	URL = "http://mip.chinaz.com/?query="
)

func main() {

	f, err := os.Open("D:/文件/2022-05-12/ip.txt")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer f.Close()

	fw, err := os.OpenFile("D:/文件/2022-05-12/AFK.csv", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer fw.Close()

	write := bufio.NewWriter(fw)

	scanner := bufio.NewScanner(f)

	countLine := int(0)
	for scanner.Scan() {
		line := scanner.Text() // or
		//line := scanner.Bytes()
		if (countLine != 0) && (countLine%20 == 0) {
			fmt.Println(countLine)
			time.Sleep(60 * time.Second)
		}
		getRes := Get(line)
		r := regexp.MustCompile(`<td class="z-tc">\s*(.*?)\s*<br />`)
		address := r.FindAllStringSubmatch(getRes, -1)

		var addressResponse string
		for _, v := range address {
			fmt.Println(v[1])
			addressResponse = fmt.Sprintln(line + "," + v[1])
		}

		// 写入文件
		// fmt.Println(addressResponse)
		write.WriteString(addressResponse)
		countLine += 1

	}
	write.Flush()

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: ip.txt, err: [%v]", err)
	}
}

func Get(ip string) string {
	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Get(URL + ip)
	if err != nil {
		log.Fatal(err)
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
			fmt.Println(ip)
		}
	}

	return result.String()
}
