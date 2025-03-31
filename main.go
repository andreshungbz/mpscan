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
	targetF    = flag.String("target", "", "the hostname or IP address to be scanned")
	startPortF = flag.Int("start-port", 1, "the lower bound port to begin scanning")
	endPortF   = flag.Int("end-port", 1024, "the upper bound port to finish scanning")
	workersF   = flag.Int("workers", 100, "the number of concurrent goroutines to launch")
	timeoutF   = flag.Int("timeout", 5, "the maximum time in seconds to wait for a connection to be established")
)

func main() {
	var portsF scan.PortList
	flag.Var(&portsF, "ports", "comma-separated list of ports (e.g., -ports=22,80,443)")

	var targetsF scan.TargetList
	flag.Var(&targetsF, "targets", "comma-separated list of targets (e.g., -targets=example.com,scanme.nmap.org)")

	flag.Parse()
	fmt.Printf("%v %v %v %v %v %v %v\n", *targetF, *startPortF, *endPortF, *workersF, *timeoutF, portsF, targetsF)

	var targets []string

	if *targetF != "" {
		targets = append(targets, *targetF)
	}

	for _, hostname := range targetsF {
		targets = append(targets, hostname)
	}

	if len(targets) == 0 {
		targets = append(targets, "localhost")
	}

	fmt.Printf("\n[SCAN START]\n\n")

	var summaries []scan.Summary
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))

	for _, target := range targets {
		wg.Add(1)

		flags := scan.Flags{
			Target:    target,
			StartPort: *startPortF,
			EndPort:   *endPortF,
			Workers:   *workersF,
			Timeout:   *timeoutF,
			Ports:     portsF,
		}

		go func() {
			defer wg.Done()
			summaries = append(summaries, connection.CreateSummary(flags, p))
		}()
	}

	wg.Wait()
	p.Wait()

	fmt.Printf("\n[BANNERS]\n\n")
	for _, summary := range summaries {
		connection.PrintBanner(summary, *timeoutF)
	}

	fmt.Printf("\n[SCAN SUMMARY]\n\n")
	for _, summary := range summaries {
		fmt.Printf("%v\n\n", summary)
	}
}
