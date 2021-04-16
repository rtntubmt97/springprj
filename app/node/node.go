package main

import (
	"os"
	"strconv"
	"time"

	"github.com/rtntubmt97/springprj/node"
	"github.com/rtntubmt97/springprj/utils"
)

func main() {
	if len(os.Args) != 3 {
		utils.LogE("Invalid arguments")
		return
	}

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		utils.LogE("Invalid port argument")
		return
	}

	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		utils.LogE("Invalid port argument")
		return
	}

	node := node.Node{}
	node.Init(int32(id))
	go node.Listen(port)

	time.Sleep(999 * time.Hour)
}