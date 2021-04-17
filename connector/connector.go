package connector

import (
	"fmt"
	"net"
	"sync"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

type OtherInfo struct {
	NodeId     int32
	ListenPort int32
}

type Connector struct {
	id              int32
	listenPort      int32
	listener        net.Listener
	ConnectedConns  map[int32]net.Conn
	handlers        map[define.ConnectorCmd]define.HandleFunc
	writeOutMutexes map[int32]*sync.Mutex
	rspMsgChannel   map[int32](chan define.MessageBuffer)
	OtherInfos      map[int32]*OtherInfo
	readyMutex      sync.Mutex
}

func (connector *Connector) Init(id int32) {
	connector.id = id
	connector.ConnectedConns = make(map[int32]net.Conn)
	connector.handlers = make(map[define.ConnectorCmd]define.HandleFunc)
	connector.writeOutMutexes = make(map[int32]*sync.Mutex)
	connector.rspMsgChannel = make(map[int32]chan define.MessageBuffer)
	connector.OtherInfos = make(map[int32]*OtherInfo)

	utils.LogI(fmt.Sprintf("connId %d initiated", id))

	connector.SetHandleFunc(define.Rsp, connector.forwardRsp)
}

// func (connector *Connector) GetConnection(id int32) net.Conn {
// 	return connector.ConnectedConns[id]
// }

func (connector *Connector) SetHandleFunc(cmd define.ConnectorCmd, f define.HandleFunc) {
	connector.handlers[cmd] = f
}

func (connector *Connector) WaitReady() {
	connector.readyMutex.Lock()
	connector.readyMutex.Unlock()
}

func (connector *Connector) IsConnected(otherId int32) bool {
	_, ok := connector.ConnectedConns[otherId]
	return ok
}

func (connector *Connector) initNewConn(otherInfo OtherInfo, conn net.Conn) {
	connId := otherInfo.NodeId
	connector.ConnectedConns[connId] = conn
	connector.writeOutMutexes[connId] = new(sync.Mutex)
	connector.rspMsgChannel[connId] = make(chan define.MessageBuffer)
	connector.OtherInfos[connId] = &otherInfo
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

	if _, exist := connector.ConnectedConns[connId]; exist {
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
	otherInfo := OtherInfo{NodeId: id, ListenPort: port}

	utils.LogI(fmt.Sprintf("Connector %d connected conn %d", connector.id, otherInfo))
	connector.initNewConn(otherInfo, conn)

	go connector.Handle(otherInfo, conn)
}

func (connector *Connector) Listen(port int) {
	connector.readyMutex.Lock()
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	connector.listener, err = net.Listen("tcp", add)
	if err != nil {
		utils.LogE(fmt.Sprintf("Connector %d got invalid listen port %d", connector.id, port))
		return
	}

	connector.listenPort = int32(port)
	connector.readyMutex.Unlock()

	utils.LogI(fmt.Sprintf("Start listening on port %d", port))
	for {
		conn, err := connector.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// msg := protocol.MessageBuffer{}
		msg := protocol.SimpleMessageBuffer{}
		msg.Read(conn)

		err, otherInfo := connector.greeting_whandle(msg, conn)
		if err != nil {
			utils.LogE(err.Error())
			continue
		}
		if _, exist := connector.ConnectedConns[otherInfo.NodeId]; exist {
			utils.LogE(fmt.Sprintf("connId %d existed", otherInfo))
			continue
		}
		// if otherInfo == -1 {
		// 	utils.LogE("Invalid message")
		// 	continue
		// }

		utils.LogI(fmt.Sprintf("Connector %d accepted conn %d", connector.id, otherInfo))
		connector.initNewConn(otherInfo, conn)

		go connector.Handle(otherInfo, conn)
	}
}

func (connector *Connector) Handle(otherInfo OtherInfo, conn net.Conn) {
	for {
		// utils.LogI(fmt.Sprintf("%d run Handle", connector.id))

		msg := protocol.SimpleMessageBuffer{}
		readErr := msg.Read(conn)
		if readErr != nil {
			break
		}
		cmd := define.ConnectorCmd(msg.ReadI32())
		f := connector.handlers[cmd]
		if f != nil {
			f(otherInfo.NodeId, msg)
		}
	}
}

func (connector *Connector) greeting_wcall(conn net.Conn) int32 {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Greeting)
	msg.WriteI32(connector.id)
	msg.WriteI32(connector.listenPort)
	msg.Write(conn)

	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Read(conn)

	cmd := define.ConnectorCmd(rspMsg.ReadI32())
	if cmd != define.GreetingRsp {
		return -1
	}

	return rspMsg.ReadI32()
}

func (connector *Connector) greeting_whandle(msg define.MessageBuffer, conn net.Conn) (error, OtherInfo) {
	cmd := define.ConnectorCmd(msg.ReadI32())
	if cmd != define.Greeting {
		return define.ErrWrongCmd, OtherInfo{}
	}

	idFromGreeting := msg.ReadI32()
	portFromGreeting := msg.ReadI32()

	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Init(define.GreetingRsp)
	rspMsg.WriteI32(connector.id)
	rspMsg.Write(conn)

	return nil, OtherInfo{NodeId: idFromGreeting, ListenPort: portFromGreeting}
}

func (connector *Connector) WriteTo(connId int32, msg define.Writeable) {
	mutex := connector.writeOutMutexes[connId]
	mutex.Lock()
	defer mutex.Unlock()
	conn := connector.ConnectedConns[connId]
	msg.Write(conn)
}

func (connector *Connector) forwardRsp(connId int32, msg define.MessageBuffer) {
	connector.rspMsgChannel[connId] <- msg
}

func (connector *Connector) WaitRsp(connId int32) define.MessageBuffer {
	return <-connector.rspMsgChannel[connId]
}
