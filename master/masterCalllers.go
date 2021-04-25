// The "w" character in the wcall post fix means that call function will wait for an
// ack message or response message with data before return
// The "Input" in the call prefix mean that this caller will be trigger by input of user.
// They will be called in the app/master/master.go when the correspond input matched

package master

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

// Send Kill signal to a node, it can be an observer
func (master *Master) InputKill_wcall(nodeId int32) {
	conn := master.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_Kill)
	master.connector.WriteTo(nodeId, &msg)
	utils.LogI(fmt.Sprintf("Sent kill to nodeId %d", nodeId))

	master.connector.WaitAckRsp(nodeId, define.Input_KillRsp)
}

// Send Send signal to a node
func (master *Master) InputSend_wcall(nodeId int32, receiver int32, money int32) {
	conn := master.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_Send)
	msg.WriteI32(receiver)
	msg.WriteI32(money)
	master.connector.WriteTo(nodeId, &msg)
	utils.LogI(fmt.Sprintf("Sent inputSend to nodeId %d, receiver is %d, money is %d", nodeId, receiver, money))

	master.connector.WaitAckRsp(nodeId, define.Input_SendRsp)
}

// Send Receive signal to a node
func (master *Master) InputReceive_wcall(receiver int32, sender int32) {
	conn := master.connector.ConnectedConns[receiver]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_receive)
	msg.WriteI32(sender)
	master.connector.WriteTo(receiver, &msg)
	utils.LogI(fmt.Sprintf("Sent inputReceive to nodeId %d, sender is %d", receiver, sender))

	master.connector.WaitAckRsp(receiver, define.Input_receiveRsp)
}

// Send ReceiveAll signal to a node
func (master *Master) InputReceiveAll_wcall() {
	for connId := range master.connector.ConnectedConns {
		if connId == define.MasterId || connId == define.ObserverId {
			continue
		}
		msg := protocol.SimpleMessageBuffer{}
		msg.Init(define.Input_receiveAll)
		master.connector.WriteTo(connId, &msg)
		utils.LogI(fmt.Sprintf("Sent inputReceiveAll to nodeId %d", connId))

		master.connector.WaitAckRsp(connId, define.Input_receiveAllRsp)
	}
}

// Send BeginSnapshot signal to a node
func (master *Master) InputBeginSnapshot_wcall(startNodeId int32) {
	conn := master.connector.ConnectedConns[startNodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_BeginSnapshot)
	master.connector.WriteTo(startNodeId, &msg)
	utils.LogI(fmt.Sprintf("Sent InputBeginSnapshot to nodeId %d", startNodeId))

	master.connector.WaitAckRsp(startNodeId, define.Input_BeginSnapshotRsp)
}

// Send PrintSnapshot signal to the observer
func (master *Master) InputPrintSnapshot_wcall() {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_PrintSnapshot)
	master.connector.WriteTo(define.ObserverId, &msg)
	utils.LogI(fmt.Sprintf("Sent InputPrintSnapshot to nodeId %d", define.ObserverId))

	master.connector.WaitAckRsp(define.ObserverId, define.Input_PrintSnapshotRsp)
}

// Send CollectState signal to the observer
func (master *Master) InputCollectState_wcall() {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_CollectState)
	master.connector.WriteTo(define.ObserverId, &msg)
	utils.LogI(fmt.Sprintf("Sent InputCollectState to nodeId %d", define.ObserverId))

	master.connector.WaitAckRsp(define.ObserverId, define.Input_CollectStateRsp)
}
