package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth string) (dirs []string, err error) {
	dirs = make([]string, 0, 20)
	//遍历目录
	err = filepath.Walk(dirPth, func(dirname string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			return nil
		}
		//将目录路径改写成正确格式
		if dirname != "./" {
			dirname = strings.Replace(dirname, "\\", "/", -1)
		}
		if strings.Contains(dirname, ".git") || strings.Contains(dirname, ".idea") {
			return nil
		}
		dirs = append(dirs, dirname)
		return nil
	})
	fmt.Println(dirs)
	return dirs, err
}

//基于windows，执行go fmt命令
func goFmtDirs(path string) {
	cmd := exec.Command("go", "fmt", path)
	stdout, err := cmd.StdoutPipe()
	//获取输出对象，可以从该对象中读取输出结果
	if err != nil {
		log.Fatal(err)
	}
	//保证关闭输出流
	defer stdout.Close()
	//运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(string(opBytes))
	}
}

func main() {
	path := flag.String("path", "D:\\Go/src/gin-kubernetes/", "The dir path")
	flag.Parse()
	dirs, err := WalkDir(*path)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range dirs {
		goFmtDirs(string(v))
	}
}
