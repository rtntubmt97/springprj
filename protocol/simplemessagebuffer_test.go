package protocol

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	fmt.Println("hi")
}

func TestReadWriteI32(t *testing.T) {
	mb := SimpleMessageBuffer{}
	mb.InitEmpty()
	in := int32(2147483)
	mb.WriteI32(in)
	out := mb.ReadI32()
	if in != out {
		t.Error("in does not equal out")
	}
}

func TestReadWriteI64(t *testing.T) {
	mb := SimpleMessageBuffer{}
	mb.InitEmpty()
	in := int64(2147483647123)
	mb.WriteI64(in)
	out := mb.ReadI64()
	if in != out {
		t.Error("in does not equal out")
	}
}

func TestReadWriteString(t *testing.T) {
	mb := SimpleMessageBuffer{}
	mb.InitEmpty()
	in := "abcd123asfsadf"
	mb.WriteString(in)
	out := mb.ReadString()
	if in != out {
		t.Error("in does not equal out")
	}
}

func TestReadWriteMessageBuffer(t *testing.T) {
	// t.Log("TestReadWriteMessageBuffer start")
	i1 := int32(123)
	i2 := int32(321)
	i3 := int64(1232142142344)
	s4 := "acsdf"
	s5 := "asfggsa"
	i6 := int64(998671)

	inMb1 := SimpleMessageBuffer{}
	inMb1.InitEmpty()
	inMb1.WriteI32(i1).
		WriteI32(i2).
		WriteI64(i3).
		WriteString(s4).
		WriteString(s5).
		WriteI64(i6)

	inMb2 := SimpleMessageBuffer{}
	inMb2.InitEmpty()
	inMb2.WriteI32(i1)

	inMb3 := SimpleMessageBuffer{}
	inMb3.InitEmpty()
	inMb3.WriteString(s5)

	stream := new(bytes.Buffer)
	inMb1.WriteMessage(stream)
	inMb2.WriteMessage(stream)
	inMb3.WriteMessage(stream)
	outMb1 := SimpleMessageBuffer{}
	err1 := outMb1.ReadMessage(stream)
	outMb2 := SimpleMessageBuffer{}
	err2 := outMb2.ReadMessage(stream)
	outMb3 := SimpleMessageBuffer{}
	outMb3.ReadMessage(stream)

	if err1 != nil {
		t.Error("cannot read message1")
	}

	if err2 != nil {
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

	if i1 != outMb2.ReadI32() {
		t.Error("wrong i1")
	}

	if s5 != outMb3.ReadString() {
		t.Error("wrong s5")
	}
}
