package protocol

import (
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
