// Node is the main actor in this project. It will receive the signal from the master
// and handle them with their own business logic. Of course its logic has to follow the
// chandy lamport algorithm.

// Node package includes 3 files: node.go, nodeCallers.go and nodeHandlers.go.
// Node.go file contains the node structure and its basic operation.
// nodeCallers.go file contains its caller which will be used to send data to other
// to nodes or observer by specific message.
// nodeHandlers.go file contains its handler which will be used to handle the
// incoming messages.

package node

import (
	"fmt"
	"os"

	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

// The Node struct contains its id, current money, a connector, an id->moneytoken_channel map, a snapshot
type Node struct {
	id            int32
	money         int64
	connector     connectorPkg.Connector
	moneyChannels map[int32]([]MoneyTokenInfo)
	snapShot      SnapShot
}

// MoneyTokenInfo struct contains the sender money and the amount of money (-1 if this is a token).
type MoneyTokenInfo struct {
	// id       int32
	SenderId int32
	Money    int32
}

// Snapshot struct contains the current money (node state) and a channel=>money map (channel states).
// This structure will be move to the define package later.
type SnapShot struct {
	NodeMoney     int64
	ChannelMoneys map[int32]int64
}

// Determine whether the info (usualy in channel) is money or a token.
func (info MoneyTokenInfo) IsToken() bool {
	return info.Money == -1
}

// Initilize the node.
func (node *Node) Init(id int32) {
	node.id = id
	connector := connectorPkg.Connector{}
	connector.Init(id)
	connector.ParticipantType = connectorPkg.NodeType

	connector.SetAfterAccept(node.afterAccept)

	connector.SetHandleFunc(define.SendInt32, node.sendInt32Nowait)
	connector.SetHandleFunc(define.SendInt64, node.sendInt64Nowait)
	connector.SetHandleFunc(define.SendString, node.sendStringNowait)
	connector.SetHandleFunc(define.Input_Kill, node.signalKillHandle)
	connector.SetHandleFunc(define.Input_Send, node.signalSendHandle)
	connector.SetHandleFunc(define.Send, node.sendHandle)
	connector.SetHandleFunc(define.Input_receive, node.signalReceiveHandle)
	connector.SetHandleFunc(define.Input_receiveAll, node.signalReceiveAllHandle)
	connector.SetHandleFunc(define.Input_BeginSnapshot, node.signalBeginSnapshotHandle)
	connector.SetHandleFunc(define.SendToken, node.sendTokenHandle)
	connector.SetHandleFunc(define.CollectState, node.collectStateHandle)

	node.connector = connector
	node.moneyChannels = make(map[int32]([]MoneyTokenInfo))

	// node.infoId = 0
	// node.nextTokenId = -1
	// node.readStatusInfoId = map[int32]bool{}
}

// Get the node id.
func (node *Node) GetId() int32 {
	return node.id
}

// Set the current node money.
func (node *Node) SetMoney(money int64) {
	node.money = money
}

// Start the listen operation.
func (node *Node) Listen(port int) {
	node.connector.Listen(port)
}

// Determin whether the Node has connected to a node specified by an id.
func (node *Node) IsConnected(otherId int32) bool {
	return node.connector.IsConnected(otherId)
}

// connect to a connector (master, observer or node).
func (node *Node) Connect(id int32, port int32) {
	node.connector.Connect(id, port)
}

// Connect to master.
func (node *Node) ConnectMaster() {
	node.Connect(define.MasterId, define.MasterPort)
}

// Connect to observer.
func (node *Node) ConnectObserver() {
	node.Connect(define.ObserverId, define.ObserverPort)
}

// Connect to other nodes, after request other nodes information from the master.
func (node *Node) ConnectPeers() {
	otherNodeListenPorts := node.RequestInfo()
	for nodeId, port := range otherNodeListenPorts {
		if nodeId == node.GetId() {
			continue
		}
		if node.IsConnected(nodeId) {
			continue
		}
		node.Connect(nodeId, port)
		node.moneyChannels[nodeId] = make([]MoneyTokenInfo, 0)
	}
}

// Wait and return when the node ready.
func (node *Node) WaitReady() {
	node.connector.WaitReady()
}

// Wait for a response message and return it.
func (node *Node) WaitRsp(connId int32) define.MessageBuffer {
	return node.connector.WaitRsp(connId)
}

// The callback will be call after a new connector accepted by this node connector.
func (node *Node) afterAccept(conInfo connectorPkg.ParticipantInfo) {
	node.moneyChannels[conInfo.NodeId] = make([]MoneyTokenInfo, 0)
}

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

// Handle simple caller sendInt32, just print it.
func (node *Node) sendInt32Nowait(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received Int32 %d", msg.ReadI32()))
}

// Handle simple caller sendInt64, just print it.
func (node *Node) sendInt64Nowait(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received Int64 %d", msg.ReadI64()))
}

// Handle simple caller sendString, just print it.
func (node *Node) sendStringNowait(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received String %s", msg.ReadString()))
}

// Handle the kill signal from master, kill the current node process.
func (node *Node) signalKillHandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received kill_handle signal", node.id))
	node.connector.SendAckRsp(connId, define.Input_KillRsp)
	os.Exit(0)
}

