package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// goreleaser
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

var thresholdMs int64
var verbose bool
var useTls bool
var tlsConfig = &tls.Config{}

func times(f func(string, time.Duration), count int, remote string, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		f(remote, timeout)
	}
}

func connect(remote string, timeout time.Duration) {
	d := net.Dialer{Timeout: timeout}
	start := time.Now()
	var conn net.Conn
	var err error
	if useTls {
		conn, err = tls.DialWithDialer(&d, "tcp", remote, tlsConfig)
	} else {
		conn, err = d.Dial("tcp4", remote)
	}
	if err != nil {
		fmt.Printf("Can't connect to %s: %v\n", remote, err)
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
	var timeout = flag.Int("timeout", 10, "Connect timeout in second")
	var count = flag.Int("count", 1, "Number of connect")
	var parallel = flag.Int("parallel", 1, "Number of go routines")
	var serverName = flag.String("servername", "", "SNI")
	flag.Int64Var(&thresholdMs, "threshold", 200, "Duration time threshold in millisecond. Show duration time if over this")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode (show duration time every time)")
	flag.BoolVar(&useTls, "tls", false, "Enable TLS mode")
	var insecure = flag.Bool("insecure", false, "Do not verify certificate")
	var versionFlag = flag.Bool("version", false, "Show version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("%v, commit %v, built at %v", version, commit, date)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s [-timeout N] [-count M] [-parallel X] [-threshold Y] [-verbose] [-tls] [-insecure] [-servername www.example.com] server:port\n", os.Args[0])
		os.Exit(1)
	}

	remoteHostPort := flag.Args()[0]

	if useTls {
		if *serverName == "" {
			*serverName = strings.Split(remoteHostPort, ":")[0]
		}
		tlsConfig = &tls.Config{ServerName: *serverName, InsecureSkipVerify: *insecure}
	}

	var wg sync.WaitGroup

	wg.Add(*parallel)

	for i := 0; i < *parallel; i++ {
		go times(connect, *count, remoteHostPort, time.Duration(*timeout)*time.Second, &wg)
	}

	wg.Wait()
}
