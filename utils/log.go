package utils

import "fmt"

func LogR(msg interface{}) {
	fmt.Println(msg)
}

func LogE(msg string) {
	fmt.Printf("[logE] %s\n", msg)
}

func LogI(msg string) {
	fmt.Printf("[LogI] %s\n", msg)
}
