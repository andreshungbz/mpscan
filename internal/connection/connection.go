package connection

import (
	"fmt"
	"net"
	"time"
)

func CreateDialer(seconds int) net.Dialer {
	return net.Dialer{
		Timeout: time.Duration(seconds) * time.Second,
	}
}

func CreateWorker(tasks chan string, dialer net.Dialer) {
	maxRetries := 3
	for addr := range tasks {
		var success bool
		for i := range maxRetries {
			conn, err := dialer.Dial("tcp", addr)
			if err == nil {
				conn.Close()
				fmt.Printf("Connection to %s was successful\n", addr)
				success = true
				break
			}
			backoff := time.Duration(1<<i) * time.Second
			fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1, addr, backoff)
			time.Sleep(backoff)
		}
		if !success {
			fmt.Printf("Failed to connect to %s after %d attempts\n", addr, maxRetries)
		}
	}
}
