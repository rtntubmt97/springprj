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

// Get the input from Standard input (usualy a keyboard), split it by space character and
// return its as an array of string
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

// Start an observer process. This is a real os process with a pid on OS.
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

// Start a node process. This is a real os process with a pid on OS.
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

// MasterObj will be used to send the command.
var masterObj master.Master

// var observerNode observer.Observer

func main() {
	// Get the configuration path if it was passed to the argument command line.
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// Reload the configuration.
	err, _ := utils.ReloadConfig(configPath)
	if err != nil {
		fmt.Println("invalid config path")
	}

	// Check whether the master port is free.
	if !utils.IsPortAvailable(int(define.MasterPort)) {
		fmt.Printf("Master cannot use port %d\n", define.MasterPort)
		return
	}

	// Check whether the observer port is free.
	if !utils.IsPortAvailable(int(define.ObserverPort)) {
		fmt.Printf("Observer cannot use port %d\n", define.ObserverPort)
		return
	}

	// Check and load the input file.
	file, err := os.Open(utils.LoadedConfig.InputFile)
	if err != nil {
		utils.LogE(err.Error())
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		// Uncomment the following line if you want the master program receive input from the stadard input.
		// input := getStdinInput()
		inputRaw, err := reader.ReadString('\n')

		// Check whether input is readable.
		if err != io.EOF {
			if inputRaw[0] == ';' {
				continue
			}
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		// Print the input line if config allows.
		if utils.LoadedConfig.PrintInput {
			fmt.Print(inputRaw)
		}

		// Split the input by space in order to process it in the future.
		input := strings.Split(inputRaw, " ")
		for i, ele := range input {
			input[i] = strings.Trim(ele, "\n\r\t ")
		}

		// Get the input command to use the switch case statement below.
		cmd := InputCmd(input[0])

		switch cmd {
		case StartMaster:
			// Start the master obj and observer proccess.
			utils.LogI("Matched StartMaster")
			masterObj = master.Master{}
			masterObj.Init()
			go masterObj.Listen()
			time.Sleep(2 * time.Second)
			createObserver(configPath)
			time.Sleep(2 * time.Second)

		case KillAll:
			// Send the Kill signal to observer and nodes it has stated, then kill the
			// master itself.
			utils.LogI(inputRaw)
			masterObj.KillAll()
			time.Sleep(1 * time.Second)
			utils.LogI("Ready to exit")
			os.Exit(0)

		case CreateNode:
			// Start a node proccess with an id and start money specified by the input
			utils.LogI(inputRaw)
			createNode(configPath, input[1], input[2])
			time.Sleep(2000 * time.Millisecond)

		case Send:
			// Send Send singal to node to command it send the money to other node.
			utils.LogI(inputRaw)
			sender, _ := strconv.Atoi(input[1])
			receiver, _ := strconv.Atoi(input[2])
			money, _ := strconv.Atoi(input[3])
			masterObj.InputSend(int32(sender), int32(receiver), int32(money))

		case Receive:
			// Send Seceive signal to node to command it receive the money from a sender channel.
			utils.LogI(inputRaw)
			sender := -1
			receiver := -1
			if len(input) == 2 {
				receiver, _ = strconv.Atoi(input[1])
			} else if len(input) == 3 {
				receiver, _ = strconv.Atoi(input[1])
				sender, _ = strconv.Atoi(input[2])
			}
			masterObj.InputReceive(int32(receiver), int32(sender))

		case ReceiveAll:
			// Send ReceiveAll signal to all nodes to command them drain all the channels.
			utils.LogI(inputRaw)
			masterObj.InputReceiveAll()

		case BeginSnapshot:
			// Send BeginSnapshot signal to a node to command it start the snapshot process.
			utils.LogI(inputRaw)
			startNodeId, _ := strconv.Atoi(input[1])
			masterObj.InputBeginSnapshot(int32(startNodeId))

		case CollectState:
			// Send CollectState signal to the Observer to command it collect states from
			// all nodes.
			utils.LogI(inputRaw)
			masterObj.InputCollectState()

		case PrintSnapshot:
			// Send PrintSnapshot signal to the Observer to command it print the states it has collected
			utils.LogI(inputRaw)
			masterObj.InputPrintSnapshot()

		default:
		}
		// time.Sleep(100 * time.Millisecond)
		if err == io.EOF {
			break
		}
	}
	time.Sleep(100 * time.Hour)
}
