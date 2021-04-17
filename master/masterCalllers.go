package master

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (master *Master) InputKill_call(nodeId int32) {
	conn := master.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_Kill)
	master.connector.WriteTo(nodeId, &msg)
	utils.LogI(fmt.Sprintf("Sent kill to nodeId %d", nodeId))
}

func (master *Master) InputSend_call(nodeId int32, receiver int32, money int32) {
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
}

func (master *Master) InputReceive_call(receiver int32, sender int32) {
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
}

func (master *Master) InputReceiveAll_call() {
	for connId := range master.connector.ConnectedConns {
		msg := protocol.SimpleMessageBuffer{}
		msg.Init(define.Input_RecieveAll)
		master.connector.WriteTo(connId, &msg)
		utils.LogI(fmt.Sprintf("Sent inputReceiveAll to nodeId %d", connId))
	}
}
