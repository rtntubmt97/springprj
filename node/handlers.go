package node

import (
	"fmt"

	"github.com/rtntubmt97/springprj/protocol"
)

func (node *Node) sendInt32_handle(msg protocol.MessageBuffer) {
	// fmt.Println("SendInt32_handle run!")
	fmt.Printf("Receive Int32 %d\n", msg.ReadI32())
}

func (node *Node) sendInt64_handle(msg protocol.MessageBuffer) {
	// fmt.Println("SendInt64_handle run!")
	fmt.Printf("Receive Int64 %d\n", msg.ReadI64())
}

func (node *Node) sendString_handle(msg protocol.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	fmt.Printf("Receive String %s\n", msg.ReadString())
}
