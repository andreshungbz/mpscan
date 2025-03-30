// Package connection contains function abtractions for scanning hostnames and ports
package connection

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/andreshungbz/mpscan/internal/scan"
)

// createDialer returns an instance of [net.Dialer] with a custom [net.Dialer.Timeout] set.
func createDialer(seconds int) net.Dialer {
	return net.Dialer{
		Timeout: time.Duration(seconds) * time.Second,
	}
}

// createWorker receives from a [scan.Address] channel and attempts to establish a TCP connection with [net.Dialer].
//
// Ports are added to [scan.Summary.OpenPorts] on successful connection.
// Three attempts are made to establish a connection, with each failed attempt enacting a timer that increases via the an [exponential backoff algorithm].
//
// [exponential backoff algorithm]: https://en.wikipedia.org/wiki/Exponential_backoff
func createWorker(addresses chan scan.Address, dialer net.Dialer, summary *scan.Summary, mu *sync.Mutex) {
	maxRetries := 3

	for address := range addresses {
		target := net.JoinHostPort(address.Hostname, strconv.Itoa(address.Port))

		for i := range maxRetries {
			conn, err := dialer.Dial("tcp", target)

			// on successful connection, add the port, close the connection, and exit loop
			if err == nil {
				mu.Lock()
				summary.AddPort(address.Port)
				mu.Unlock()

				conn.Close()
				break
			}

			backoff := time.Duration(1<<i) * time.Second
			time.Sleep(backoff) // apply exponential backoff timer
		}
	}
}

// CreateSummary concurrently scans the ports of a target hostname based on [scan.Flags] and returns a [scan.Summary].
func CreateSummary(flags scan.Flags) scan.Summary {
	var wg sync.WaitGroup
	var mu sync.Mutex

	addresses := make(chan scan.Address, 100)
	dialer := createDialer(*flags.Timeout)

	var summary scan.Summary
	summary.Hostname = *flags.Target
	summary.TotalPortsScanned = *flags.EndPort - *flags.StartPort + 1

	startTime := time.Now() // timer for the concurrent scan

	// launch goroutines
	for i := 1; i <= *flags.Workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			createWorker(addresses, dialer, &summary, &mu)
		}()
	}

	// send addresses to channel
	for port := *flags.StartPort; port <= *flags.EndPort; port++ {
		addresses <- scan.Address{Hostname: *flags.Target, Port: port}
	}

	close(addresses)
	wg.Wait() // wait for all goroutines to complete

	summary.TimeTaken = time.Since(startTime) // record the total scan time

	return summary
}
