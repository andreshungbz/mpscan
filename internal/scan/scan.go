// Package scan contains structs for operating on scans and their results.
package scan

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Address represents a simple pairing of hostname and port.
type Address struct {
	Hostname string
	Port     int
}

// Flags represents the used parameters when conducting a scan.
type Flags struct {
	Target    string
	StartPort int
	EndPort   int
	Workers   int
	Timeout   int

	Ports PortList
}

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
		"[%s]\nTotal Ports Scanned: %d\nOpen Ports Count: %d\nOpen Ports: %v\nTime Taken: %.3fs",
		s.Hostname, s.TotalPortsScanned, s.OpenPortCount, s.OpenPorts, s.TimeTaken.Seconds(),
	)
}

// AddPort appends a port and sorts the slice so that it is in ascending order.
func (s *Summary) AddPort(port int) {
	s.OpenPorts = append(s.OpenPorts, port)
	sort.Ints(s.OpenPorts)
	s.OpenPortCount = len(s.OpenPorts)
}

// CUSTOM FLAGS

// PortsList is a custom flag type for specifying multiple particular ports.
type PortList []int

// Set is the method called by [flag.Var] to parsing values.
func (pl *PortList) Set(value string) error {
	ports := strings.Split(value, ",")
	for _, port := range ports {
		var portNum int
		_, err := fmt.Sscanf(port, "%d", &portNum)
		if err != nil || portNum < 1 || portNum > 65535 {
			continue
		}
		*pl = append(*pl, portNum)
	}
	return nil
}

func (p *PortList) String() string {
	strPorts := make([]string, len(*p))
	for i, port := range *p {
		strPorts[i] = fmt.Sprintf("%d", port)
	}
	return strings.Join(strPorts, ",")
}

// TargetList is a custom flag type for specifying multiple targets.
type TargetList []string

// Set is the method called by [flag.Var] to parsing values.
func (t *TargetList) Set(value string) error {
	*t = strings.Split(value, ",")
	return nil
}

// Print format of a [TargetList].
func (t *TargetList) String() string {
	return strings.Join(*t, ",")
}
