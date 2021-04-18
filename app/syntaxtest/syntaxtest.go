package main

import (
	"fmt"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/utils"
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
	_, config := utils.ReloadConfig("")
	fmt.Println(config.IsProduction)
	fmt.Println(define.MasterId)
}
