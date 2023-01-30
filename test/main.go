// package main

// // import "C"

// func main() {
// 	var ch chan struct{}
// 	<-ch
// }
package main

import (
	"bytes"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/go-ping/ping"
	"github.com/urfave/cli"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	base := "http://127.0.0.1:28080"
	tp := jira.BasicAuthTransport{
		Username: "silent",
		Password: "123456",
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}
	fmt.Println("ok")
	//issue, _, _ := jiraClient.Issue.Get("SRE-10", nil)
	////currentStatus := issue.Fields.Status.Name
	//fmt.Println(issue.Fields.Assignee)
	//fmt.Println(issue.Fields.Type)
	//fmt.Println(issue.Fields.Status)
	issue := jira.Issue{
		Fields: &jira.IssueFields{
			Type: jira.IssueType{
				Name: "任务",
			},
			Project: jira.Project{
				Key: "SRE",
			},
			Summary:     "Test Issue from Go-JIRA library",
			Description: "Description of my bug report",
			Reporter: &jira.User{
				Name:         "silent",
				EmailAddress: "silent.mo@xx.com",
				DisplayName:  "silent.mo",
			},
			Priority: &jira.Priority{
				Name: "Medium",
			},

			//Status: &jira.Status{
			//	Name: "待办",
			//},
		},
	}
	fmt.Println(issue)
	// Create the issue
	newIssue, resp, err := jiraClient.Issue.Create(&issue)
	if err != nil {
		if resp != nil {
			fmt.Printf("%s\n", resp.Status)
			fmt.Println(err)
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			newStr := buf.String()

			fmt.Printf(newStr)
		}
		fmt.Println(resp)
		os.Exit(0)
	}
	fmt.Printf("%s: %s\n", newIssue.Key, newIssue.Fields.Summary)

}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container",
	Flags: []cli.Flag{
		// integrate -i and -t for convenience
		&cli.BoolFlag{
			Name:  "it",
			Usage: "open an interactive tty(pseudo terminal)",
		},
	},
	Action: func(context *cli.Context) error {
		args := context.Args()
		//if args.Len() == 0 {
		//	return errors.New("Run what?")
		//}
		containerCmd := args.Get(0) // command

		// check whether type `-it`
		tty := context.Bool("it") // presudo terminal

		// 这个函数在下面定义
		fmt.Printf("cmd %v, tty: %v", containerCmd, tty)

		return nil
	},
}

func CobraCommand() {
	app := cli.NewApp()
	app.Name = "test"
	app.Usage = "test"
	app.Commands = []cli.Command{
		runCommand,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Other() {
	addressResponse := "巴西圣保罗 华为"
	cloud := regexp.MustCompile(`(微软云)|(谷歌云)|(亚马逊云)|(华为云)|(阿里云)|(腾讯云)`)
	judgeTmp := cloud.FindAllStringSubmatch(addressResponse, -1)
	fmt.Println(judgeTmp)
	var judge string
	for _, v := range judgeTmp {
		fmt.Println(v)
		judge = v[1]
	}
	if strings.Contains(addressResponse, judge) {
		fmt.Println("test")
	}
	// Ping()
	// return
	// var wg sync.WaitGroup
	// foo := make(chan int)
	// bar := make(chan int)
	// closing := make(chan struct{})
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	select {
	// 	case foo <- <-bar:
	// 		println("bar")
	// 	case <-closing:
	// 		println("closing")
	// 	}
	// }()
	// // bar <- 123
	// close(closing)
	// close(bar)
	// <-foo
	// close(foo)
	// wg.Wait()
	var ar = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
	// var a []byte
	a := ar[2:5]
	a[0] = 'd'
	fmt.Println(ar)
	defer fmt.Println(1)
	defer fmt.Println(2)
	defer fmt.Println(3)
	// type Human struct {
	// 	name   string
	// 	age    int
	// 	weight int
	// }

	// type Student struct {
	// 	Human      // 匿名字段，那么默认Student就包含了Human的所有字段，
	// 	speciality string
	// 	int        // 内置类型作为匿名字段，变量名就是 int
	// }
	// jane := Student{Human: Human{"Jane", 35, 100}, speciality: "Biology", int: 1}
	// fmt.Println("Her preferred number is", jane.int)

	// Bob := Human{"Bob", 39, "sssfgf"}
	// fmt.Println("This Human is : ", Bob)

	// A interview questions
	m := make(map[int]int, 10)
	for i := 1; i <= 10; i++ {
		m[i] = i
	}
	// 闭包使用外部变量，输出会有问题
	// for k, v := range m {
	// 	go func() {
	// 		fmt.Println("k ->", k, "v ->", v)
	// 	}()
	// }
	// 解决，不要闭包直接使用外部变量，通过传参就能解决了
	for k, v := range m {
		go func(a, b int) {
			fmt.Println("k ->", a, "v ->", b)
		}(k, v)
	}

	doc := []interface{}{
		"test",
		"name",
	}
	fmt.Printf("%T\n", doc)
}

type Human struct {
	name  string
	age   int
	phone string
}

// 通过这个方法 Human 实现了 fmt.Stringer
func (h Human) String() string {
	return "❰" + h.name + " - " + strconv.Itoa(h.age) + " years -  ✆ " + h.phone + "❱"
}

func Ping() {
	// pinger, err := ping.NewPinger("www.google.com")
	// if err != nil {
	// 	panic(err)
	// }

	// // Listen for Ctrl-C.
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// go func() {
	// 	for _ = range c {
	// 		pinger.Stop()
	// 	}
	// }()

	// pinger.OnRecv = func(pkt *ping.Packet) {
	// 	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
	// 		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	// }

	pinger, err := ping.NewPinger("39.101.244.245")
	pinger.SetPrivileged(true)
	if err != nil {
		log.Println(err)
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.Count = 5
	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		log.Println(err)
	}

	stats := pinger.Statistics()
	fmt.Println(stats)
}
