// The "w" character in the wcall post fix means that call function will wait for an
// ack message or response message with data before return
// The "Input" in the call prefix mean that this caller will be trigger by input of user.
// They will be called in the app/master/master.go when the correspond input matched

// Node can send a simple int32, int64, string to master/observer/nodes.
// Node can send a request information call to the master and get the data.
// Node can send Send(_call) to other node (that node will receive it, and put the money
// in the corresponding channel to the sender). Node also can send a token to other nodes
// to propagate the beginSnapshot command of the master.

package node

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

// Send a simple message contains an int32
func (node *Node) SendInt32Nowait(otherNodeId int32, i int32) {
	conn := node.connector.ConnectedConns[otherNodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.BinaryProtocol{}
	msg.Init(define.SendInt32)
	msg.WriteI32(i)
	node.connector.WriteTo(otherNodeId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent Int32 %d", i))
}

// Send a simple message contains an int64.
func (node *Node) SendInt64Nowait(otherNodeId int32, i int64) {
	conn := node.connector.ConnectedConns[otherNodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.BinaryProtocol{}
	msg.Init(define.SendInt64)
	msg.WriteI64(i)
	node.connector.WriteTo(otherNodeId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent Int64 %d", i))
}

// Send a simple message contains a string.
func (node *Node) SendStringNowait(otherNodeId int32, s string) {
	conn := node.connector.ConnectedConns[otherNodeId]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.BinaryProtocol{}
	msg.Init(define.SendString)
	msg.WriteString(s)
	node.connector.WriteTo(otherNodeId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent String %s", s))
}

// Send an information request to master node.
func (node *Node) RequestInfo() map[int32]int32 {
	conn := node.connector.ConnectedConns[define.MasterId]
	if conn == nil {
		utils.LogE("nil conn")
		return nil
	}
	msg := protocol.BinaryProtocol{}
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

// Send money to others node.
func (node *Node) Send(receiver int32, money int32) {
	conn := node.connector.ConnectedConns[receiver]
	if conn == nil {
		utils.LogE("nil conn")
		return
	}
	msg := protocol.BinaryProtocol{}
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

// Send token to other node.
func (node *Node) SendToken(connId int32) {
	msg := protocol.BinaryProtocol{}
	msg.Init(define.SendToken)
	node.connector.WriteTo(connId, &msg)
	// msg.Write(conn)
	utils.LogI(fmt.Sprintf("Sent token to node %d", connId))

	node.connector.WaitAckRsp(connId, define.SendTokenRsp)
}

// Propagate token to all nodes.
func (node *Node) propagateToken() {
	for nodeId := range node.connector.ConnectedConns {
		if nodeId == define.MasterId || nodeId == define.ObserverId {
			continue
		}
		node.SendToken(nodeId)
	}
}
