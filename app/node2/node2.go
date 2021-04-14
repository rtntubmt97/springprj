package main

import (
	"time"

	"github.com/rtntubmt97/springprj/node"
)

func main() {
	node := node.Node{}
	node.Init()
	node.ConnectNode(9090)
	time.Sleep(999 * time.Second)
}
