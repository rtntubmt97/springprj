package node

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (node *Node) SendInt32_call(nodeId int32, i int32) {
	conn := node.connector.GetConnection(nodeId)
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.MessageBuffer{}
	msg.InitEmpty()
	msg.WriteI32(define.SendInt32)
	msg.WriteI32(i)
	protocol.WriteMessage(conn, msg)
	utils.LogI(fmt.Sprintf("Sent Int32 %d", i))
}

func (node *Node) SendInt64_call(nodeId int32, i int64) {
	conn := node.connector.GetConnection(nodeId)
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.MessageBuffer{}
	msg.InitEmpty()
	msg.WriteI32(define.SendInt64)
	msg.WriteI64(i)
	protocol.WriteMessage(conn, msg)
	utils.LogI(fmt.Sprintf("Sent Int64 %d", i))
}

func (node *Node) SendString_call(nodeId int32, s string) {
	conn := node.connector.GetConnection(nodeId)
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.MessageBuffer{}
	msg.InitEmpty()
	msg.WriteI32(define.SendString)
	msg.WriteString(s)
	protocol.WriteMessage(conn, msg)
	utils.LogI(fmt.Sprintf("Sent String %s", s))
}
