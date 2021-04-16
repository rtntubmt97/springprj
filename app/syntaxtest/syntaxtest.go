package main

import (
	"fmt"
)

type TestInterface interface {
	Print(s string)
}

type TestStruct struct {
}

func (st TestStruct) Print(s string) {
	fmt.Println("foo")
}

func foo(obj TestInterface) {
	obj.Print("foo")
}

func main() {
	obj := TestStruct{}
	foo(obj)
}
