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

	"github.com/rtntubmt97/springprj/define"
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

func createObserver(configPath string) {
	var exeCmd *exec.Cmd
	if utils.LoadedConfig.UseBin {
		exeCmd = exec.Command("bin/observer", configPath, "somethingtohtop")
	} else {
		exeCmd = exec.Command("go", "run", "app/observer/observer.go", configPath, "somethingtohtop")
	}
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr
	err := exeCmd.Start()
	if err != nil {
		utils.LogE(err.Error())
	}
}

func createNode(configPath string, id string, initMoney string) {
	var exeCmd *exec.Cmd
	if utils.LoadedConfig.UseBin {
		exeCmd = exec.Command("bin/node", configPath, id, initMoney, "somethingtohtop")
	} else {
		exeCmd = exec.Command("go", "run", "app/node/node.go", configPath, id, initMoney, "somethingtohtop")
	}
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
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	err, _ := utils.ReloadConfig(configPath)
	if err != nil {
		fmt.Println("invalid config path")
	}

	define.MasterId = utils.LoadedConfig.MasterId
	define.MasterPort = utils.LoadedConfig.MasterPort
	define.ObserverId = utils.LoadedConfig.ObserverId
	define.ObserverPort = utils.LoadedConfig.ObserverPort
	utils.UseLog = utils.LoadedConfig.UseLog

	if !utils.IsPortAvailable(int(define.MasterPort)) {
		fmt.Printf("Master node cannot using port %d", define.MasterPort)
		return
	}

	if !utils.IsPortAvailable(int(define.ObserverPort)) {
		fmt.Printf("Observer node cannot using port %d", define.ObserverPort)
		return
	}

	file, err := os.Open(utils.LoadedConfig.InputFile)
	if err != nil {
		utils.LogE(err.Error())
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		// input := getStdinInput()
		inputRaw, err := reader.ReadString('\n')

		if err != io.EOF {
			if inputRaw[0] == ';' {
				continue
			}
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		if utils.LoadedConfig.PrintInput {
			fmt.Print(inputRaw)
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
			createObserver(configPath)
			time.Sleep(1 * time.Second)

		case KillAll:
			utils.LogI(inputRaw)
			masterNode.KillAll()
			time.Sleep(1 * time.Second)
			utils.LogI("Ready to exit")
			os.Exit(0)

		case CreateNode:
			utils.LogI(inputRaw)
			createNode(configPath, input[1], input[2])
			time.Sleep(1000 * time.Millisecond)

		case Send:
			utils.LogI(inputRaw)
			sender, _ := strconv.Atoi(input[1])
			receiver, _ := strconv.Atoi(input[2])
			money, _ := strconv.Atoi(input[3])
			masterNode.InputSend_wcall(int32(sender), int32(receiver), int32(money))

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
			masterNode.InputReceive_wcall(int32(receiver), int32(sender))

		case ReceiveAll:
			utils.LogI(inputRaw)
			masterNode.InputReceiveAll_wcall()

		case BeginSnapshot:
			utils.LogI(inputRaw)
			startNodeId, _ := strconv.Atoi(input[1])
			masterNode.InputBeginSnapshot_wcall(int32(startNodeId))

		case CollectState:
			utils.LogI(inputRaw)
			masterNode.InputCollectState_wcall()

		case PrintSnapshot:
			utils.LogI(inputRaw)
			masterNode.InputPrintSnapshot_wcall()

		default:
		}
		// time.Sleep(100 * time.Millisecond)
		if err == io.EOF {
			break
		}
	}
	time.Sleep(100 * time.Hour)
}
