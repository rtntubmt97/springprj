package protocol

import (
	"bytes"
	"encoding/binary"
)

type MessageBuffer struct {
	Buf *bytes.Buffer
}

func (mb *MessageBuffer) InitEmpty() {
	mb.Buf = new(bytes.Buffer)
}

func (mb *MessageBuffer) Init(command int32, ver int32) {

}

func (mb *MessageBuffer) WriteI32(i int32) {
	binary.Write(mb.Buf, binary.BigEndian, i)
}

func (mb *MessageBuffer) WriteI64(i int64) {
	binary.Write(mb.Buf, binary.BigEndian, i)
}

func (mb *MessageBuffer) WriteString(s string) {
	sLen := int32(len(s))
	mb.WriteI32(sLen)
	mb.Buf.WriteString(s)
}

func (mb *MessageBuffer) ReadI32() int32 {
	var out int32
	binary.Read(mb.Buf, binary.BigEndian, &out)
	return out
}

func (mb *MessageBuffer) ReadI64() int64 {
	var out int64
	binary.Read(mb.Buf, binary.BigEndian, &out)
	return out
}

func (mb *MessageBuffer) ReadString() string {
	sLen := int(mb.ReadI32())
	return string(mb.Buf.Next(sLen))
}

// func (mb *MessageBuffer) SetReadyToGet() bool {
// 	return true
// }

// func (mb *MessageBuffer) SetReadyToWrite() bool {
// 	return true
// }
