// Package connection contains function abstractions for scanning hostnames and ports
package connection

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andreshungbz/mpscan/internal/scan"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// CreateSummary concurrently scans the ports of a target hostname based on [scan.Flags] and returns a [scan.Summary].
//
// If [scan.Flags.Ports] contains values, it overrides scanning based on [scan.Flags.StartPort] and [scan.Flags.EndPort].
func CreateSummary(flags scan.Flags, p *mpb.Progress) scan.Summary {
	var wg sync.WaitGroup
	var mu sync.Mutex

	addresses := make(chan scan.Address, 100)
	dialer := createDialer(flags.Timeout)

	// portsOverride defines whether or not -ports overrides -start-port and -end-port
	portsOverride := len(flags.Ports) > 0

	var summary scan.Summary
	summary.Hostname = flags.Target

	if !portsOverride {
		summary.TotalPortsScanned = flags.EndPort - flags.StartPort + 1
	} else {
		summary.TotalPortsScanned = len(flags.Ports)
	}

	startTime := time.Now() // timer for the concurrent scan

	bar := p.AddBar(int64(summary.TotalPortsScanned),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("%s: scanning ", summary.Hostname)),
			decor.CountersNoUnit("%d/%d ports", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)

	// launch goroutines
	for i := 1; i <= flags.Workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			createWorker(addresses, dialer, &summary, &mu, bar)
		}()
	}

	// send addresses to channel
	if !portsOverride {
		for port := flags.StartPort; port <= flags.EndPort; port++ {
			addresses <- scan.Address{Hostname: flags.Target, Port: port}
		}
	} else {
		for _, port := range flags.Ports {
			addresses <- scan.Address{Hostname: flags.Target, Port: port}
		}
	}

	close(addresses)
	wg.Wait() // wait for all goroutines to complete

	summary.TimeTaken = time.Since(startTime) // record the total scan time

	return summary
}

// PrintBanner iterates through [scan.Summary.OpenPorts] and attempts to print banner information of each
//
// For HTTP port 80 an HTTP request is sent in order to read the response. For all other ports, it is assumed
// that the server sends a response on TCP connection establishment.
func PrintBanner(summary scan.Summary, timeout int) {
	dialer := createDialer(timeout)

	for _, port := range summary.OpenPorts {
		target := net.JoinHostPort(summary.Hostname, strconv.Itoa(port))

		conn, connErr := dialer.Dial("tcp", target)

		if connErr == nil {
			defer conn.Close()
			conn.SetDeadline(time.Now().Add(dialer.Timeout))
			result := make([]byte, 1024)

			var err error
			var n int

			switch port {
			case 80: // HTTP usually requires manually sending a request to get a banner
				_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + summary.Hostname + "\r\n\r\n")) // send an HTTP request
				if err == nil {
					_, err = conn.Read(result) // read the HTTP response
					if err == nil {
						res, resErr := http.ReadResponse(bufio.NewReader(bytes.NewReader(result)), nil) // convert to http.Response
						err = resErr
						if resErr == nil {
							serverHeader := res.Header.Get("Server") // get Server header

							if serverHeader == "" {
								serverHeader = res.Header.Get("X-Powered-By") // check X-Powered-By if Server is empty
							}

							if serverHeader == "" {
								serverHeader = res.Header.Get("Date") // get date if none others match.
							}

							if serverHeader == "" {
								err = fmt.Errorf("no banner information found in headers") // assign arbitrary error
							}

							copy(result[:], []byte(serverHeader)) // write contents to result
							n = len(serverHeader)                 // adjust length to read
						}
					}
				}

			default: // banners can be protocol-dependent, so default to assuming one is automatically sent on connection
				n, err = conn.Read(result)
			}

			// print the banner if nothing went wrong
			if err == nil {
				banner := strings.TrimSpace(string(result[:n]))
				fmt.Printf("[%s] %s\n", target, banner)
			}
		}
	}

}

// HELPER FUNCTIONS

// createWorker receives from a [scan.Address] channel and attempts to establish a TCP connection with [net.Dialer].
//
// Ports are added to [scan.Summary.OpenPorts] on successful connection.
// Three attempts are made to establish a connection, with each failed attempt enacting a timer that increases via the an [exponential backoff algorithm].
// A random value from 0.0 to 1.0 is also multiplied to the timer to add variance.
//
// [exponential backoff algorithm]: https://en.wikipedia.org/wiki/Exponential_backoff
func createWorker(addresses chan scan.Address, dialer net.Dialer, summary *scan.Summary, mu *sync.Mutex, bar *mpb.Bar) {
	maxRetries := 1

	for address := range addresses {
		target := net.JoinHostPort(address.Hostname, strconv.Itoa(address.Port))

		for i := 1; i <= maxRetries; i++ {
			conn, err := dialer.Dial("tcp", target)

			// on successful connection, add the port, close the connection, and exit loop
			if err == nil {
				mu.Lock()
				summary.AddPort(address.Port)
				mu.Unlock()

				conn.Close()
				break
			}

			inital := 1 << i
			backoff := time.Duration(float64(inital)*rand.Float64()) * time.Second
			time.Sleep(backoff) // apply exponential backoff timer
		}

		bar.Increment()
	}
}

// createDialer returns an instance of [net.Dialer] with a custom [net.Dialer.Timeout] set.
func createDialer(seconds int) net.Dialer {
	return net.Dialer{
		Timeout: time.Duration(seconds) * time.Second,
	}
}
