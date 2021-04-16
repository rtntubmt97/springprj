package utils

import (
	"net"
	"strconv"
)

func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return false
	}

	err = ln.Close()
	if err != nil {
		return false
	}

	return true
}

func GetAvailablePort(start int, end int) int {
	for port := start; port < end; port++ {
		if IsPortAvailable(port) {
			return port
		}
	}

	return -1
}
