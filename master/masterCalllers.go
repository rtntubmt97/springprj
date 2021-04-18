package master

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

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

func (master *Master) InputReceive_wcall(receiver int32, sender int32) {
	conn := master.connector.ConnectedConns[receiver]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_Recieve)
	msg.WriteI32(sender)
	master.connector.WriteTo(receiver, &msg)
	utils.LogI(fmt.Sprintf("Sent inputReceive to nodeId %d, sender is %d", receiver, sender))

	master.connector.WaitAckRsp(receiver, define.Input_RecieveRsp)
}

func (master *Master) InputReceiveAll_wcall() {
	for connId := range master.connector.ConnectedConns {
		if connId == define.MasterId || connId == define.ObserverId {
			continue
		}
		msg := protocol.SimpleMessageBuffer{}
		msg.Init(define.Input_RecieveAll)
		master.connector.WriteTo(connId, &msg)
		utils.LogI(fmt.Sprintf("Sent inputReceiveAll to nodeId %d", connId))

		master.connector.WaitAckRsp(connId, define.Input_RecieveAllRsp)
	}
}

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

func (master *Master) InputPrintSnapshot_wcall() {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_PrintSnapshot)
	master.connector.WriteTo(define.ObserverId, &msg)
	utils.LogI(fmt.Sprintf("Sent InputPrintSnapshot to nodeId %d", define.ObserverId))

	master.connector.WaitAckRsp(define.ObserverId, define.Input_PrintSnapshotRsp)
}

func (master *Master) InputCollectState_wcall() {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_CollectState)
	master.connector.WriteTo(define.ObserverId, &msg)
	utils.LogI(fmt.Sprintf("Sent InputCollectState to nodeId %d", define.ObserverId))

	master.connector.WaitAckRsp(define.ObserverId, define.Input_CollectStateRsp)
}
