// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

import (
	"flag"
	"fmt"

	"github.com/andreshungbz/mpscan/internal/connection"
	"github.com/andreshungbz/mpscan/internal/scan"
)

var (
	targetF    = flag.String("target", "localhost", "the hostname or IP address to be scanned")
	startPortF = flag.Int("start-port", 1, "the lower bound port to begin scanning")
	endPortF   = flag.Int("end-port", 1024, "the upper bound port to finish scanning")
	workersF   = flag.Int("workers", 100, "the number of concurrent goroutines to launch")
	timeoutF   = flag.Int("timeout", 5, "the maximum time in seconds to wait for a connection to be established")
)

func main() {
	flag.Parse()
	fmt.Printf("%v %v %v %v %v\n", *targetF, *startPortF, *endPortF, *workersF, *timeoutF)

	flags := scan.Flags{
		Target:    targetF,
		StartPort: startPortF,
		EndPort:   endPortF,
		Workers:   workersF,
		Timeout:   timeoutF,
	}

	summary := connection.CreateSummary(flags)
	fmt.Printf("%v\n", summary)
}
