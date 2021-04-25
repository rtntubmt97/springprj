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
	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
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
	connector.SetHandleFunc(define.Input_Kill, node.inputKillHandle)
	connector.SetHandleFunc(define.Input_Send, node.inputSendHandle)
	connector.SetHandleFunc(define.Send, node.sendHandle)
	connector.SetHandleFunc(define.Input_receive, node.inputReceiveHandle)
	connector.SetHandleFunc(define.Input_receiveAll, node.inputReceiveAllHandle)
	connector.SetHandleFunc(define.Input_BeginSnapshot, node.inputBeginSnapshotHandle)
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
