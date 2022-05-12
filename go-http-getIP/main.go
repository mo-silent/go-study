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
	Client *http.Client
	URL    = "http://mip.chinaz.com/?query="
)

func main() {
	Client = &http.Client{Timeout: 10 * time.Second}

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
	for scanner.Scan() {
		line := scanner.Text() // or
		//line := scanner.Bytes()
		getRes := Get(line)
		r := regexp.MustCompile(`<td class="z-tc">\s*(.*?)\s*<br />`)
		address := r.FindAllStringSubmatch(getRes, -1)

		var addressResponse string
		for _, v := range address {
			fmt.Println(v[1])
			addressResponse = fmt.Sprintln(line + "," + v[1])
		}

		// 写入文件

		write.WriteString(addressResponse)

	}
	write.Flush()

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: ip.txt, err: [%v]", err)
	}
}

func Get(ip string) string {
	res, err := Client.Get(URL + ip)
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
			panic(err)
		}
	}

	return result.String()
}
