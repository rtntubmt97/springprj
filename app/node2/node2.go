package main

import (
	"time"

	"github.com/rtntubmt97/springprj/node"
)

func main() {
	node := node.Node{}
	node.Init(2)
	node.Connect(1, 9090)
	node.SendInt32(1, 123)
	node.SendInt32(1, 123123)
	node.SendInt64(1, 121412412345)
	node.SendString(1, "asfsdfdaa\"\\/\\")
	node.SendString(1, "@#$%^&88sdfs")
	node.SendInt64(1, 2345678986)
	time.Sleep(999 * time.Second)
}
