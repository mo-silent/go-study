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
	OPT        string
)

func main() {
	flag.StringVar(&INPUTFILE, "infile", "./tmp.txt", "input file")
	flag.StringVar(&OUTPUTFILE, "outfile", "./tmp.csv", "output file for csv")
	flag.StringVar(&OPT, "opt", "ip", "domain or ip")
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
		s := scanner.Text() // or
		ip := Get(s)

		// 写入文件
		write.WriteString(fmt.Sprintln(s + "," + ip))

	}
	write.Flush()
	fmt.Println("write success!")

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: ip.txt, err: [%v]", err)
	}

}

// Get Get the IP or AvgRtt by pinging
//
// param domain or ip string
//
// return string ip or NAT
func Get(s string) string {
	pinger, err := ping.NewPinger(s)
	pinger.SetPrivileged(true)
	if err != nil {
		log.Println(err)
		return "NAT"
	}
	pinger.Count = 5
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Println(err)
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
