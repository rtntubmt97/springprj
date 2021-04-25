// In this project, master, observer and nodes are different processes. They will
// communicate with each other through the network. This package provides basic materials
// which will be used by those processes to connect to other, communicate by sending
// and recieving messages.

// All the data sent by connector will be serilized by the a message buffer with an cmd (int32) in front
// of it in order to the reciever know what to do with that message.
// The receiver will use the cmd in front of each message to determine the handler to handle it.

package impl

import (
	"fmt"
	"net"
	"sync"
)

// Participant Type, this data show role of an process in the network (master, observer, node).
type ParticipantType int32

// Integer value of (master, observer, node) type
const (
	MasterType   ParticipantType = 1
	ObserverType ParticipantType = 2
	NodeType     ParticipantType = 3
)

// Basic information of a proccess, it contains role, id, and listenning port for connection.
type ParticipantInfo struct {
	ParticipantType
	NodeId     int32
	ListenPort int32
}

// Connector is the struct which (master, observer, node) will be used to connect others.
type Connector struct {
	ParticipantType
	id              int32
	listenPort      int32
	listener        net.Listener
	ConnectedConns  map[int32]net.Conn
	handlers        map[ConnectorCmd]HandleFunc
	writeOutMutexes map[int32]*sync.Mutex
	rspMsgChannel   map[int32](chan MessageBuffer)
	OtherInfos      map[int32]*ParticipantInfo
	readyMutex      *sync.Mutex
	afterAccept     AfterAccept
}

// The interface of a call back function, this function will be called after an connection
// is accepted by the connector.
type AfterAccept func(connInfo ParticipantInfo)

// Initilize the connector
func (connector *Connector) Init(id int32) {
	connector.id = id
	connector.ConnectedConns = make(map[int32]net.Conn)
	connector.handlers = make(map[ConnectorCmd]HandleFunc)
	connector.writeOutMutexes = make(map[int32]*sync.Mutex)
	connector.rspMsgChannel = make(map[int32]chan MessageBuffer)
	connector.OtherInfos = make(map[int32]*ParticipantInfo)

	LogI(fmt.Sprintf("connId %d initiated", id))
	connector.readyMutex = new(sync.Mutex)

	connector.SetHandleFunc(Rsp, connector.forwardRsp)
}

// The connector communicates by sending and recieving messages. This function set the handler
// for each kind of messages. Pay attention to the use of this function in the node or
// master structure to understand the concept of it.
func (connector *Connector) SetHandleFunc(cmd ConnectorCmd, f HandleFunc) {
	connector.handlers[cmd] = f
}

// Set the AfterAccept function for connector.
func (connector *Connector) SetAfterAccept(afterAccept AfterAccept) {
	connector.afterAccept = afterAccept
}

// Wait the Connector for availble by trying to lock and unlock the mutex.
func (connector *Connector) WaitReady() {
	connector.readyMutex.Lock()
	connector.readyMutex.Unlock()
}

// Check if the connector has connected to other connector by the other's id.
func (connector *Connector) IsConnected(otherId int32) bool {
	_, ok := connector.ConnectedConns[otherId]
	return ok
}

// Initialize and save new connector information, prepare for sending and recieving data
// from it
func (connector *Connector) initNewConn(peerInfo ParticipantInfo, conn net.Conn) {
	connId := peerInfo.NodeId
	connector.ConnectedConns[connId] = conn
	connector.writeOutMutexes[connId] = new(sync.Mutex)
	connector.rspMsgChannel[connId] = make(chan MessageBuffer)
	connector.OtherInfos[connId] = &peerInfo
}

// Connect to other connector by its id and listen port
func (connector *Connector) Connect(id int32, port int32) {
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", add)
	if err != nil {
		LogE("Invalid connect port")
		return
	}

	err, info := connector.greeting_wcall(conn)
	connId := info.NodeId

	if err != nil {
		LogE(err.Error())
		return
	}
	if _, exist := connector.ConnectedConns[connId]; exist {
		LogE(fmt.Sprintf("connId %d existed", connId))
		return
	}
	if connId != id {
		LogE(fmt.Sprintf("Invalid connId %d", connId))
		return
	}
	otherInfo := ParticipantInfo{info.ParticipantType, connId, port}

	LogI(fmt.Sprintf("Connector %d connected conn %d", connector.id, otherInfo))
	connector.initNewConn(otherInfo, conn)

	go connector.Handle(otherInfo, conn)
}

