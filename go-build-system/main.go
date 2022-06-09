// Generate executable files for different operating systems
//
// 生成不同操作系统的可执行文件
//
// Author  mogd  2022-05-06 CST
//
// Update  mogd  2022-05-06 CST

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

// CGO_ENABLED C version of the GO compiler, 1 or 0
//
// GOOS operating system
//
// GOARCH architecture
//
// File input file
//
// OutputFile output file
var (
	CGO_ENABLED string
	GOOS        string
	GOARCH      string
	File        string
	OutputFile  string
	GOTRACEBACK string
)

func main() {
	flag.StringVar(&CGO_ENABLED, "cgo", "0", "1 or 0")
	flag.StringVar(&GOOS, "goos", "linux", "The operating system of the target platform(darwin or freebsd or linux or windows)")
	flag.StringVar(&GOARCH, "goarch", "amd64", "The architecture of the target platform(386 or amd64 or arm)")
	flag.StringVar(&GOTRACEBACK, "gotraceback", "crash", "controls the amount of output generated when a Go program fails due to an unrecovered panic or an unexpected runtime condition")
	flag.StringVar(&File, "infile", "../go-ping/main.go", "Files that need to be compiled")
	flag.StringVar(&OutputFile, "outfile", "D:/文件/2022-05-16/go-ping.exe", "File that need to be output")
	flag.Parse()

	fmt.Println(File)
	defer buildFile()
}

// buildFile Edit environment variables and build the executable file
func buildFile() {
	os.Setenv("CGO_ENABLED", CGO_ENABLED)
	os.Setenv("GOOS", GOOS)
	os.Setenv("GOARCH", GOARCH)
	os.Setenv("GOTRACEBACK", "crash")

	var (
		out, stderr bytes.Buffer
	)
	fmt.Println(OutputFile, File)
	cmd := exec.Command("go", "build", "-o", OutputFile, File)
	cmd.Env = os.Environ()
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("err:", stderr.String())
		return
	}
}
