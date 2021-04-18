package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/node"
	"github.com/rtntubmt97/springprj/utils"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Invalid arguments for node")
		return
	}

	configPath := os.Args[1]

	err, _ := utils.ReloadConfig(configPath)
	if err != nil {
		fmt.Println("invalid config path")
	}

	define.MasterId = utils.LoadedConfig.MasterId
	define.MasterPort = utils.LoadedConfig.MasterPort
	define.ObserverId = utils.LoadedConfig.ObserverId
	define.ObserverPort = utils.LoadedConfig.ObserverPort
	utils.UseLog = utils.LoadedConfig.UseLog

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		utils.LogE("Invalid id argument")
		return
	}

	initMoney, err := strconv.Atoi(os.Args[3])
	if err != nil {
		utils.LogE("Invalid initMoney argument")
		return
	}

	port := utils.GetAvailablePort(9000, 9500)

	node := node.Node{}
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
