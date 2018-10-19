package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type cArgs []string

func (i *cArgs) String() string {
	return strings.Join(*i, "\",\"")
}

func (c *cArgs) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func main() {
	flag.CommandLine.Usage = func() {
		flag.Usage()
		os.Exit(0)
	}

	var commands cArgs
	flag.Var(&commands, "c", "需要守护的指令")
	daemon := flag.Bool("d", false, "是否进入后台模式")
	flag.Parse()
	if 0 == len(commands) {
		flag.CommandLine.Usage()
	}
	if *daemon {
		os.Args = append(os.Args, "-d=false")
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Start()
		os.Exit(0)
	}

	fmt.Printf("Command \"%v\" to startup...\n", &commands)
	for _, command := range commands[1:] {
		command := command
		go func() {
			startup(command)
		}()
	}
	startup(commands[0])
}

func startup(command string) {
recycle:
	stime := time.Now()
	cmd := exec.Command("sh", "-c", command)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		println(err.Error())
		return
	}
	err = cmd.Wait()
	if err != nil {
		print(stderr.String())
	}

	span := 10 * time.Second
	if d := span - time.Since(stime); d > 0 {
		fmt.Printf("Waiting for Sleep to finish of \"%v\"...\n", command)
		time.Sleep(d)
	}
	goto recycle
}
