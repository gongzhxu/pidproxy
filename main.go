package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	if len(os.Args) <= 1 {
		fmt.Println("usage: pidproxy /path/to/binary -some-flag")
		return
	}

	argsCmd := os.Args[1:]
	log.Printf("pidproxy start with cmd: %v\n", argsCmd)

	cmd := exec.Command(argsCmd[0], argsCmd[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		log.Printf("pidproxy start err: %v\n", err)
		return
	}

	pid := cmd.Process.Pid
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = cmd.Wait()

		var err error
		for {
			if err = syscall.Kill(-pid, 0); err != nil {
				log.Printf("pidproxy kill pid: %d err: %v\n", -pid, err)
				break
			}

			time.Sleep(time.Second)
		}
	}()

	go func() {
		sigProxy := make(chan os.Signal, 1)
		signal.Notify(sigProxy,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGUSR1,
			syscall.SIGUSR2)
		defer signal.Stop(sigProxy)
		var err error
		for sig := range sigProxy {
			log.Printf("pidproxy kill pid: %d, sig: %s begin\n", -pid, sig.String())
			if err = syscall.Kill(-pid, sig.(syscall.Signal)); err != nil {
				log.Printf("pidproxy kill pid: %d, sig: %s, err: %v\n", -pid, sig.String(), err)
			} else {
				log.Printf("pidproxy kill pid: %d, sig: %s success\n", -pid, sig.String())
			}
		}
	}()

	wg.Wait()
	log.Println("pidproxy exit")
}
