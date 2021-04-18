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

type Wraper struct {
	i int32
}

func main() {
	m := make(map[int32]bool)
	m[0] = false
	delete(m, 0)
	a, exist := m[0]
	fmt.Println(a)
	fmt.Println(exist)
}
