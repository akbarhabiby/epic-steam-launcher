package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

func exit(code int, err error) {
	fmt.Printf("Error: %s\n", err.Error())
	fmt.Println("Press 'Enter' to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(code)
}

func getPID(name string) (uint32, error) {
	h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if e != nil {
		return 0, e
	}
	p := windows.ProcessEntry32{Size: 568}
	for {
		e := windows.Process32Next(h, &p)
		if e != nil {
			return 0, fmt.Errorf("Could not find a single process with name: %s", name)
		}
		if windows.UTF16ToString(p.ExeFile[:]) == name {
			return p.ProcessID, nil
		}
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	if len(os.Args[1:]) != 2 {
		exit(1, fmt.Errorf("Needs launch URL and EXE name"))
	}
	epicUrl := os.Args[1]
	exeName := os.Args[2]
	if epicUrl == "" {
		exit(1, fmt.Errorf("Empty URL"))
	}
	if !strings.Contains(epicUrl, "com.epicgames.launcher") {
		exit(1, fmt.Errorf("Invalid URL"))
	}
	if exeName == "" {
		exit(1, fmt.Errorf("Empty EXE name"))
	}
	go exec.Command(epicUrl).Run()
	fmt.Printf("Starting %s\n", epicUrl)
	time.Sleep(time.Second * 5)
	fmt.Printf("Checking for %s\n", exeName)
	pid, e := getPID(exeName)
	if e != nil {
		exit(2, e)
	}
	fmt.Printf("%s is running with PID %d\n", exeName, pid)
	fmt.Printf("Checking if PID %d is running\n", pid)
	proc, e := os.FindProcess(int(pid))
	if e != nil {
		exit(2, e)
	}

	fmt.Println("Game started. Waiting to exit")
	proc.Wait()
	fmt.Println("Game exited.")
}
