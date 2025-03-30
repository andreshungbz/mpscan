// Package scan contains structs for operating on scans and their results.
package scan

import (
	"fmt"
	"sort"
	"time"
)

// Summary represents the statistics of a particular scan.
type Summary struct {
	Hostname          string
	TotalPortsScanned int
	OpenPortCount     int
	OpenPorts         []int
	TimeTaken         time.Duration
}

// Print format of a [Summary].
func (s Summary) String() string {
	return fmt.Sprintf(
		"[summary/%s]\nTotal Ports Scanned: %d\nOpen Ports Count: %d\nOpen Ports: %v\nTime Taken: %.3fs",
		s.Hostname, s.TotalPortsScanned, s.OpenPortCount, s.OpenPorts, s.TimeTaken.Seconds(),
	)
}

// AddPort appends a port and sorts the slice so that it is in ascending order.
func (s *Summary) AddPort(port int) {
	s.OpenPorts = append(s.OpenPorts, port)
	sort.Ints(s.OpenPorts)
	s.OpenPortCount = len(s.OpenPorts)
}

// Address represents a simple pairing of hostname and port.
type Address struct {
	Hostname string
	Port     int
}

// Flags represents the used paramters when conducting a scan.
type Flags struct {
	Target    *string
	StartPort *int
	EndPort   *int
	Workers   *int
	Timeout   *int
}
