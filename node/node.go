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
	moneyChannels map[int32](chan MoneyTokenInfo)
}

type MoneyTokenInfo struct {
	SenderId int32
	Money    int32
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

	connector.SetAfterAccept(node.afterAccept)

	connector.SetHandleFunc(define.SendInt32, node.sendInt32_handle)
	connector.SetHandleFunc(define.SendInt64, node.sendInt64_handle)
	connector.SetHandleFunc(define.SendString, node.sendString_handle)
	connector.SetHandleFunc(define.Input_Kill, node.kill_handle)
	connector.SetHandleFunc(define.Input_Send, node.inputSend_handle)
	connector.SetHandleFunc(define.Send, node.send_whandle)
	connector.SetHandleFunc(define.Input_Recieve, node.inputReceive_handle)
	connector.SetHandleFunc(define.Input_RecieveAll, node.inputReceiveAll_handle)

	node.connector = connector
	node.moneyChannels = make(map[int32](chan MoneyTokenInfo))
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
		node.moneyChannels[nodeId] = make(chan MoneyTokenInfo, 1000)
	}
}

func (node *Node) WaitReady() {
	node.connector.WaitReady()
}

func (node *Node) WaitRsp(connId int32) define.MessageBuffer {
	return node.connector.WaitRsp(connId)
}

func (node *Node) afterAccept(conInfo connector.OtherInfo) {
	node.moneyChannels[conInfo.NodeId] = make(chan MoneyTokenInfo, 1000)
}
