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
	mb := new(MessageBuffer)
	mb.InitEmpty()
	in := int32(2147483)
	mb.WriteI32(in)
	out := mb.ReadI32()
	if in != out {
		t.Error("in does not equal out")
	}
}

func TestReadWriteI64(t *testing.T) {
	mb := new(MessageBuffer)
	mb.InitEmpty()
	in := int64(2147483647123)
	mb.WriteI64(in)
	out := mb.ReadI64()
	if in != out {
		t.Error("in does not equal out")
	}
}

func TestReadWriteString(t *testing.T) {
	mb := new(MessageBuffer)
	mb.InitEmpty()
	in := "abcd123asfsadf"
	mb.WriteString(in)
	out := mb.ReadString()
	if in != out {
		t.Error("in does not equal out")
	}
}

func TestReadWriteMessageBuffer(t *testing.T) {
	inMb := new(MessageBuffer)
	inMb.InitEmpty()

	i1 := int32(123)
	i2 := int32(321)
	i3 := int64(1232142142344)
	s4 := "acsdf"
	s5 := "asfggsa"
	i6 := int64(998671)

	inMb.WriteI32(i1).
		WriteI32(i2).
		WriteI64(i3).
		WriteString(s4).
		WriteString(s5).
		WriteI64(i6)

	stream := new(bytes.Buffer)
	WriteMessage(stream, *inMb)
	outMb := ReadMessage(stream)

	if i1 != outMb.ReadI32() {
		t.Error("wrong i1")
	}
	if i2 != outMb.ReadI32() {
		t.Error("wrong i2")
	}
	if i3 != outMb.ReadI64() {
		t.Error("wrong i3")
	}
	if s4 != outMb.ReadString() {
		t.Error("wrong s4")
	}
	if s5 != outMb.ReadString() {
		t.Error("wrong s5")
	}
	if i6 != outMb.ReadI64() {
		t.Error("wrong i6")
	}
}
