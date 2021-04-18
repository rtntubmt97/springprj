package node

import (
	"fmt"
	"os"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (node *Node) sendInt32_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received Int32 %d", msg.ReadI32()))
}

func (node *Node) sendInt64_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received Int64 %d", msg.ReadI64()))
}

func (node *Node) sendString_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received String %s", msg.ReadString()))
}

func (node *Node) inputKill_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received kill_handle signal", node.id))
	node.connector.SendAckRsp(connId, define.Input_KillRsp)
	os.Exit(0)
}

func (node *Node) inputReceive_whandle(connId int32, msg define.MessageBuffer) {
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

	node.connector.SendAckRsp(connId, define.Input_RecieveRsp)
}

func (node *Node) inputReceiveAll_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputReceiveAll signal", node.id))
	count := 0
	for nodeId, channel := range node.moneyChannels {
		if nodeId == define.ObserverId ||
			nodeId == define.MasterId ||
			len(channel) == 0 {
			continue
		}
		for i := 0; i < len(channel); i++ {
			count++
			if count > 5 {
				break
			}
			moneyTokenInfo := channel[i]
			node.processInfo(moneyTokenInfo, false)
		}
		node.moneyChannels[nodeId] = channel[:0]
	}

	node.connector.SendAckRsp(connId, define.Input_RecieveAllRsp)
}

func (node *Node) processInfo(info MoneyTokenInfo, print bool) {
	money := info.Money
	sender := info.SenderId
	var output define.ProjectOutput

	if info.IsToken() {
		utils.LogI(fmt.Sprintf("Node %d received token from %d", node.id, sender))
		node.updateSnapShot()
		output = utils.CreateReceiveSnapshotOutput(sender)
	} else {
		node.money += int64(money)
		utils.LogI(fmt.Sprintf("Node %d added %d money from sender %d", node.id, money, sender))
		output = utils.CreateTransferOutput(sender, money)
	}

	if print {
		utils.LogR(output)
	}
}

func (node *Node) updateSnapShot() {
	newSnapShot := SnapShot{}
	newSnapShot.NodeMoney = node.money
	channels := make(map[int32]int64)
	for nodeId, channel := range node.moneyChannels {
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
}

func (node *Node) inputSend_whandle(connId int32, msg define.MessageBuffer) {
	receiver := msg.ReadI32()
	money := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received inputSend signal, receiver is %d, money is %d", node.id, receiver, money))
	if money > int32(node.money) {
		utils.LogR(define.ERR_SEND)
		return
	}
	node.Send_wcall(receiver, money)

	node.connector.SendAckRsp(connId, define.Input_SendRsp)
}

func (node *Node) inputBeginSnapshot_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputPrintSnapshot signal", node.id))
	utils.LogR(utils.CreateBeginSnapshotOutput(node.id))
	node.updateSnapShot()
	node.propagateToken()

	node.connector.SendAckRsp(connId, define.Input_BeginSnapshotRsp)
}

func (node *Node) sendToken_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received sendToken from node %d", node.id, connId))
	newInfo := MoneyTokenInfo{SenderId: connId, Money: -1}
	node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)

	node.connector.SendAckRsp(connId, define.SendTokenRsp)
}

func (node *Node) collectState_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received collectState", node.id))
	rspMsg := protocol.SimpleMessageBuffer{}
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

// func getMoneySlice(infos []MoneyTokenInfo) []int32 {
// 	ret := make([]int32, 0)
// 	for _, info := range infos {
// 		if info.IsToken() {
// 			continue
// 		}
// 		ret = append(ret, info.Money)
// 	}
// 	return ret
// }

func (node *Node) send_whandle(connId int32, msg define.MessageBuffer) {
	money := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received money %d from node %d", node.id, money, connId))
	newInfo := MoneyTokenInfo{SenderId: connId, Money: money}
	node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)
	// node.moneyChannels[connId] <- MoneyTokenInfo{SenderId: connId, Money: money}
	// rspMsg := new(protocol.SimpleMessageBuffer)
	// rspMsg.Init(define.Rsp)
	// rspMsg.WriteI32(int32(define.SendRsp))
	// node.connector.WriteTo(connId, rspMsg)

	node.connector.SendAckRsp(connId, define.SendRsp)
}
