package node

import (
	"fmt"
	"os"

	"github.com/rtntubmt97/springprj/define"
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
