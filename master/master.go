package master

import (
	"bufio"
	"fmt"
	"net"
)

type Master struct{}

func (master *Master) Start() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn)
	}
}
func handleUserConnection(c net.Conn) {
	fmt.Println("foo")
	defer c.Close()

	reader := bufio.NewReader(c)
	for {
		fmt.Println("foobar")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Println(userInput)
	}
}
