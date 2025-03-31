// Package helper contains abstracted functions used in the main program
package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/andreshungbz/mpscan/internal/connection"
	"github.com/andreshungbz/mpscan/internal/scan"
)

// createTargets aggregates all hostnames from the -target and -targets flags.
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

// printResults prints the banners and summaries from the scan.
func PrintResults(summaries []scan.Summary, timeout int, outputJSON bool) {
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
