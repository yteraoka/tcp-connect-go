package main

import (
	"crypto/tls"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

// goreleaser
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

var tlsConfig = &tls.Config{}
var opts Options

type Options struct {
	Timeout         int    `short:"t" long:"timeout" default:"3" description:"Connect timeout in second."`
	Count           int    `short:"c" long:"count" default:"1" description:"Number of request per threads."`
	Parallel        int    `short:"p" long:"parallel" default:"1" description:"Number of threads."`
	ServerName      string `long:"servername" description:"Server Name Indication extension in TLS handshake."`
	ShowThresholdMs int64  `short:"s" long:"show-threshold" default:"200" description:"Duration time threshold in millisecond. Show duration time if over this."`
	Verbose         bool   `short:"v" long:"verbose" description:"Enable verbose output."`
	Version         bool   `short:"V" long:"version" description:"Show version and exit."`
	SleepMs         int    `short:"S" long:"sleep" default:"0" description:"Sleep time in millisecond after each connect."`
	SleepRandom     bool   `short:"R" long:"sleep-random" description:"Randomize sleep time."`
	UseTLS          bool   `long:"tls" description:"Do TLS handshake."`
	Insecure        bool   `short:"k" long:"insecure" description:"Disable certificate verification."`
	Args            struct {
		Destination string `description:"servername:port"`
	} `positional-args:"yes"`
}

func times(id int, f func(int, int, string, time.Duration), count int, remote string, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= count; i++ {
		f(id, i, remote, timeout)
	}
}

func connect(routine_id, counter int, remote string, timeout time.Duration) {
	d := net.Dialer{Timeout: timeout}
	start := time.Now()
	var conn net.Conn
	var err error
	if opts.UseTLS {
		conn, err = tls.DialWithDialer(&d, "tcp", remote, tlsConfig)
	} else {
		conn, err = d.Dial("tcp4", remote)
	}
	if err != nil {
		log.Printf("[%03d-%05d] Can't connect to %s: %v\n", routine_id, counter, remote, err)
		return
	}
	end := time.Now()
	diff := end.Sub(start).Milliseconds()
	if opts.Verbose || diff > opts.ShowThresholdMs {
		log.Printf("[%03d-%05d] %d ms\n", routine_id, counter, diff)
	}
	defer conn.Close()
	if opts.SleepMs > 0 {
		if opts.SleepRandom {
			time.Sleep(time.Duration(rand.Intn(opts.SleepMs)) * time.Millisecond)
		} else {
			time.Sleep(time.Duration(opts.SleepMs) * time.Millisecond)
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("version: %s\ncommit: %s\ndate: %s\n", version, commit, date)
		os.Exit(0)
	}

	if opts.Args.Destination == "" {
		log.Fatal("destination parameter is required")
	}

	if opts.UseTLS {
		tlsConfig = &tls.Config{ServerName: opts.ServerName, InsecureSkipVerify: opts.Insecure}
	}

	var wg sync.WaitGroup

	wg.Add(opts.Parallel)

	for i := 0; i < opts.Parallel; i++ {
		go times(i, connect, opts.Count, opts.Args.Destination, time.Duration(opts.Timeout)*time.Second, &wg)
	}

	wg.Wait()
}
