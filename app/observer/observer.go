package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rtntubmt97/springprj/impl"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid arguments for observer")
		fmt.Println(len(os.Args))
		return
	}

	err, _ := impl.ReloadConfig()
	if err != nil {
		fmt.Println("invalid config path")
	}

	observer := impl.Observer{}
	observer.Init()
	go observer.Listen()
	// time.Sleep(1000 * time.Millisecond)
	observer.ConnectMaster()

	time.Sleep(999 * time.Hour)
}
