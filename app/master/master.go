package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rtntubmt97/springprj/utils"
)

type InputCmd string

const (
	StartMaster   InputCmd = "StartMaster"
	KillAll       InputCmd = "KillAll"
	CreateNode    InputCmd = "CreateNode"
	Send          InputCmd = "Send"
	Receive       InputCmd = "Receive"
	ReceiveAll    InputCmd = "ReceiveAll"
	BeginSnapshot InputCmd = "BeginSnapshot"
	CollectState  InputCmd = "CollectState"
	PrintSnapshot InputCmd = "PrintSnapshot"
)

func main() {
	sm := InputCmd("StartMaster")
	t := string(sm)
	if sm == StartMaster {
		fmt.Println("true")
	}
	if StartMaster == InputCmd(t) {
		fmt.Println("true2")
	}
	fmt.Println(StartMaster)
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			utils.LogE("Invalid input")
			return
		}
		inputCmd := InputCmd(text)
		switch inputCmd {
		case StartMaster:
		default:
		}
	}
}
