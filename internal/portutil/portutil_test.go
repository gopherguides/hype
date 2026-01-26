package portutil

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindAvailable(t *testing.T) {
	r := require.New(t)

	addr, _ := FindAvailable(":3000")
	r.NotEmpty(addr)
	r.True(len(addr) > 1)

	ln, err := net.Listen("tcp", addr)
	r.NoError(err)
	_ = ln.Close()
}

func TestFindAvailable_PortInUse(t *testing.T) {
	r := require.New(t)

	ln, err := net.Listen("tcp", ":0")
	r.NoError(err)
	defer func() { _ = ln.Close() }()

	port := ln.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf(":%d", port)

	newAddr, tried := FindAvailable(addr)
	r.NotEmpty(newAddr)
	r.NotEqual(addr, newAddr)
	r.Contains(tried, fmt.Sprintf("%d", port))
}

func TestFindAvailablePort(t *testing.T) {
	r := require.New(t)

	port, _ := FindAvailablePort(3000)
	r.Greater(port, 0)

	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	r.NoError(err)
	_ = ln.Close()
}

func TestFindAvailablePort_PortInUse(t *testing.T) {
	r := require.New(t)

	ln, err := net.Listen("tcp", ":0")
	r.NoError(err)
	defer func() { _ = ln.Close() }()

	usedPort := ln.Addr().(*net.TCPAddr).Port

	port, tried := FindAvailablePort(usedPort)
	r.NotEqual(usedPort, port)
	r.Contains(tried, usedPort)
}
