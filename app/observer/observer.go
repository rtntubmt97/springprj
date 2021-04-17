package main

import (
	"os"
	"time"

	"github.com/rtntubmt97/springprj/observer"
	"github.com/rtntubmt97/springprj/utils"
)

func main() {
	if len(os.Args) != 2 {
		utils.LogE("Invalid arguments")
		return
	}

	observer := observer.Observer{}
	observer.Init()
	go observer.Listen()
	// time.Sleep(1000 * time.Millisecond)
	observer.ConnectMaster()

	time.Sleep(999 * time.Hour)
}
