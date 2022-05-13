// Author mogd 2022-05-13
// Update mogd 2022-05-13

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-ping/ping"
)

var (
	INPUTFILE  string
	OUTPUTFILE string
)

func main() {
	flag.StringVar(&INPUTFILE, "infile", "./tmp.txt", "input file")
	flag.StringVar(&OUTPUTFILE, "outfile", "./tmp.csv", "output file for csv")
	flag.Parse()

	f, err := os.Open(INPUTFILE)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer f.Close()

	fw, err := os.OpenFile(OUTPUTFILE, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer fw.Close()
	write := bufio.NewWriter(fw)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		domain := scanner.Text() // or
		ip := GetIP(domain)

		// 写入文件
		write.WriteString(fmt.Sprintln(domain + "," + ip))

	}
	write.Flush()
	fmt.Println("write success!")

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: ip.txt, err: [%v]", err)
	}

}

// GetIP Get the IP by pinging the domain name
//
// param domain string
//
// return string ip or NAT
func GetIP(domain string) string {
	pinger, err := ping.NewPinger(domain)
	pinger.SetPrivileged(true)
	if err != nil {
		log.Println(err)
		return "NAT"
	}
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Println(err)
		return "NAT"
	}
	stats := pinger.Statistics()
	return stats.IPAddr.String()
}
