// The "w" character in the whandle post fix means that handle function will send an
// ack message or response message to the caller before return
// The "Input" in the call prefix mean that this caller will be call by the master after
// receive a correspond input from user.

// Node can handle the simpler callers sendInt32, sendInt64, sendString from other
// nodes/master/observer
// Node can handle the Send, receive, receiveAll signal from the master and BeginSnapShot,
// CollectState from the observer as the specification described.

package node

import (
	"fmt"
	"os"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

// Handle simple caller sendInt32, just print it.
func (node *Node) sendInt32_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received Int32 %d", msg.ReadI32()))
}

// Handle simple caller sendInt64, just print it.
func (node *Node) sendInt64_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received Int64 %d", msg.ReadI64()))
}

// Handle simple caller sendString, just print it.
func (node *Node) sendString_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Received String %s", msg.ReadString()))
}

// Handle the kill signal from master, kill the current node process.
func (node *Node) inputKill_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received kill_handle signal", node.id))
	node.connector.SendAckRsp(connId, define.Input_KillRsp)
	os.Exit(0)
}

// Handle the Receive signal from master, add money from a channel or capture a snapshot
// whether it is a token.
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

	node.connector.SendAckRsp(connId, define.Input_receiveRsp)
}

// Handle the ReceiveAll signal from master, drain out all the channel.
func (node *Node) inputReceiveAll_whandle(connId int32, msg define.MessageBuffer) {
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

// Handle the BeginSnapshot signal from master, update the token and propagate it to other nodes.
func (node *Node) inputBeginSnapshot_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputPrintSnapshot signal", node.id))
	utils.LogR(utils.CreateBeginSnapshotOutput(node.id))
	// newInfo := MoneyTokenInfo{SenderId: connId, Money: -1}
	// node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)
	node.updateSnapShot(connId)
	node.propagateToken()

	node.connector.SendAckRsp(connId, define.Input_BeginSnapshotRsp)
}

// Handle the token received from other node, put it to corresponding moneytoken channel.
func (node *Node) sendToken_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received sendToken from node %d", node.id, connId))
	newInfo := MoneyTokenInfo{SenderId: connId, Money: -1}
	node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)

	node.connector.SendAckRsp(connId, define.SendTokenRsp)
}

// Handle the CollectState signal from observer, resonse the last updated snapshot
// to the observer.
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

// Handle a Send signal from other node, retrieve the money amount from the message, and
// put it to the corresponding money channel.
func (node *Node) send_whandle(connId int32, msg define.MessageBuffer) {
	money := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received send_call money %d from node %d", node.id, money, connId))
	newInfo := MoneyTokenInfo{SenderId: connId, Money: money}
	node.moneyChannels[connId] = append(node.moneyChannels[connId], newInfo)

	node.connector.SendAckRsp(connId, define.SendRsp)
}
