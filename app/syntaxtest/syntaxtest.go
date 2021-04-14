package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/rtntubmt97/springprj/protocol"
)

func main() {
	strReader := strings.NewReader("aasfs")
	p := make([]byte, 2)
	n, err := io.ReadFull(strReader, p)
	if err != nil || n != 2 {
		fmt.Println("read fail")
	}
	fmt.Println(p)
	n, err = io.ReadFull(strReader, p)
	if err != nil || n != 2 {
		fmt.Println("read fail")
	}
	fmt.Println(p)

	var t *testing.T

	i1 := int32(123)
	i2 := int32(321)
	i3 := int64(1232142142344)
	s4 := "acsdf"
	s5 := "asfggsa"
	i6 := int64(998671)

	inMb1 := protocol.MessageBuffer{}
	inMb1.InitEmpty()
	inMb1.WriteI32(i1).
		WriteI32(i2).
		WriteI64(i3).
		WriteString(s4).
		WriteString(s5).
		WriteI64(i6)

	inMb2 := protocol.MessageBuffer{}
	inMb2.InitEmpty()
	inMb2.WriteI32(i1)

	stream := new(bytes.Buffer)
	fmt.Println(stream.Bytes())
	protocol.WriteMessage(stream, inMb1)
	fmt.Println(stream.Bytes())
	protocol.WriteMessage(stream, inMb2)
	fmt.Println(stream.Bytes())
	outMb1 := protocol.ReadMessage(stream)
	fmt.Println(stream.Bytes())
	outMb2 := protocol.ReadMessage(stream)
	fmt.Println(stream.Bytes())

	if outMb1 == nil {
		t.Error("cannot read message1")
	}

	if outMb2 == nil {
		t.Error("cannot read message2")
	}

	// temp := outMb2.ReadI32()
	// fmt.Println(temp)

	if i1 != outMb1.ReadI32() {
		t.Error("wrong i1")
	}
	if i2 != outMb1.ReadI32() {
		t.Error("wrong i2")
	}
	if i3 != outMb1.ReadI64() {
		t.Error("wrong i3")
	}
	if s4 != outMb1.ReadString() {
		t.Error("wrong s4")
	}
	if s5 != outMb1.ReadString() {
		t.Error("wrong s5")
	}
	if i6 != outMb1.ReadI64() {
		t.Error("wrong i6")
	}

}
