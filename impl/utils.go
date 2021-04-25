package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
	"strconv"
)

type Config struct {
	IsProduction bool
	InputFile    string
	UseBin       bool
	UseLog       bool
	PrintInput   bool
	MasterId     int32
	MasterPort   int32
	ObserverId   int32
	ObserverPort int32
}

var LoadedConfig Config

func ReloadConfig(filePath string) (error, Config) {
	if filePath == "" {
		filePath = "config.json"
	}
	file, _ := os.Open(filePath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
		return err, config
	}

	LoadedConfig = config

	MasterId = LoadedConfig.MasterId
	MasterPort = LoadedConfig.MasterPort
	ObserverId = LoadedConfig.ObserverId
	ObserverPort = LoadedConfig.ObserverPort
	UseLog = LoadedConfig.UseLog

	return nil, config
}

func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return false
	}

	err = ln.Close()
	if err != nil {
		return false
	}

	return true
}

func GetAvailablePort(start int, end int) int {
	for port := start; port < end; port++ {
		if IsPortAvailable(port) {
			return port
		}
	}

	return -1
}

var innerLog *log.Logger
var UseLog bool

func init() {
	innerLog = log.New(os.Stdout, "", log.Lshortfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	UseLog = false
}

func LogR(msg ProjectOutput) {
	fmt.Println(msg)
}

func LogE(msg string) {
	if !UseLog {
		return
	}
	innerLog.Output(2, fmt.Sprintf("[logE] %s\n", msg))
}

func LogD(msg string) {
	if !UseLog {
		return
	}
	innerLog.Output(2, fmt.Sprintf("[LogD] %s\n", msg))
}

func LogI(msg string) {
	if !UseLog {
		return
	}
	innerLog.Output(2, fmt.Sprintf("[LogI] %s\n", msg))
}

func PrintStack() {
	debug.PrintStack()
}

func CreateTransferOutput(sender int32, money int32) ProjectOutput {
	msg := fmt.Sprintf("%d Transfer %d", sender, money)
	return ProjectOutput(msg)
}

func CreateBeginSnapshotOutput(beginer int32) ProjectOutput {
	msg := fmt.Sprintf("Started by Node %d", beginer)
	return ProjectOutput(msg)
}

func CreateReceiveSnapshotOutput(beginer int32) ProjectOutput {
	msg := fmt.Sprintf("%d SnapshotToken -1", beginer)
	return ProjectOutput(msg)
}
