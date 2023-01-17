package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

/*
Description: Use go language to implement Linux namespace isolation
Author: Silent.mo
Create: 2023-01-16
Change Activity:
   2023/01/16: init by silent.mo
Annotation: run only linux
	Mount Namespace CLONE_NEWNS
	UTS Namespace CLONE_NEWUTS
	IPC Namespace CLONE_NEWIPC
	PID Namespace CLONE_NEWPID
	Network Namespace CLONE_NEWNET
	User Namespace CLONE_NEWUSER
*/
func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//Cloneflags: syscall.CLONE_NEWUTS, // Run linux
	}
	// user namespace use
	//cmd.SysProcAttr.Credential = &syscall.Credential{
	//	Uid: uint32(1),
	//	Gid: uint32(1),
	//}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
