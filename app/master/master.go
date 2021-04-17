package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
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

func getStdinInput() []string {
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

func createObserver() {
	exeCmd := exec.Command("go", "run", "app/observer/observer.go", "somethingtohtop")
	// exeCmd := exec.Command("./node.exe", id, initMoney, "somethingtohtop")
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr
	err := exeCmd.Start()
	if err != nil {
		utils.LogE(err.Error())
	}
}

func createNode(id string, initMoney string) {
	exeCmd := exec.Command("go", "run", "app/node/node.go", id, initMoney, "somethingtohtop")
	// exeCmd := exec.Command("./node.exe", id, initMoney, "somethingtohtop")
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr
	err := exeCmd.Start()
	if err != nil {
		utils.LogE(err.Error())
	}
}

var masterNode master.Master

// var observerNode observer.Observer

func main() {
	file, err := os.Open("input.ini")
	if err != nil {
		utils.LogE(err.Error())
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		// input := getStdinInput()
		inputRaw, err := reader.ReadString('\n')
		// fmt.Print(inputRaw)

		if err != io.EOF {
			if inputRaw[0] == ';' {
				continue
			}
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		input := strings.Split(inputRaw, " ")
		for i, ele := range input {
			input[i] = strings.Trim(ele, "\n\r\t ")
		}

		cmd := InputCmd(input[0])

		switch cmd {
		case StartMaster:
			utils.LogI("Matched StartMaster")
			masterNode = master.Master{}
			masterNode.Init()
			go masterNode.Listen()
			time.Sleep(1 * time.Second)
			createObserver()
			// observerNode = observer.Observer{}
			// observerNode.Init()
			// go observerNode.Listen()
			time.Sleep(1 * time.Second)

		case KillAll:
			utils.LogI(inputRaw)
			masterNode.KillAll()
			time.Sleep(1 * time.Second)
			utils.LogI("Ready to exit")
			os.Exit(0)

		case CreateNode:
			utils.LogI(inputRaw)
			createNode(input[1], input[2])
			time.Sleep(900 * time.Millisecond)

		case Send:
			utils.LogI(inputRaw)
			sender, _ := strconv.Atoi(input[1])
			receiver, _ := strconv.Atoi(input[2])
			money, _ := strconv.Atoi(input[3])
			masterNode.InputSend_call(int32(sender), int32(receiver), int32(money))

		case Receive:
			utils.LogI(inputRaw)
			sender := -1
			receiver := -1
			if len(input) == 2 {
				receiver, _ = strconv.Atoi(input[1])
			} else if len(input) == 3 {
				receiver, _ = strconv.Atoi(input[1])
				sender, _ = strconv.Atoi(input[2])
			}
			masterNode.InputReceive_call(int32(receiver), int32(sender))

		case ReceiveAll:
			utils.LogI(inputRaw)
			masterNode.InputReceiveAll_call()

		case BeginSnapshot:
			utils.LogI(inputRaw)
			startNodeId, _ := strconv.Atoi(input[1])
			masterNode.InputBeginSnapshot_call(int32(startNodeId))

		case CollectState:
			utils.LogI(inputRaw)
			masterNode.InputCollectState_call()

		case PrintSnapshot:
			utils.LogI(inputRaw)
			masterNode.InputPrintSnapshot_call()

		default:
		}
		time.Sleep(100 * time.Millisecond)
		if err == io.EOF {
			break
		}
	}
	time.Sleep(100 * time.Hour)
}
