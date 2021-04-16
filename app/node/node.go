package main

import (
	"os"
	"strconv"
	"time"

	"github.com/rtntubmt97/springprj/node"
	"github.com/rtntubmt97/springprj/utils"
)

func main() {
	if len(os.Args) != 4 {
		utils.LogE("Invalid arguments")
		return
	}

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		utils.LogE("Invalid port argument")
		return
	}

	initMoney, err := strconv.Atoi(os.Args[2])
	if err != nil {
		utils.LogE("Invalid port argument")
		return
	}

	port := utils.GetAvailablePort(9000, 9500)

	node := node.Node{}
	node.Init(int32(id))
	node.SetMoney(int64(initMoney))
	go node.Listen(port)
	node.ConnectMaster()

	time.Sleep(999 * time.Hour)
}
