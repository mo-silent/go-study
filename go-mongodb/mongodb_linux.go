package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	os.Setenv("CGO_ENABLED", "0")
	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")
	cmd := exec.Command("go", "build", "-o", "mongodb_linux", "./mongodb.go")
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		fmt.Println("***err:", err)
		return
	}
}
