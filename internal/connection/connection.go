// Package connection contains function abtractions for scanning hostnames and ports
package connection

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
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
				printBanner(conn, address.Hostname, address.Port)

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

// printBanner reads the possible response for an open port, extracts banner information, then prints the result
//
// For HTTP port 80 an HTTP request is sent in order to read the response. For all other ports, it is assumed
// that the server sends a response on TCP connection establishment.
func printBanner(conn net.Conn, hostname string, port int) {
	target := net.JoinHostPort(hostname, strconv.Itoa(port))

	result := make([]byte, 1024)

	var err error
	var n int

	switch port {
	case 80: // HTTP usually requires manually sending a request to get a banner
		_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + hostname + "\r\n\r\n")) // send an HTTP request
		if err == nil {
			_, err = conn.Read(result) // read the HTTP response
			if err == nil {
				res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(result)), nil) // convert to http.Response
				if err == nil {
					serverHeader := res.Header.Get("Server") // get Server header
					copy(result[:], []byte(serverHeader))    // write contents to result
					n = len(serverHeader)                    // adjust length to read
				}
			}
		}
	default: // banners can be protocol-dependent, so default to assuming one is automatically sent on connection
		n, err = conn.Read(result)
	}

	// print the banner if nothing went wrong
	if err == nil {
		banner := strings.TrimSpace(string(result[:n]))
		fmt.Printf("[banner] %s: %s\n", target, banner)
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
