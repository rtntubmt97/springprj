package utils

import "fmt"

func LogR(msg interface{}) {
	fmt.Println(msg)
}

func LogE(msg interface{}) {
	fmt.Printf("[logE] %s\n", msg)
}
