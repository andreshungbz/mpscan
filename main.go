// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/andreshungbz/mpscan/internal/connection"
	"github.com/andreshungbz/mpscan/internal/scan"
)

var targetF = flag.String("target", "localhost", "the hostname or IP address to be scanned")
var startPortF = flag.Int("start-port", 1, "the lower bound port to begin scanning")
var endPortF = flag.Int("end-port", 1024, "the upper bound port to finish scanning")
var workersF = flag.Int("workers", 100, "the number of concurrent goroutines to launch")
var timeoutF = flag.Int("timeout", 5, "the maximum time in seconds to wait for a connection to be established")

func main() {
	flag.Parse()
	fmt.Printf("%v %v %v %v %v\n", *targetF, *startPortF, *endPortF, *workersF, *timeoutF)

	var wg sync.WaitGroup
	tasks := make(chan scan.Address, 100)
	dialer := connection.CreateDialer(*timeoutF)

	var summary scan.Summary
	summary.Hostname = *targetF

	startTime := time.Now()
	for i := 1; i <= *workersF; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connection.CreateWorker(tasks, dialer, &summary)
		}()
	}

	for port := *startPortF; port <= *endPortF; port++ {
		tasks <- scan.Address{Hostname: *targetF, Port: port}
	}

	close(tasks)

	wg.Wait()

	summary.TimeTaken = time.Since(startTime)
	fmt.Printf("%v\n", summary)
}
