package node

import (
	"fmt"
	"os"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (node *Node) sendInt32_handle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendInt32_handle run!")
	utils.LogI(fmt.Sprintf("Received Int32 %d", msg.ReadI32()))
}

func (node *Node) sendInt64_handle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendInt64_handle run!")
	utils.LogI(fmt.Sprintf("Received Int64 %d", msg.ReadI64()))
}

func (node *Node) sendString_handle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	utils.LogI(fmt.Sprintf("Received String %s", msg.ReadString()))
}

func (node *Node) kill_handle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	utils.LogI(fmt.Sprintf("Node %d Received kill_handle signal", node.id))
	os.Exit(0)
}

func (node *Node) inputReceive_handle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	sender := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received inputReceive signal, sender is %d", node.id, sender))
	var selectedChannel chan MoneyTokenInfo
	if sender != int32(-1) {
		selectedChannel = node.moneyChannels[sender]
	} else {
		for _, channel := range node.moneyChannels {
			if len(channel) != 0 {
				selectedChannel = channel
				break
			}
		}
	}

	moneyTokenInfo := <-selectedChannel
	node.processInfo(moneyTokenInfo, true)
}

func (node *Node) inputReceiveAll_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputReceiveAll signal", node.id))
	for connId, channel := range node.moneyChannels {
		if connId == define.ObserverId ||
			connId == define.MasterId ||
			len(channel) == 0 {
			continue
		}
		for len(channel) > 0 {
			moneyTokenInfo := <-channel
			node.processInfo(moneyTokenInfo, false)
		}
	}

}

func (node *Node) processInfo(info MoneyTokenInfo, print bool) {
	money := info.Money
	sender := info.SenderId

	if info.IsToken() {
	} else {
		node.money += int64(money)
		utils.LogI(fmt.Sprintf("Node %d added %d money from sender %d", node.id, money, sender))
	}

	if print {
		output := utils.CreateTransferOutput(sender, money)
		utils.LogR(output)
	}
}

func (node *Node) inputSend_handle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	receiver := msg.ReadI32()
	money := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received inputSend signal, receiver is %d, money is %d", node.id, receiver, money))
	if money > int32(node.money) {
		utils.LogR(define.ERR_SEND)
		return
	}
	node.Send_wcall(receiver, money)

}

func (node *Node) send_whandle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	money := msg.ReadI32()
	utils.LogI(fmt.Sprintf("Node %d Received money %d from node %d", node.id, money, connId))
	node.moneyChannels[connId] <- MoneyTokenInfo{SenderId: connId, Money: money}
	rspMsg := new(protocol.SimpleMessageBuffer)
	rspMsg.Init(define.Rsp)
	rspMsg.WriteI32(int32(define.SendRsp))
	node.connector.WriteTo(connId, rspMsg)

}
