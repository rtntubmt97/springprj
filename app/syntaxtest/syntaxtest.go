package main

import (
	"fmt"
	"net"
)

func main() {
	m := make(map[int32]net.Conn)
	m[1] = &net.IPConn{}
	i := m[2]
	if i == nil {
		fmt.Println(i)
	}
}
