package connection

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/andreshungbz/mpscan/internal/scan"
)

func createDialer(seconds int) net.Dialer {
	return net.Dialer{
		Timeout: time.Duration(seconds) * time.Second,
	}
}

func createWorker(tasks chan scan.Address, dialer net.Dialer, summary *scan.Summary, mu *sync.Mutex) {
	maxRetries := 3
	for address := range tasks {
		summary.TotalPortsScanned++
		target := net.JoinHostPort(address.Hostname, strconv.Itoa(address.Port))
		var success bool

		for i := range maxRetries {
			conn, err := dialer.Dial("tcp", target)
			if err == nil {
				conn.Close()
				fmt.Printf("Connection to %s was successful\n", target)
				mu.Lock()
				summary.AddPort(address.Port)
				mu.Unlock()
				success = true
				break
			}

			backoff := time.Duration(1<<i) * time.Second
			fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1, target, backoff)
			time.Sleep(backoff)
		}

		if !success {
			fmt.Printf("Failed to connect to %s after %d attempts\n", target, maxRetries)
		}
	}
}

func CreateSummary(summaries chan scan.Summary, flags scan.Flags) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	tasks := make(chan scan.Address, 100)
	dialer := createDialer(*flags.Timeout)

	var summary scan.Summary
	summary.Hostname = *flags.Target

	startTime := time.Now()

	for i := 1; i <= *flags.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createWorker(tasks, dialer, &summary, &mu)
		}()
	}

	for port := *flags.StartPort; port <= *flags.EndPort; port++ {
		tasks <- scan.Address{Hostname: *flags.Target, Port: port}
	}
	close(tasks)

	wg.Wait()
	summary.TimeTaken = time.Since(startTime)
	summaries <- summary
}