// Handle the Receive signal from master, add money from a channel or capture a snapshot
// whether it is a token.
func (node *Node) signalReceiveHandle(connId int32, msg define.MessageBuffer) {
	sender := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received inputReceive signal, sender is %d", node.id, sender))
	var selectedChannel []MoneyTokenInfo
	if sender != int32(-1) {
		selectedChannel = node.moneyChannels[sender]
	} else {
		for id, channel := range node.moneyChannels {
			if len(channel) != 0 {
				sender = id
				selectedChannel = channel
				break
			}
		}
	}

	moneyTokenInfo := selectedChannel[0]
	node.moneyChannels[sender] = selectedChannel[1:]
	node.processInfo(moneyTokenInfo, true)

	node.connector.SendAckRsp(connId, define.Input_receiveRsp)
}

// Handle the ReceiveAll signal from master, drain out all the channel.
func (node *Node) signalReceiveAllHandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputReceiveAll signal", node.id))

	// node.receiveAllProccess(define.MasterId)

	for nodeId, channel := range node.moneyChannels {
		if nodeId == define.ObserverId ||
			nodeId == define.MasterId ||
			len(channel) == 0 {
			continue
		}
		node.receiveAllProccess(nodeId)
	}

	node.connector.SendAckRsp(connId, define.Input_receiveAllRsp)
}

//  Drain out a channel specific by its id.
func (node *Node) receiveAllProccess(nodeId int32) {
	channel := node.moneyChannels[nodeId]
	for i := 0; i < len(channel); i++ {
		moneyTokenInfo := channel[i]
		node.processInfo(moneyTokenInfo, false)
	}
	node.moneyChannels[nodeId] = channel[:0]
}

// Process money/token in channel. Add the money to the node money or update snapshot.
func (node *Node) processInfo(info MoneyTokenInfo, print bool) {
	money := info.Money
	sender := info.SenderId
	var output define.ProjectOutput

	if info.IsToken() {
		utils.LogI(fmt.Sprintf("Node %d received token from %d", node.id, sender))
		node.updateSnapShot(sender)
		output = utils.CreateReceiveSnapshotOutput(sender)
	} else {
		node.money += int64(money)
		utils.LogI(fmt.Sprintf("Node %d added %d money from sender %d", node.id, money, sender))
		output = utils.CreateTransferOutput(sender, money)
	}

	if print && sender != define.MasterId {
		utils.LogR(output)
	}
}

// Update node snapshot of current stage which later will be sent to the observer.
func (node *Node) updateSnapShot(tokenSender int32) {
	newSnapShot := SnapShot{}
	newSnapShot.NodeMoney = node.money
	channels := make(map[int32]int64)
	for nodeId, channel := range node.moneyChannels {
		if nodeId == tokenSender {
			channels[nodeId] = 0
			continue
		}
		totalMoney := int64(0)
		for _, money := range channel {
			if money.Money == -1 {
				continue
			}
			totalMoney += int64(money.Money)
		}
		channels[nodeId] = totalMoney
	}
	newSnapShot.ChannelMoneys = channels
	node.snapShot = newSnapShot
	utils.LogI("newSnapShot")
}

// Handle the Send signal from master, send money to other node.
func (node *Node) signalSendHandle(connId int32, msg define.MessageBuffer) {
	receiver := msg.ReadI32()
	money := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received inputSend signal, receiver is %d, money is %d", node.id, receiver, money))
	if money > int32(node.money) {
		utils.LogR(define.ERR_SEND)
		return
	}
	node.Send(receiver, money)

	node.connector.SendAckRsp(connId, define.Input_SendRsp)
}

// Handle the BeginSnapshot signal from master, update the token and propagate it to other nodes.
func (node *Node) signalBeginSnapshotHandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputPrintSnapshot signal", node.id))
	utils.LogR(utils.CreateBeginSnapshotOutput(node.id))
	// newInfo := MoneyTokenInfo{SenderId: connId, Money: -1}
	// node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)
	node.updateSnapShot(connId)
	node.propagateToken()

	node.connector.SendAckRsp(connId, define.Input_BeginSnapshotRsp)
}

// Handle the token received from other node, put it to corresponding moneytoken channel.
func (node *Node) sendTokenHandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received sendToken from node %d", node.id, connId))
	newInfo := MoneyTokenInfo{SenderId: connId, Money: -1}
	node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)

	node.connector.SendAckRsp(connId, define.SendTokenRsp)
}

// Handle the CollectState signal from observer, resonse the last updated snapshot
// to the observer.
func (node *Node) collectStateHandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received collectState", node.id))
	rspMsg := protocol.BinaryProtocol{}
	rspMsg.Init(define.Rsp)
	rspMsg.WriteI32(int32(define.CollectStateRsp))

	snapShot := node.snapShot
	rspMsg.WriteI64(snapShot.NodeMoney)
	rspMsg.WriteI32(int32(len(snapShot.ChannelMoneys)))
	for channelId, channelMoney := range snapShot.ChannelMoneys {
		rspMsg.WriteI32(channelId)
		// moneySlice := getMoneySlice(channel)
		rspMsg.WriteI64(channelMoney)
	}
	node.connector.WriteTo(connId, &rspMsg)
}

// Handle a Send signal from other node, retrieve the money amount from the message, and
// put it to the corresponding money channel.
func (node *Node) sendHandle(connId int32, msg define.MessageBuffer) {
	money := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received send_call money %d from node %d", node.id, money, connId))
	newInfo := MoneyTokenInfo{SenderId: connId, Money: money}
	node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)

	node.connector.SendAckRsp(connId, define.SendRsp)
}
