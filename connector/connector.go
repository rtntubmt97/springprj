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
	handlers       map[define.ConnectorCmd]define.HandleFunc
}

func (connector *Connector) Init(id int32) {
	connector.id = id
	connector.connectedConns = make(map[int32]net.Conn)
	connector.handlers = make(map[define.ConnectorCmd]define.HandleFunc)

	utils.LogI(fmt.Sprintf("connId %d initiated", id))
}

func (connector *Connector) GetConnection(id int32) net.Conn {
	return connector.connectedConns[id]
}

func (connector *Connector) SetHandleFunc(cmd define.ConnectorCmd, f define.HandleFunc) {
	connector.handlers[cmd] = f
}

func (connector *Connector) Connect(id int32, port int32) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", add)
	if err != nil {
		utils.LogE("Invalid connect port")
		return
	}

	connId := connector.greeting_wcall(conn)

	if _, exist := connector.connectedConns[connId]; exist {
		utils.LogE(fmt.Sprintf("connId %d existed", connId))
		return
	}
	if connId == -1 {
		utils.LogE("Invalid message")
		return
	}
	if connId != id {
		utils.LogE(fmt.Sprintf("Invalid connId %d", connId))
		return
	}

	utils.LogI(fmt.Sprintf("Connected connId %d", connId))
	connector.connectedConns[connId] = conn

	go connector.Handle(connId, conn)
}

func (connector *Connector) Listen(port int) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	connector.listener, err = net.Listen("tcp", add)
	if err != nil {
		utils.LogE("Invalid listen port")
		return
	}

	utils.LogI(fmt.Sprintf("Start listening on port %d", port))
	for {
		conn, err := connector.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// msg := protocol.MessageBuffer{}
		msg := protocol.SimpleMessageBuffer{}
		msg.ReadMessage(conn)

		connId := connector.greeting_whandle(msg, conn)
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

		go connector.Handle(connId, conn)
	}
}

func (connector *Connector) Handle(connId int32, conn net.Conn) {
	for {
		// utils.LogI(fmt.Sprintf("%d run Handle", connector.id))

		msg := protocol.SimpleMessageBuffer{}
		readErr := msg.ReadMessage(conn)
		if readErr != nil {
			break
		}
		cmd := define.ConnectorCmd(msg.ReadI32())
		f := connector.handlers[cmd]
		if f != nil {
			f(connId, msg)
		}
	}
}

func (connector *Connector) greeting_wcall(conn net.Conn) int32 {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Greeting)
	msg.WriteI32(connector.id)
	msg.WriteMessage(conn)

	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.ReadMessage(conn)

	cmd := define.ConnectorCmd(rspMsg.ReadI32())
	if cmd != define.GreetingRsp {
		return -1
	}

	return rspMsg.ReadI32()
}

func (connector *Connector) greeting_whandle(msg define.MessageBuffer, conn net.Conn) int32 {
	cmd := define.ConnectorCmd(msg.ReadI32())
	if cmd != define.Greeting {
		return -1
	}

	idFromGreeting := msg.ReadI32()

	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Init(define.GreetingRsp)
	rspMsg.WriteI32(connector.id)
	rspMsg.WriteMessage(conn)

	return idFromGreeting
}
