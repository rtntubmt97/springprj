package main

import (
	"fmt"

	"github.com/rtntubmt97/springprj/protocol"
)

func main() {
	mb := new(protocol.MessageBuffer)
	mb.InitEmpty()
	in := int32(3)
	mb.WriteI32(in)
	out := mb.ReadI32()
	fmt.Println(out)
	if in != out {
		fmt.Println("in does not equal out")
	}
	mb.WriteString("asfsadfsas \nsfksal;f; /asd'''\"\"")
	outS := mb.ReadString()
	fmt.Println(outS)
	fmt.Println("asfsadfsas \nsfksal;f; /asd")
}
