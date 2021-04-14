package connector

import (
	"fmt"
	"net"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
)

type Connector struct {
	id            int32
	listener      net.Listener
	acceptedConns map[int32]net.Conn
	handlers      map[int32]define.HandleFunc
}

func (connector *Connector) Init() {
	connector.acceptedConns = make(map[int32]net.Conn)
	connector.handlers = make(map[int32]define.HandleFunc)
}

func (connector *Connector) SetHandleFunc(cmd int32, f define.HandleFunc) {
	connector.handlers[cmd] = f
}

func (connector *Connector) Connect(port int32) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", add)
	if err != nil {
		return
	}

	msg := protocol.ReadMessage(conn)
	isValidConn := connector.greetingBack_handle(*msg, conn)
	if isValidConn {
		go connector.Handle(conn)
	}
}

func (connector *Connector) Listen(port int) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	connector.listener, err = net.Listen("tcp", add)
	if err != nil {
		return
	}

	for {
		conn, err := connector.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		msg := protocol.MessageBuffer{}
		isValidConn := connector.greeting_handle(msg, conn)
		if isValidConn {
			go connector.Handle(conn)
		}
	}
}

func (connector *Connector) Handle(conn net.Conn) {
	for {
		msg := protocol.ReadMessage(conn)
		if msg == nil {
			break
		}
		cmd := msg.ReadI32()
		f := connector.handlers[cmd]
		if f != nil {
			f(*msg)
		}
	}
}

func (connector *Connector) greeting_handle(msg protocol.MessageBuffer, conn net.Conn) bool {
	cmd := msg.ReadI32()
	if cmd != define.Greeting {
		return false
	}

	id := msg.ReadI32()
	connector.acceptedConns[id] = conn

	sendMsg := protocol.MessageBuffer{}
	sendMsg.InitEmpty()
	sendMsg.WriteI32(define.GreetingBack)
	sendMsg.WriteI32(connector.id)
	protocol.WriteMessage(conn, sendMsg)

	return true
}

func (connector *Connector) greetingBack_handle(msg protocol.MessageBuffer, conn net.Conn) bool {
	cmd := msg.ReadI32()
	if cmd != define.GreetingBack {
		return false
	}

	id := msg.ReadI32()
	connector.acceptedConns[id] = conn

	return true
}
