package scan

import (
	"fmt"
	"sort"
	"time"
)

type Summary struct {
	Hostname          string
	TotalPortsScanned int
	OpenPortCount     int
	OpenPorts         []int
	TimeTaken         time.Duration
}

func (s Summary) String() string {
	return fmt.Sprintf(
		"[summary/%s]\nTotal Ports Scanned: %d\nOpen Ports Count: %d\nOpen Ports: %v\nTime Taken: %s",
		s.Hostname, s.TotalPortsScanned, s.OpenPortCount, s.OpenPorts, s.TimeTaken,
	)
}

func (s *Summary) AddPort(port int) {
	s.OpenPorts = append(s.OpenPorts, port)
	sort.Ints(s.OpenPorts)
	s.OpenPortCount = len(s.OpenPorts)
}

type Address struct {
	Hostname string
	Port     int
}

type Flags struct {
	Target    *string
	StartPort *int
	EndPort   *int
	Workers   *int
	Timeout   *int
}
