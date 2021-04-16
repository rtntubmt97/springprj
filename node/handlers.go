package node

import (
	"fmt"

	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (node *Node) sendInt32_handle(connId int32, msg protocol.MessageBuffer) {
	// fmt.Println("SendInt32_handle run!")
	utils.LogI(fmt.Sprintf("Received Int32 %d", msg.ReadI32()))
}

func (node *Node) sendInt64_handle(connId int32, msg protocol.MessageBuffer) {
	// fmt.Println("SendInt64_handle run!")
	utils.LogI(fmt.Sprintf("Received Int64 %d", msg.ReadI64()))
}

func (node *Node) sendString_handle(connId int32, msg protocol.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	utils.LogI(fmt.Sprintf("Received String %s", msg.ReadString()))
}
