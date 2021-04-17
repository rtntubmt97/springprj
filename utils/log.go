package utils

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/rtntubmt97/springprj/define"
)

var innerLog *log.Logger

func init() {
	innerLog = log.New(os.Stdout, "", log.Lshortfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LogR(msg define.ProjectOutput) {
	fmt.Println(msg)
}

func LogE(msg string) {
	// innerLog.Output(2, fmt.Sprintf("[logE] %s\n", msg))
}

func LogD(msg string) {
	// innerLog.Output(2, fmt.Sprintf("[LogD] %s\n", msg))
}

func LogI(msg string) {
	// innerLog.Output(2, fmt.Sprintf("[LogI] %s\n", msg))
}

func PrintStack() {
	debug.PrintStack()
}

func CreateTransferOutput(sender int32, money int32) define.ProjectOutput {
	msg := fmt.Sprintf("%d Transfer %d", sender, money)
	return define.ProjectOutput(msg)
}

func CreateTokenOutput(sender int32) define.ProjectOutput {
	msg := fmt.Sprintf("%d SnapshotToken -1", sender)
	return define.ProjectOutput(msg)
}

func CreateBeginSnapshotOutput(beginer int32) define.ProjectOutput {
	msg := fmt.Sprintf("BeginSnapshot %d", beginer)
	return define.ProjectOutput(msg)
}

func CreateReceiveSnapshotOutput(beginer int32) define.ProjectOutput {
	msg := fmt.Sprintf("%d SnapshotToken -1", beginer)
	return define.ProjectOutput(msg)
}
