package main

import (
	"time"

	"github.com/rtntubmt97/springprj/node"
)

func main() {
	node := node.Node{}
	node.Init(2)
	go node.Listen(9091)
	otherNodeId := int32(1)
	node.Connect(otherNodeId, 9090)
	node.SendInt32_call(otherNodeId, 123)
	node.SendInt32_call(otherNodeId, 123123)
	node.SendInt64_call(otherNodeId, 121412412345)
	node.SendString_call(otherNodeId, "asfsdfdaa\"\\/\\")
	node.SendString_call(otherNodeId, "@#$%^&88sdfs")
	node.SendInt64_call(otherNodeId, 2345678986)
	time.Sleep(999 * time.Second)
}
