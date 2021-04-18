package node

import (
	"github.com/rtntubmt97/springprj/connector"
	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
)

type Node struct {
	id            int32
	money         int64
	connector     connectorPkg.Connector
	moneyChannels map[int32]([]MoneyTokenInfo)
	snapShot      SnapShot

	// this flow can be implemented more easy if golang has tree structure
	// infoId int32
	// nextTokenId      int32
	// readStatusInfoId map[int32]bool
}

type MoneyTokenInfo struct {
	// id       int32
	SenderId int32
	Money    int32
}

type SnapShot struct {
	NodeMoney     int64
	ChannelMoneys map[int32]int64
}

func (info MoneyTokenInfo) IsToken() bool {
	return info.Money == -1
}

func (node *Node) IsConnected(otherId int32) bool {
	return node.connector.IsConnected(otherId)
}

func (node *Node) Init(id int32) {
	node.id = id
	connector := connectorPkg.Connector{}
	connector.Init(id)
	connector.ParticipantType = connectorPkg.NodeType

	connector.SetAfterAccept(node.afterAccept)

	connector.SetHandleFunc(define.SendInt32, node.sendInt32_handle)
	connector.SetHandleFunc(define.SendInt64, node.sendInt64_handle)
	connector.SetHandleFunc(define.SendString, node.sendString_handle)
	connector.SetHandleFunc(define.Input_Kill, node.inputKill_whandle)
	connector.SetHandleFunc(define.Input_Send, node.inputSend_whandle)
	connector.SetHandleFunc(define.Send, node.send_whandle)
	connector.SetHandleFunc(define.Input_Recieve, node.inputReceive_whandle)
	connector.SetHandleFunc(define.Input_RecieveAll, node.inputReceiveAll_whandle)
	connector.SetHandleFunc(define.Input_BeginSnapshot, node.inputBeginSnapshot_whandle)
	connector.SetHandleFunc(define.SendToken, node.sendToken_whandle)
	connector.SetHandleFunc(define.CollectState, node.collectState_whandle)

	node.connector = connector
	node.moneyChannels = make(map[int32]([]MoneyTokenInfo))

	// node.infoId = 0
	// node.nextTokenId = -1
	// node.readStatusInfoId = map[int32]bool{}
}

func (node *Node) GetId() int32 {
	return node.id
}

func (node *Node) SetMoney(money int64) {
	node.money = money
}

func (node *Node) Start() {

}

func (node *Node) Listen(port int) {
	node.connector.Listen(port)
}

func (node *Node) Connect(id int32, port int32) {
	node.connector.Connect(id, port)
}

func (node *Node) ConnectMaster() {
	node.Connect(define.MasterId, define.MasterPort)
}

func (node *Node) ConnectObserver() {
	node.Connect(define.ObserverId, define.ObserverPort)
}

func (node *Node) ConnectPeers() {
	otherNodeListenPorts := node.RequestInfo_wcall()
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

func (node *Node) WaitReady() {
	node.connector.WaitReady()
}

func (node *Node) WaitRsp(connId int32) define.MessageBuffer {
	return node.connector.WaitRsp(connId)
}

func (node *Node) afterAccept(conInfo connector.ParticipantInfo) {
	node.moneyChannels[conInfo.NodeId] = make([]MoneyTokenInfo, 0)
}

// func (node *Node) nextInfoId() int32 {
// 	ret := node.infoId
// 	node.infoId++
// 	return ret
// }

// func (node *Node) setReadInfoId(infoId int32) {
// 	delete(node.readStatusInfoId, infoId)
// }

// func (node *Node) shouldReadMasterToken(infoId int32) bool {
// 	tokenId := int32(-1)
// 	masterChannel := node.moneyChannels[define.MasterId]
// 	if len(masterChannel) == 0 {
// 		utils.LogI(fmt.Sprintf("Node %d check for infoId %d, tokenId is %d\n", node.id, infoId, tokenId))
// 		return false
// 	}

// 	nextMasterTokenId := masterChannel[0].id
// 	tokenId = nextMasterTokenId
// 	utils.LogI(fmt.Sprintf("Node %d check for infoId %d, tokenId is %d\n", node.id, infoId, tokenId))
// 	if infoId == nextMasterTokenId {
// 		utils.LogE("Wrong infoId")
// 	}

// 	return infoId > nextMasterTokenId
// }

// func (node *Node) updateNextTokenId() {
// 	masterChannel := node.moneyChannels[define.MasterId]
// 	if len(masterChannel) == 0 {
// 		return
// 	}
// 	node.nextTokenId = masterChannel[0].id
// 	node.moneyChannels[define.MasterId] = node.moneyChannels[define.MasterId][1:]
// }
