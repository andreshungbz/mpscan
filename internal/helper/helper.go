// Package helper contains abstracted functions used in the main program
package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/andreshungbz/mpscan/internal/connection"
	"github.com/andreshungbz/mpscan/internal/scan"
)

// CreateTargets aggregates all hostnames from the -target and -targets flags.
//
// If there are no viable values, a default "localhost" is added.
func CreateTargets(target string, targets []string) []string {
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

// PrintResults prints the banners and summaries from the scan.
func PrintResults(summaries []scan.Summary, timeout int, outputJSON bool) {
	var wg sync.WaitGroup
	fmt.Printf("\n[BANNERS]\n\n")
	for _, summary := range summaries {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connection.PrintBanner(summary, timeout)
		}()
	}
	wg.Wait()

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

// ValidateSEPorts adjusts the values of -start-port and -end-port if they are out of the valid 1-65535 port range
func ValidateSEPorts(startPort *int, endPort *int) {
	if *startPort < 1 || *startPort > 65535 {
		*startPort = 1
	}

	if *endPort < 1 || *endPort > 65535 {
		*endPort = 1024
	}
}
