package connection

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/andreshungbz/mpscan/internal/scan"
)

func CreateDialer(seconds int) net.Dialer {
	return net.Dialer{
		Timeout: time.Duration(seconds) * time.Second,
	}
}

func CreateWorker(tasks chan scan.Address, dialer net.Dialer, summary *scan.Summary) {
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
				summary.AddPort(address.Port)
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
