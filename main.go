// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer) {
	defer wg.Done()
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

func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	target := "scanme.nmap.org"

	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	workers := 100

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, dialer)
	}

	ports := 512

	for p := 1; p <= ports; p++ {
		port := strconv.Itoa(p)
		address := net.JoinHostPort(target, port)
		tasks <- address
	}
	close(tasks)
	wg.Wait()
}
