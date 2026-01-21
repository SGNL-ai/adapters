// Copyright 2026 SGNL.ai, Inc.

package common

import (
	"errors"
	"net"
	"testing"
)

func availablePort() (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("type assertion to `*net.TCPAddr` failed")
	}

	return tcpAddr.Port, nil
}

func AvailableTestPort(t *testing.T) int {
	p, err := availablePort()
	if err != nil {
		t.Error("Failed to find available port", err)

		return 0
	}

	return p
}
