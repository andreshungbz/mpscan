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
		targetF    = flag.String("target", "", "The hostname or IP address to be scanned.")
		startPortF = flag.Int("start-port", 1, "The lower bound port to begin scanning.")
		endPortF   = flag.Int("end-port", 1024, "The upper bound port to finish scanning.")
		workersF   = flag.Int("workers", 100, "The number of concurrent goroutines to launch per target.")
		timeoutF   = flag.Int("timeout", 5, "The maximum time in seconds to wait for connections to be established.")

		portsF   scan.PortList
		targetsF scan.TargetList
		jsonF    = flag.Bool("json", false, "Indicates whether to also output a JSON file of the scan results.")

		debugF = flag.Bool("debug", false, "Displays flag values for debugging.")
	)

	flag.Var(&portsF, "ports", "Comma-separated list of ports (e.g., -ports=22,80,443). Setting this overrides -start-port and -end-port.")
	flag.Var(&targetsF, "targets", "Comma-separated list of targets (e.g., -targets=localhost,scanme.nmap.org). Targets are aggregated with -target.")

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