// Start listen at port for other connector to connector to connect
func (connector *Connector) Listen(port int) {
	connector.readyMutex.Lock()
	var err error
	add := fmt.Sprintf("localhost:%d", port)
	connector.listener, err = net.Listen("tcp", add)
	if err != nil {
		LogE(fmt.Sprintf("Connector %d got invalid listen port %d", connector.id, port))
		return
	}

	connector.listenPort = int32(port)
	connector.readyMutex.Unlock()

	LogI(fmt.Sprintf("Start listening on port %d", port))
	for {
		conn, err := connector.listener.Accept()
		if err != nil {
			LogE(err.Error())
			continue
		}

		msg := BinaryProtocol{}
		msg.Read(conn)

		err, otherInfo := connector.greeting_whandle(msg, conn)
		if err != nil {
			LogE(err.Error())
			continue
		}
		if _, exist := connector.ConnectedConns[otherInfo.NodeId]; exist {
			LogE(fmt.Sprintf("connId %d existed", otherInfo))
			continue
		}

		LogI(fmt.Sprintf("Connector %d accepted conn %d", connector.id, otherInfo))
		connector.initNewConn(otherInfo, conn)
		if connector.afterAccept != nil {
			connector.afterAccept(otherInfo)
		}

		go connector.Handle(otherInfo, conn)
	}
}

// Handle incoming connector by its info and net.Conn
func (connector *Connector) Handle(otherInfo ParticipantInfo, conn net.Conn) {
	for {
		// LogI(fmt.Sprintf("%d run Handle", connector.id))

		msg := BinaryProtocol{}
		readErr := msg.Read(conn)
		if readErr != nil {
			break
		}
		cmd := ConnectorCmd(msg.ReadI32())
		f := connector.handlers[cmd]
		if f != nil {
			f(otherInfo.NodeId, msg)
		}
	}
}

// Send current connector's information to other connector and recieve its information
func (connector *Connector) greeting_wcall(conn net.Conn) (error, ParticipantInfo) {
	msg := BinaryProtocol{}
	msg.Init(Greeting)
	msg.WriteI32(int32(connector.ParticipantType))
	msg.WriteI32(connector.id)
	msg.WriteI32(connector.listenPort)
	msg.Write(conn)

	rspMsg := BinaryProtocol{}
	rspMsg.Read(conn)

	cmd := ConnectorCmd(rspMsg.ReadI32())
	if cmd != GreetingRsp {
		return ErrFailGreeting, ParticipantInfo{}
	}

	cType := rspMsg.ReadI32()
	cId := rspMsg.ReadI32()

	return nil, ParticipantInfo{ParticipantType(cType), cId, 0}
}

// Recieve accepted connector's information and send current connector information to that connector
func (connector *Connector) greeting_whandle(msg MessageBuffer, conn net.Conn) (error, ParticipantInfo) {
	cmd := ConnectorCmd(msg.ReadI32())
	if cmd != Greeting {
		return ErrWrongCmd, ParticipantInfo{}
	}

	typeFromGreeting := ParticipantType(msg.ReadI32())
	idFromGreeting := msg.ReadI32()
	portFromGreeting := msg.ReadI32()

	rspMsg := BinaryProtocol{}
	rspMsg.Init(GreetingRsp)
	rspMsg.WriteI32(int32(connector.ParticipantType))
	rspMsg.WriteI32(connector.id)
	rspMsg.Write(conn)

	return nil, ParticipantInfo{
		ParticipantType: typeFromGreeting,
		NodeId:          idFromGreeting,
		ListenPort:      portFromGreeting}
}

// Write a writeable message to a connector specified by its id
func (connector *Connector) WriteTo(connId int32, msg Writeable) {
	mutex := connector.writeOutMutexes[connId]
	mutex.Lock()
	defer mutex.Unlock()
	conn := connector.ConnectedConns[connId]
	msg.Write(conn)
}

// Forward the response message to its corresponding channel in order to the handler
// receive and continue proccess it
func (connector *Connector) forwardRsp(connId int32, msg MessageBuffer) {
	connector.rspMsgChannel[connId] <- msg
}

// Wait an incoming message from a specific connector
func (connector *Connector) WaitRsp(connId int32) MessageBuffer {
	return <-connector.rspMsgChannel[connId]
}

// Send an ack message to a specific connector
func (connector *Connector) SendAckRsp(connId int32, cmd ConnectorCmd) {
	rspMsg := BinaryProtocol{}
	rspMsg.Init(Rsp)
	rspMsg.WriteI32(int32(cmd))
	connector.WriteTo(connId, &rspMsg)
}

// Wait fo an Ack message from a specific connector and check where it is correct ack by the cmd
func (connector *Connector) WaitAckRsp(nodeId int32, cmd ConnectorCmd) {
	rspMsg := connector.WaitRsp(nodeId)
	rspCmd := ConnectorCmd(rspMsg.ReadI32())
	if cmd == rspCmd {
		LogI(fmt.Sprintf("Connecter %d received correct response for cmd %d", connector.id, cmd))
	} else {
		LogI(fmt.Sprintf("Connecter %d  received %d, wrong response for cmd %d", connector.id, rspCmd, cmd))
	}
}
