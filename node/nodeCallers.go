package node

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (node *Node) SendInt32_call(otherNodeId int32, i int32) {
	conn := node.connector.ConnectedConns[otherNodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.SendInt32)
	msg.WriteI32(i)
	node.connector.WriteTo(otherNodeId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent Int32 %d", i))
}

func (node *Node) SendInt64_call(otherNodeId int32, i int64) {
	conn := node.connector.ConnectedConns[otherNodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.SendInt64)
	msg.WriteI64(i)
	node.connector.WriteTo(otherNodeId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent Int64 %d", i))
}

func (node *Node) SendString_call(otherNodeId int32, s string) {
	conn := node.connector.ConnectedConns[otherNodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.SendString)
	msg.WriteString(s)
	node.connector.WriteTo(otherNodeId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent String %s", s))
}

func (node *Node) RequestInfo_wcall() map[int32]int32 {
	conn := node.connector.ConnectedConns[define.MasterId]
	if conn == nil {
		utils.LogE("nil conn")
		return nil
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.RequestInfo)
	node.connector.WriteTo(define.MasterId, &msg)
	// msg.Write(conn)

	utils.LogI(fmt.Sprintf("Node %d request info", node.id))

	utils.LogI(fmt.Sprintf("Requested from address %s", conn.LocalAddr().String()))

	rspMsg := node.connector.WaitRsp(define.MasterId)

	utils.LogI("Received")

	ret := make(map[int32]int32)
	cmd := define.ConnectorCmd(rspMsg.ReadI32())
	if cmd != define.RequestInfoRsp {
		utils.LogE(fmt.Sprintf("node %d got invalid response for RequestInfo_wcall", node.id))
		return nil
	}

	connN := rspMsg.ReadI32()
	for i := int32(0); i < connN; i++ {
		otherId := rspMsg.ReadI32()
		otherListenPort := rspMsg.ReadI32()
		ret[otherId] = otherListenPort
		utils.LogI(fmt.Sprintf("connId %d", otherId))
		utils.LogI(fmt.Sprintf("port %d", otherListenPort))
	}

	return ret
}

func (node *Node) Send_wcall(receiver int32, money int32) {
	conn := node.connector.ConnectedConns[receiver]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.Send)
	msg.WriteI32(money)
	node.connector.WriteTo(receiver, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Send_wcall to %d, money is %d", receiver, money))
	rspMsg := node.WaitRsp(receiver)
	cmd := define.ConnectorCmd(rspMsg.ReadI32())
	if cmd != define.SendRsp {
		utils.LogE("Wrong send response")
	} else {
		utils.LogI("Success send")
	}
	node.money -= int64(money)
}

func (node *Node) SendToken_wcall(connId int32) {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.SendToken)
	node.connector.WriteTo(connId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent token to node %d", connId))

	node.connector.WaitAckRsp(connId, define.SendTokenRsp)
}

func (node *Node) propagateToken() {
	for nodeId := range node.connector.ConnectedConns {
		if nodeId == define.MasterId || nodeId == define.ObserverId {
			continue
		}
		node.SendToken_wcall(nodeId)
	}
}
