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

	"github.com/rtntubmt97/springprj/impl"
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
		impl.LogE("Invalid input")
		impl.LogE(err.Error())
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
	if impl.LoadedConfig.UseBin {
		exeCmd = exec.Command("bin/observer", configPath, "somethingtohtop")
	} else {
		exeCmd = exec.Command("go", "run", "app/observer/observer.go", configPath, "somethingtohtop")
	}
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr
	err := exeCmd.Start()
	if err != nil {
		impl.LogE(err.Error())
	}
}

// Start a node process. This is a real os process with a pid on OS.
func createNode(configPath string, id string, initMoney string) {
	var exeCmd *exec.Cmd
	if impl.LoadedConfig.UseBin {
		exeCmd = exec.Command("bin/node", configPath, id, initMoney, "somethingtohtop")
	} else {
		exeCmd = exec.Command("go", "run", "app/node/node.go", configPath, id, initMoney, "somethingtohtop")
	}
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr
	err := exeCmd.Start()
	if err != nil {
		impl.LogE(err.Error())
	}
}

// MasterObj will be used to send the command.
var masterObj impl.Master

// var observerNode observer.Observer

func main() {
	// Get the configuration path if it was passed to the argument command line.
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// Reload the configuration.
	err, _ := impl.ReloadConfig(configPath)
	if err != nil {
		fmt.Println("invalid config path")
	}

	// Check whether the master port is free.
	if !impl.IsPortAvailable(int(impl.MasterPort)) {
		fmt.Printf("Master cannot use port %d\n", impl.MasterPort)
		return
	}

	// Check whether the observer port is free.
	if !impl.IsPortAvailable(int(impl.ObserverPort)) {
		fmt.Printf("Observer cannot use port %d\n", impl.ObserverPort)
		return
	}

	// Check and load the input file.
	file, err := os.Open(impl.LoadedConfig.InputFile)
	if err != nil {
		impl.LogE(err.Error())
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
		if impl.LoadedConfig.PrintInput {
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
			impl.LogI("Matched StartMaster")
			masterObj = impl.Master{}
			masterObj.Init()
			go masterObj.Listen()
			time.Sleep(2 * time.Second)
			createObserver(configPath)
			time.Sleep(2 * time.Second)

		case KillAll:
			// Send the Kill signal to observer and nodes it has stated, then kill the
			// master itself.
			impl.LogI(inputRaw)
			masterObj.KillAll()
			time.Sleep(1 * time.Second)
			impl.LogI("Ready to exit")
			os.Exit(0)

		case CreateNode:
			// Start a node proccess with an id and start money specified by the input
			impl.LogI(inputRaw)
			createNode(configPath, input[1], input[2])
			time.Sleep(2000 * time.Millisecond)

		case Send:
			// Send Send singal to node to command it send the money to other node.
			impl.LogI(inputRaw)
			sender, _ := strconv.Atoi(input[1])
			receiver, _ := strconv.Atoi(input[2])
			money, _ := strconv.Atoi(input[3])
			masterObj.SignalSend(int32(sender), int32(receiver), int32(money))

		case Receive:
			// Send Seceive signal to node to command it receive the money from a sender channel.
			impl.LogI(inputRaw)
			sender := -1
			receiver := -1
			if len(input) == 2 {
				receiver, _ = strconv.Atoi(input[1])
			} else if len(input) == 3 {
				receiver, _ = strconv.Atoi(input[1])
				sender, _ = strconv.Atoi(input[2])
			}
			masterObj.SignalReceive(int32(receiver), int32(sender))

		case ReceiveAll:
			// Send ReceiveAll signal to all nodes to command them drain all the channels.
			impl.LogI(inputRaw)
			masterObj.SignalReceiveAll()

		case BeginSnapshot:
			// Send BeginSnapshot signal to a node to command it start the snapshot process.
			impl.LogI(inputRaw)
			startNodeId, _ := strconv.Atoi(input[1])
			masterObj.SignalBeginSnapshot(int32(startNodeId))

		case CollectState:
			// Send CollectState signal to the Observer to command it collect states from
			// all nodes.
			impl.LogI(inputRaw)
			masterObj.SignalCollectState()

		case PrintSnapshot:
			// Send PrintSnapshot signal to the Observer to command it print the states it has collected
			impl.LogI(inputRaw)
			masterObj.SignalPrintSnapshot()

		default:
		}
		// time.Sleep(100 * time.Millisecond)
		if err == io.EOF {
			break
		}
	}
	time.Sleep(100 * time.Hour)
}
