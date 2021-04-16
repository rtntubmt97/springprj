package master

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (master *Master) inputKill_call(nodeId int32) {
	conn := master.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_Kill)
	msg.WriteMessage(conn)
	utils.LogI(fmt.Sprintf("Sent kill to nodeId %d", nodeId))
}

func (master *Master) inputSend_call(nodeId int32, money int32) {
	conn := master.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Input_Send)
	msg.WriteI32(money)
	msg.WriteMessage(conn)
	utils.LogI(fmt.Sprintf("Sent inputSend to nodeId %d", nodeId))
}
