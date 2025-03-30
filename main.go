// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

import (
	"net"
	"strconv"
	"sync"

	"github.com/andreshungbz/mpscan/internal/connection"
)

func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	target := "scanme.nmap.org"

	dialer := connection.CreateDialer(5)

	workers := 100

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connection.CreateWorker(tasks, dialer)
		}()
	}

	ports := 100

	for p := 1; p <= ports; p++ {
		port := strconv.Itoa(p)
		address := net.JoinHostPort(target, port)
		tasks <- address
	}
	close(tasks)
	wg.Wait()
}
