package utils

import (
	"fmt"
	"log"
	"runtime/debug"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LogR(msg interface{}) {
	fmt.Println(msg)
}

func LogE(msg string) {
	fmt.Printf("[logE] %s\n", msg)
}

func LogD(msg string) {
	fmt.Printf("[LogD] %s\n", msg)
}

func LogI(msg string) {
	fmt.Printf("[LogI] %s\n", msg)
}

func PrintStack() {
	debug.PrintStack()
}
