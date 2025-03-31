// mpscan is a command-line utility that scans open ports on a target IP address or hostname.
package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/andreshungbz/mpscan/internal/connection"
	"github.com/andreshungbz/mpscan/internal/scan"
	"github.com/vbauerster/mpb/v8"
)

var (
	targetF    = flag.String("target", "localhost", "the hostname or IP address to be scanned")
	startPortF = flag.Int("start-port", 1, "the lower bound port to begin scanning")
	endPortF   = flag.Int("end-port", 1024, "the upper bound port to finish scanning")
	workersF   = flag.Int("workers", 100, "the number of concurrent goroutines to launch")
	timeoutF   = flag.Int("timeout", 5, "the maximum time in seconds to wait for a connection to be established")
)

func main() {
	var portsF scan.PortList
	flag.Var(&portsF, "ports", "comma-separated list of ports (e.g., -ports=22,80,443)")

	flag.Parse()
	fmt.Printf("%v %v %v %v %v %v\n", *targetF, *startPortF, *endPortF, *workersF, *timeoutF, portsF)

	flags := scan.Flags{
		Target:    *targetF,
		StartPort: *startPortF,
		EndPort:   *endPortF,
		Workers:   *workersF,
		Timeout:   *timeoutF,
		Ports:     portsF,
	}

	flags2 := scan.Flags{
		Target:    "localhost",
		StartPort: *startPortF,
		EndPort:   *endPortF,
		Workers:   *workersF,
		Timeout:   *timeoutF,
		Ports:     portsF,
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var summary scan.Summary
	var summary2 scan.Summary

	// Create a shared progress manager
	p := mpb.New(mpb.WithWaitGroup(&wg))

	go func() {
		defer wg.Done()
		summary = connection.CreateSummary(flags, p)
	}()

	go func() {
		defer wg.Done()
		summary2 = connection.CreateSummary(flags2, p)
	}()

	wg.Wait()
	p.Wait()

	connection.PrintBanner(summary, flags.Timeout)
	connection.PrintBanner(summary2, flags2.Timeout)

	fmt.Printf("%v\n", summary)
	fmt.Printf("%v\n", summary2)
}
