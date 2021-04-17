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
	innerLog.Output(2, fmt.Sprintf("[logE] %s\n", msg))
	// fmt.Printf("[logE] %s\n", msg)
}

func LogD(msg string) {
	innerLog.Output(2, fmt.Sprintf("[LogD] %s\n", msg))
	// fmt.Printf("[LogD] %s\n", msg)
}

func LogI(msg string) {
	innerLog.Output(2, fmt.Sprintf("[LogI] %s\n", msg))
	// fmt.Printf("[LogI] %s\n", msg)
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
