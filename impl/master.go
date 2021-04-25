// In this project, master only send the signal to observer and nodes to do action.
// Master can send Kill signal to a node, an observer, or all nodes, to force them exit.
// Master can send Send signal to a node with reciver and money data (then that node
// will send the money to the receiver).
// Master can send Receive signal to a node to make it take money in the money channel
// and add that them to its current money.
// Similarly, master can send the ReceiveAll, BeginSnapShot, CollectState, PrintSnapshot
// signal to observer/nodes to make them do their own business.

// Master package includes 3 files: master.go, masterCallers.go and masterHandlers.go.
// Master.go file contains the master structure and its basic operation.
// masterCallers.go file contains its caller which will be used to send the signal
// to nodes/observer by sending specific message.
// masterHandlers.go file contains its handler which will be used to handle the incoming messages.

package impl

import (
	"fmt"
)

// Master struct, it contain an id (which will be assign to an defined int32) and an connector
type Master struct {
	id        int32
	connector Connector
}

// Initilize master
func (master *Master) Init() {
	master.id = MasterId
	master.connector = Connector{}
	master.connector.Init(MasterId)
	master.connector.ParticipantType = MasterType

	master.connector.SetHandleFunc(RequestInfo, master.requestInfoHandle)
}

// Start the listen operation of master
func (master *Master) Listen() {
	master.connector.Listen(int(MasterPort))
}

// Connect to a connector by id and port
func (master *Master) Connect(id int32, port int32) {
	master.connector.Connect(id, port)
}

// Send Kill signal to all connector it known
func (master *Master) KillAll() {
	for nodeId := range master.connector.ConnectedConns {
		master.SignalKill(nodeId)
	}

}

// Send Kill signal to a node, it can be an observer
func (master *Master) SignalKill(nodeId int32) {
	conn := master.connector.ConnectedConns[nodeId]
	if conn == nil {
		LogE("nil conn")
		return
	}
	msg := BinaryProtocol{}
	msg.Init(Input_Kill)
	master.connector.WriteTo(nodeId, &msg)
	LogI(fmt.Sprintf("Sent kill to nodeId %d", nodeId))

	master.connector.WaitAckRsp(nodeId, Input_KillRsp)
}

// Send Send signal to a node
func (master *Master) SignalSend(nodeId int32, receiver int32, money int32) {
	conn := master.connector.ConnectedConns[nodeId]
	if conn == nil {
		LogE("nil conn")
		return
	}
	msg := BinaryProtocol{}
	msg.Init(Input_Send)
	msg.WriteI32(receiver)
	msg.WriteI32(money)
	master.connector.WriteTo(nodeId, &msg)
	LogI(fmt.Sprintf("Sent inputSend to nodeId %d, receiver is %d, money is %d", nodeId, receiver, money))

	master.connector.WaitAckRsp(nodeId, Input_SendRsp)
}

// Send Receive signal to a node
func (master *Master) SignalReceive(receiver int32, sender int32) {
	conn := master.connector.ConnectedConns[receiver]
	if conn == nil {
		LogE("nil conn")
		return
	}
	msg := BinaryProtocol{}
	msg.Init(Input_receive)
	msg.WriteI32(sender)
	master.connector.WriteTo(receiver, &msg)
	LogI(fmt.Sprintf("Sent inputReceive to nodeId %d, sender is %d", receiver, sender))

	master.connector.WaitAckRsp(receiver, Input_receiveRsp)
}

// Send ReceiveAll signal to a node
func (master *Master) SignalReceiveAll() {
	for connId := range master.connector.ConnectedConns {
		if connId == MasterId || connId == ObserverId {
			continue
		}
		msg := BinaryProtocol{}
		msg.Init(Input_receiveAll)
		master.connector.WriteTo(connId, &msg)
		LogI(fmt.Sprintf("Sent inputReceiveAll to nodeId %d", connId))

		master.connector.WaitAckRsp(connId, Input_receiveAllRsp)
	}
}

// Send BeginSnapshot signal to a node
func (master *Master) SignalBeginSnapshot(startNodeId int32) {
	conn := master.connector.ConnectedConns[startNodeId]
	if conn == nil {
		LogE("nil conn")
		return
	}
	msg := BinaryProtocol{}
	msg.Init(Input_BeginSnapshot)
	master.connector.WriteTo(startNodeId, &msg)
	LogI(fmt.Sprintf("Sent InputBeginSnapshot to nodeId %d", startNodeId))

	master.connector.WaitAckRsp(startNodeId, Input_BeginSnapshotRsp)
}

// Send PrintSnapshot signal to the observer
func (master *Master) SignalPrintSnapshot() {
	msg := BinaryProtocol{}
	msg.Init(Input_PrintSnapshot)
	master.connector.WriteTo(ObserverId, &msg)
	LogI(fmt.Sprintf("Sent InputPrintSnapshot to nodeId %d", ObserverId))

	master.connector.WaitAckRsp(ObserverId, Input_PrintSnapshotRsp)
}

// Send CollectState signal to the observer
func (master *Master) SignalCollectState() {
	msg := BinaryProtocol{}
	msg.Init(Input_CollectState)
	master.connector.WriteTo(ObserverId, &msg)
	LogI(fmt.Sprintf("Sent InputCollectState to nodeId %d", ObserverId))

	master.connector.WaitAckRsp(ObserverId, Input_CollectStateRsp)
}

// Handle the request info message from a connector, response it all master-known
// connector information
func (master *Master) requestInfoHandle(connId int32, msg MessageBuffer) {
	LogI("requestInfo_whandle run")
	rspMsg := BinaryProtocol{}
	rspMsg.Init(Rsp)
	rspMsg.WriteI32(int32(RequestInfoRsp))
	rspMsg.WriteI32(int32(len(master.connector.ConnectedConns)))
	for otherConnId := range master.connector.ConnectedConns {
		rspMsg.WriteI32(otherConnId)
		port := master.connector.OtherInfos[otherConnId].ListenPort
		rspMsg.WriteI32(int32(port))
	}
	conn := master.connector.ConnectedConns[connId]
	rspMsg.Write(conn)
}
