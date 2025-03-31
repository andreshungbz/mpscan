// mpscan is a command-line utility that scans open ports on a target IP address or hostname.
package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/andreshungbz/mpscan/internal/connection"
	"github.com/andreshungbz/mpscan/internal/helper"
	"github.com/andreshungbz/mpscan/internal/scan"
	"github.com/vbauerster/mpb/v8"
)

func main() {

	// FLAG VARIABLES SETUP

	var (
		targetF    = flag.String("target", "", "the hostname or IP address to be scanned")
		startPortF = flag.Int("start-port", 1, "the lower bound port to begin scanning")
		endPortF   = flag.Int("end-port", 1024, "the upper bound port to finish scanning")
		workersF   = flag.Int("workers", 100, "the number of concurrent goroutines to launch")
		timeoutF   = flag.Int("timeout", 5, "the maximum time in seconds to wait for a connection to be established")

		portsF   scan.PortList
		targetsF scan.TargetList
		jsonF    = flag.Bool("json", false, "indicates whether to also output JSON")

		debugF = flag.Bool("debug", false, "displays flag values for debugging")
	)

	flag.Var(&portsF, "ports", "comma-separated list of ports (e.g., -ports=22,80,443)")
	flag.Var(&targetsF, "targets", "comma-separated list of targets (e.g., -targets=localhost,scanme.nmap.org)")

	flag.Parse()

	// change start and end ports to defaults if they are outside the valid port range 1-65535
	helper.ValidateSEPorts(startPortF, endPortF)

	// SYNCHRONIZATION SETUP & SCANS

	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg)) // multiple progress bars setup

	targets := helper.CreateTargets(*targetF, targetsF)
	var summaries []scan.Summary

	// print all flag values if -debug is set
	if *debugF {
		fmt.Printf("\n[DEBUG]\n\n")
		fmt.Printf("%-15s %s\n", "[-target]", *targetF)
		fmt.Printf("%-15s %d\n", "[-start-port]", *startPortF)
		fmt.Printf("%-15s %d\n", "[-end-port]", *endPortF)
		fmt.Printf("%-15s %d\n", "[-workers]", *workersF)
		fmt.Printf("%-15s %d\n", "[-timeout]", *timeoutF)
		fmt.Printf("%-15s %v\n", "[-ports]", portsF)
		fmt.Printf("%-15s %v\n", "[-targets]", targetsF)
		fmt.Printf("%-15s %v\n", "[-json]", *jsonF)
		fmt.Printf("%-15s %v\n", "[-debug]", *debugF)
		fmt.Printf("%-15s %v\n", "[TARGETS]", targets)
	}

	fmt.Printf("\n[SCAN START]\n\n")
	for _, target := range targets {
		flags := scan.Flags{
			Target:    target,
			StartPort: *startPortF,
			EndPort:   *endPortF,
			Workers:   *workersF,
			Timeout:   *timeoutF,
			Ports:     portsF,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			summaries = append(summaries, connection.CreateSummary(flags, p))
		}()
	}

	wg.Wait()
	p.Wait()

	// RESULTS

	helper.PrintResults(summaries, *timeoutF, *jsonF)
}
