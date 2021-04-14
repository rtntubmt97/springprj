package main

import (
	"time"

	"github.com/rtntubmt97/springprj/node"
)

func main() {
	node := node.Node{}
	node.Init()
	node.ConnectNode(9090)
	node.SendInt32(0, 123)
	node.SendInt32(0, 123123)
	node.SendInt64(0, 121412412345)
	node.SendString(0, "asfsdfdaa\"\\/\\")
	node.SendString(0, "@#$%^&88sdfs")
	node.SendInt64(0, 2345678986)
	time.Sleep(999 * time.Second)
}
