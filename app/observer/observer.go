package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/observer"
	"github.com/rtntubmt97/springprj/utils"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Invalid arguments for observer")
		fmt.Println(len(os.Args))
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

	observer := observer.Observer{}
	observer.Init()
	go observer.Listen()
	// time.Sleep(1000 * time.Millisecond)
	observer.ConnectMaster()

	time.Sleep(999 * time.Hour)
}
