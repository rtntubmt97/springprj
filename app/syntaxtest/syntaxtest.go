package main

import (
	"fmt"

	"github.com/rtntubmt97/springprj/node"
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
	info := node.MoneyTokenInfo{}
	info.Money = 100
	fmt.Println(info.IsToken())
	info2 := new(node.MoneyTokenInfo)
	info2.Money = -1
	fmt.Println(info2.IsToken())
}
