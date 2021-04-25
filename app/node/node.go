package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rtntubmt97/springprj/impl"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Invalid arguments for node")
		return
	}

	err, _ := impl.ReloadConfig()
	if err != nil {
		fmt.Println("invalid config path")
	}

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		impl.LogE("Invalid id argument")
		return
	}

	initMoney, err := strconv.Atoi(os.Args[2])
	if err != nil {
		impl.LogE("Invalid initMoney argument")
		return
	}

	port := impl.GetAvailablePort(9000, 9500)

	node := impl.Node{}
	node.Init(int32(id))
	node.SetMoney(int64(initMoney))
	go node.Listen(port)
	time.Sleep(100 * time.Millisecond)
	node.WaitReady()
	node.ConnectMaster()
	node.ConnectObserver()
	node.ConnectPeers()

	time.Sleep(999 * time.Hour)
}
