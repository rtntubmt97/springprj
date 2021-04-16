package main

import (
	"time"

	"github.com/rtntubmt97/springprj/node"
)

func main() {
	node := node.Node{}
	node.Init(1)
	go node.Listen(9090)
	time.Sleep(999 * time.Minute)
}
