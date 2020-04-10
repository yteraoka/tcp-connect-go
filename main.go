package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var thresholdMs int64
var verbose bool

func times(f func(string, time.Duration), count int, remote string, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		f(remote, timeout)
	}
}

func connect(remote string, timeout time.Duration) {
	d := net.Dialer{Timeout: timeout}
	start := time.Now()
	conn, err := d.Dial("tcp4", remote)
	if err != nil {
		fmt.Printf("Can't connect to %s: %v", remote, err)
		return
	}
	end := time.Now()
	diff := end.Sub(start).Milliseconds()
	if verbose || diff > thresholdMs {
		fmt.Printf("%d ms\n", diff)
	}
	defer conn.Close()
}

func main() {
	var timeout = flag.Int("timeout", 10, "Connect Timeout in Second")
	var count = flag.Int("count", 1, "Number of connect")
	var parallel = flag.Int("parallel", 1, "Number of go routines")
	flag.Int64Var(&thresholdMs, "threshold", 200, "show duration time if over this")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s [-timeout N] [-count M] [-parallel X] [-threshold Y] [-verbose] server:port", os.Args[0])
		os.Exit(1)
	}

	remoteHostPort := flag.Args()[0]

	var wg sync.WaitGroup

	wg.Add(*parallel)

	for i := 0; i < *parallel; i++ {
		go times(connect, *count, remoteHostPort, time.Duration(*timeout)*time.Second, &wg)
	}

	wg.Wait()
}
