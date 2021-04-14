package node

import (
	"fmt"
	"net"

	"github.com/rtntubmt97/springprj/protocol"
)

type Node struct {
	id             int32
	listener       net.Listener
	masterConn     net.Conn
	otherNodesConn map[int32]net.Conn
}

func (node *Node) Init() {
	node.otherNodesConn = make(map[int32]net.Conn)
}

func (node *Node) Start() {
	c, err := net.Dial("tcp", "localhost:9090")
	defer c.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	n, err := c.Write([]byte("foo123\nasf\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n)
	c.Write([]byte("foo12345\n"))

}

func (node *Node) Listen(port int) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	node.listener, err = net.Listen("tcp", add)
	if err != nil {
		return
	}

	for {
		conn, err := node.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		node.otherNodesConn[count] = conn
		msg := protocol.MessageBuffer{}
		msg.InitEmpty()
		msg.WriteString("Listen")
		protocol.WriteMessage(conn, msg)
		go node.Handle(conn)
	}
}

var count = int32(0)

func (node *Node) Handle(conn net.Conn) {
	count++
	for {
		msg := protocol.ReadMessage(conn)
		if msg == nil {
			break
		}
		// cmd := define.Command(msg.ReadI32())
		// switch (cmd){
		// case define.SendInt32 :
		// }
		fmt.Println(msg.ReadString())
	}
}

func (node *Node) ConnectNode(port int32) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", add)
	if err != nil {
		return
	}
	node.otherNodesConn[count] = conn
	msg := protocol.MessageBuffer{}
	msg.InitEmpty()
	msg.WriteString("ConnectNode")
	protocol.WriteMessage(conn, msg)
	go node.Handle(conn)
}

func (node *Node) ConnectMaster() {

}
