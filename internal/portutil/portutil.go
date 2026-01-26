package portutil

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// FindAvailable finds an available port starting from the given address.
// It returns the available address and a list of ports that were tried but in use.
// The addr parameter should be in the format ":port" (e.g., ":3000").
func FindAvailable(addr string) (string, []string) {
	var triedPorts []string

	port := 3000
	if strings.HasPrefix(addr, ":") {
		if p, err := strconv.Atoi(addr[1:]); err == nil {
			port = p
		}
	}

	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		testAddr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", testAddr)
		if err == nil {
			_ = ln.Close()
			return testAddr, triedPorts
		}
		triedPorts = append(triedPorts, strconv.Itoa(port))
		port++
	}

	return fmt.Sprintf(":%d", port), triedPorts
}

// FindAvailablePort is a convenience function that takes an int port
// and returns the available port number.
func FindAvailablePort(startPort int) (int, []int) {
	var triedPorts []int

	port := startPort
	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		testAddr := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", testAddr)
		if err == nil {
			_ = ln.Close()
			return port, triedPorts
		}
		triedPorts = append(triedPorts, port)
		port++
	}

	return port, triedPorts
}
