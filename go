StartMaster
CreateNode 12 12345
package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rtntubmt97/springprj/master"
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

func getInput() []string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		utils.LogE("Invalid input")
		utils.LogE(err.Error())
		os.Exit(1)
	}

	splitInput := strings.Split(input, " ")
	for i, ele := range splitInput {
		splitInput[i] = strings.Trim(ele, "\n ")
	}

	return splitInput
}

func createNode(id string, initMoney string) {
	exeCmd := exec.Command("go", "run", "app/node/node.go", id, initMoney, "somethingtohtop")
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr
	exeCmd.Start()
}

var masterNode master.Master

func main() {

	for {
		input := getInput()
		cmd := InputCmd(input[0])

		switch cmd {
		case StartMaster:
			utils.LogI("Matched StartMaster")
			masterNode = master.Master{}
			masterNode.Init()
			go masterNode.Listen()

		case KillAll:
			utils.LogI("Matched KillAll")
			masterNode.KillAll()
			os.Exit(0)

		case CreateNode:
			utils.LogI("Matched CreateNode")
			createNode(input[1], input[2])

		default:
		}
		time.Sleep(1 * time.Second)
	}
}
