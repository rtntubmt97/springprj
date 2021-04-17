package utils

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

var innerLog *log.Logger

func init() {
	innerLog = log.New(os.Stdout, "", log.Lshortfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LogR(msg interface{}) {
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
