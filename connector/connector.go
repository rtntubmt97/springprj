package connector

import (
	"fmt"
	"net"
	"sync"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

type ParticipantType int32

const (
	MasterType   ParticipantType = 1
	ObserverType ParticipantType = 2
	NodeType     ParticipantType = 3
)

type ParticipantInfo struct {
	ParticipantType
	NodeId     int32
	ListenPort int32
}

type Connector struct {
	ParticipantType
	id              int32
	listenPort      int32
	listener        net.Listener
	ConnectedConns  map[int32]net.Conn
	handlers        map[define.ConnectorCmd]define.HandleFunc
	writeOutMutexes map[int32]*sync.Mutex
	rspMsgChannel   map[int32](chan define.MessageBuffer)
	OtherInfos      map[int32]*ParticipantInfo
	readyMutex      *sync.Mutex
	afterAccept     AfterAccept
}

type AfterAccept func(connInfo ParticipantInfo)

func (connector *Connector) Init(id int32) {
	connector.id = id
	connector.ConnectedConns = make(map[int32]net.Conn)
	connector.handlers = make(map[define.ConnectorCmd]define.HandleFunc)
	connector.writeOutMutexes = make(map[int32]*sync.Mutex)
	connector.rspMsgChannel = make(map[int32]chan define.MessageBuffer)
	connector.OtherInfos = make(map[int32]*ParticipantInfo)

	utils.LogI(fmt.Sprintf("connId %d initiated", id))
	connector.readyMutex = new(sync.Mutex)

	connector.SetHandleFunc(define.Rsp, connector.forwardRsp)
}

// func (connector *Connector) GetConnection(id int32) net.Conn {
// 	return connector.ConnectedConns[id]
// }

func (connector *Connector) SetHandleFunc(cmd define.ConnectorCmd, f define.HandleFunc) {
	connector.handlers[cmd] = f
}

func (connector *Connector) SetAfterAccept(afterAccept AfterAccept) {
	connector.afterAccept = afterAccept
}

func (connector *Connector) WaitReady() {
	connector.readyMutex.Lock()
	connector.readyMutex.Unlock()
}

func (connector *Connector) IsConnected(otherId int32) bool {
	_, ok := connector.ConnectedConns[otherId]
	return ok
}

func (connector *Connector) initNewConn(peerInfo ParticipantInfo, conn net.Conn) {
	connId := peerInfo.NodeId
	connector.ConnectedConns[connId] = conn
	connector.writeOutMutexes[connId] = new(sync.Mutex)
	connector.rspMsgChannel[connId] = make(chan define.MessageBuffer)
	connector.OtherInfos[connId] = &peerInfo
}

func (connector *Connector) Connect(id int32, port int32) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", add)
	if err != nil {
		utils.LogE("Invalid connect port")
		return
	}

	err, info := connector.greeting_wcall(conn)
	connId := info.NodeId

	if err != nil {
		utils.LogE(err.Error())
		return
	}
	if _, exist := connector.ConnectedConns[connId]; exist {
		utils.LogE(fmt.Sprintf("connId %d existed", connId))
		return
	}
	if connId != id {
		utils.LogE(fmt.Sprintf("Invalid connId %d", connId))
		return
	}
	otherInfo := ParticipantInfo{info.ParticipantType, connId, port}

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
			utils.LogE(err.Error())
			continue
		}

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

		utils.LogI(fmt.Sprintf("Connector %d accepted conn %d", connector.id, otherInfo))
		connector.initNewConn(otherInfo, conn)
		if connector.afterAccept != nil {
			connector.afterAccept(otherInfo)
		}

		go connector.Handle(otherInfo, conn)
	}
}

func (connector *Connector) Handle(otherInfo ParticipantInfo, conn net.Conn) {
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

func (connector *Connector) greeting_wcall(conn net.Conn) (error, ParticipantInfo) {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Greeting)
	msg.WriteI32(int32(connector.ParticipantType))
	msg.WriteI32(connector.id)
	msg.WriteI32(connector.listenPort)
	msg.Write(conn)

	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Read(conn)

	cmd := define.ConnectorCmd(rspMsg.ReadI32())
	if cmd != define.GreetingRsp {
		return define.ErrFailGreeting, ParticipantInfo{}
	}

	cType := rspMsg.ReadI32()
	cId := rspMsg.ReadI32()

	return nil, ParticipantInfo{ParticipantType(cType), cId, 0}
}

func (connector *Connector) greeting_whandle(msg define.MessageBuffer, conn net.Conn) (error, ParticipantInfo) {
	cmd := define.ConnectorCmd(msg.ReadI32())
	if cmd != define.Greeting {
		return define.ErrWrongCmd, ParticipantInfo{}
	}

	typeFromGreeting := ParticipantType(msg.ReadI32())
	idFromGreeting := msg.ReadI32()
	portFromGreeting := msg.ReadI32()

	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Init(define.GreetingRsp)
	rspMsg.WriteI32(int32(connector.ParticipantType))
	rspMsg.WriteI32(connector.id)
	rspMsg.Write(conn)

	return nil, ParticipantInfo{
		ParticipantType: typeFromGreeting,
		NodeId:          idFromGreeting,
		ListenPort:      portFromGreeting}
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
