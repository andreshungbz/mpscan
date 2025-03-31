// mpscan is a command-line utility that scans open ports on a target IP address or hostname.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/andreshungbz/mpscan/internal/connection"
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
	)

	flag.Var(&portsF, "ports", "comma-separated list of ports (e.g., -ports=22,80,443)")
	flag.Var(&targetsF, "targets", "comma-separated list of targets (e.g., -targets=localhost,scanme.nmap.org)")

	flag.Parse()

	// fmt.Printf("%v %v %v %v %v %v %v\n", *targetF, *startPortF, *endPortF, *workersF, *timeoutF, portsF, targetsF) // values check

	// SYNCHRONIZATION SETUP & SCANS

	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg)) // multiple progress bars setup

	targets := createTargets(*targetF, targetsF)
	var summaries []scan.Summary

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

	printResults(summaries, *timeoutF, *jsonF)
}

// createTargets aggregates all hostnames from the -target and -targets flags.
//
// If there are no viable values, a default "localhost" is added.
func createTargets(target string, targets []string) []string {
	var results []string

	if target != "" {
		results = append(results, target)
	}

	results = append(results, targets...)

	if len(results) == 0 {
		results = append(results, "localhost")
	}

	return results
}

// printResults prints the banners and summaries from the scan.
func printResults(summaries []scan.Summary, timeout int, outputJSON bool) {
	fmt.Printf("\n[BANNERS]\n\n")
	for _, summary := range summaries {
		connection.PrintBanner(summary, timeout)
	}

	fmt.Printf("\n[SCAN SUMMARY]\n\n")
	for _, summary := range summaries {
		fmt.Printf("%v\n\n", summary)
	}

	if outputJSON {
		// convert slice to summaries into an array of objects
		jsonData, err := json.MarshalIndent(summaries, "", "  ")
		if err != nil {
			fmt.Printf("Error converting summaries to JSON: %v\n", err)
			return
		}

		// write JSON to a file
		filename := time.Now().Format("20060102-150405") + "-mpscan.json"
		err = os.WriteFile(filename, jsonData, 0644) // -rw-r--r-- permissions
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Printf("[JSON OUTPUT SAVED: %s]\n", filename)
	}
}
