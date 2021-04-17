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
	m := make(map[int32]*Wraper)
	w := Wraper{i: 1}
	m[111] = &w
	w.i = 10
	fmt.Println(m[111].i)
	for _, i := range m {
		fmt.Println(i)
	}
}
