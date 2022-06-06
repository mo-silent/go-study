package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {
	files, _ := ioutil.ReadDir("/root/zabbix-review-export-import/history/")
	for _, f := range files {
		cmd := "zabbix_sender  -z 172.16.30.16 -p10051 -NT -vv -i /root/zabbix-review-export-import/history/" + f.Name() +
			" > /root/zabbix-review-export-import/history_log/" + f.Name() + ".log"
		out, err := exec.Command("/bin/sh", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(out))
	}

}
