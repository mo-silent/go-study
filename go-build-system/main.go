// Generate executable files for different operating systems
//
// 生成不同操作系统的可执行文件
//
// Author  mogd  2022-05-06 CST
//
// Update  mogd  2022-05-06 CST

package main

import (
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
)

func main() {
	flag.StringVar(&CGO_ENABLED, "cgo", "0", "1 or 0")
	flag.StringVar(&GOOS, "goos", "linux", "The operating system of the target platform(darwin or freebsd or linux or windows)")
	flag.StringVar(&GOARCH, "goarch", "amd64", "The architecture of the target platform(386 or amd64 or arm)")
	flag.StringVar(&File, "infile", "../go-mongodb/mongodb.go", "Files that need to be compiled")
	flag.StringVar(&OutputFile, "outfile", "main.exe", "File that need to be output")
	flag.Parse()

	buildFile()
}

// buildFile Edit environment variables and build the executable file
func buildFile() {
	os.Setenv("CGO_ENABLED", CGO_ENABLED)
	os.Setenv("GOOS", GOOS)
	os.Setenv("GOARCH", GOARCH)

	cmd := exec.Command("go", "build", "-o", OutputFile, File)
	cmd.Env = os.Environ()

	err := cmd.Run()
	if err != nil {
		fmt.Println("***err:", err)
		return
	}
}
