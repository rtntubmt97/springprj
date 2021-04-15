package connector

import (
	"fmt"
	"net"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

type Connector struct {
	id             int32
	listener       net.Listener
	connectedConns map[int32]net.Conn
	handlers       map[int32]define.HandleFunc
}

func (connector *Connector) Init(id int32) {
	connector.id = id
	connector.connectedConns = make(map[int32]net.Conn)
	connector.handlers = make(map[int32]define.HandleFunc)
}

func (connector *Connector) GetConnection(id int32) net.Conn {
	return connector.connectedConns[id]
}

func (connector *Connector) SetHandleFunc(cmd int32, f define.HandleFunc) {
	connector.handlers[cmd] = f
}

func (connector *Connector) Connect(id int32, port int32) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", add)
	if err != nil {
		utils.LogE("Invalid port")
		return
	}

	connector.greeting_call(conn)
	msg := protocol.ReadMessage(conn)

	connId := connector.greetingBack_handle(*msg, conn)
	if _, exist := connector.connectedConns[connId]; exist {
		utils.LogE(fmt.Sprintf("connId %d existed", connId))
		return
	}
	if connId == -1 {
		utils.LogE("Invalid message")
		return
	}
	if connId != id {
		utils.LogE("Invalid connId")
		return
	}

	utils.LogI(fmt.Sprintf("Connected connId %d", connId))
	connector.connectedConns[connId] = conn

	go connector.Handle(conn)
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

		// msg := protocol.MessageBuffer{}
		msg := protocol.ReadMessage(conn)

		connId := connector.greeting_handle(*msg, conn)
		if _, exist := connector.connectedConns[connId]; exist {
			utils.LogE(fmt.Sprintf("connId %d existed", connId))
			continue
		}
		if connId == -1 {
			utils.LogE("Invalid message")
			continue
		}

		utils.LogI(fmt.Sprintf("Accepted connId %d", connId))
		connector.connectedConns[connId] = conn

		go connector.Handle(conn)
	}
}

func (connector *Connector) Handle(conn net.Conn) {
	for {
		// utils.LogI(fmt.Sprintf("%d run Handle", connector.id))
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

func (connector *Connector) greeting_call(conn net.Conn) {
	msg := protocol.MessageBuffer{}
	msg.InitEmpty()
	msg.WriteI32(define.Greeting)
	msg.WriteI32(connector.id)
	protocol.WriteMessage(conn, msg)
}

func (connector *Connector) greetingBack_call(conn net.Conn) {
	msg := protocol.MessageBuffer{}
	msg.InitEmpty()
	msg.WriteI32(define.GreetingBack)
	msg.WriteI32(connector.id)
	protocol.WriteMessage(conn, msg)
}

func (connector *Connector) greeting_handle(msg protocol.MessageBuffer, conn net.Conn) int32 {
	cmd := msg.ReadI32()
	if cmd != define.Greeting {
		return -1
	}

	id := msg.ReadI32()

	connector.greetingBack_call(conn)

	return id
}

func (connector *Connector) greetingBack_handle(msg protocol.MessageBuffer, conn net.Conn) int32 {
	cmd := msg.ReadI32()
	if cmd != define.GreetingBack {
		return -1
	}

	id := msg.ReadI32()

	return id
}
