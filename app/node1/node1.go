package main

import "github.com/rtntubmt97/springprj/node"

func main() {
	node := node.Node{}
	node.Init()
	node.Listen(9090)
}
