package node

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (node *Node) SendInt32_call(nodeId int32, i int32) {
	conn := node.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.SendInt32)
	msg.WriteI32(i)
	msg.WriteMessage(conn)
	utils.LogI(fmt.Sprintf("Sent Int32 %d", i))
}

func (node *Node) SendInt64_call(nodeId int32, i int64) {
	conn := node.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.SendInt64)
	msg.WriteI64(i)
	msg.WriteMessage(conn)
	utils.LogI(fmt.Sprintf("Sent Int64 %d", i))
}

func (node *Node) SendString_call(nodeId int32, s string) {
	conn := node.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.SendString)
	msg.WriteString(s)
	msg.WriteMessage(conn)
	utils.LogI(fmt.Sprintf("Sent String %s", s))
}

func (node *Node) RequestInfo_wcall(nodeId int32, s string) map[int32]int32 {
	conn := node.connector.ConnectedConns[nodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return nil
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.RequestInfo)
	msg.WriteMessage(conn)
	utils.LogI(fmt.Sprintf("Node %d request info", nodeId))

	ret := make(map[int32]int32)

	return ret
}
