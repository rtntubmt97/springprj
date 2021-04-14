package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type MessageBuffer struct {
	Buf *bytes.Buffer
}

func (mb *MessageBuffer) InitEmpty() {
	mb.Buf = new(bytes.Buffer)
}

func (mb *MessageBuffer) Init(command int32, ver int32) {

}

func (mb *MessageBuffer) WriteI32(i int32) *MessageBuffer {
	binary.Write(mb.Buf, binary.BigEndian, i)
	return mb
}

func (mb *MessageBuffer) WriteI64(i int64) *MessageBuffer {
	binary.Write(mb.Buf, binary.BigEndian, i)
	return mb
}

func (mb *MessageBuffer) WriteString(s string) *MessageBuffer {
	sLen := int32(len(s))
	mb.WriteI32(sLen)
	mb.Buf.WriteString(s)
	return mb
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

const magicByte byte = 'x'

func ReadMessage(reader io.Reader) *MessageBuffer {
	bufReader := bufio.NewReader(reader)
	initByte, err := bufReader.ReadByte()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if initByte != magicByte {
		fmt.Println("wrong initByte")
		return nil
	}

	var len int32
	binary.Read(bufReader, binary.BigEndian, &len)
	// data, err := ioutil.ReadAll(io.LimitReader(reader, int64(len)))
	data := make([]byte, len)
	n, err := io.ReadFull(bufReader, data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if int32(n) != len {
		fmt.Println("wrong len")
		return nil
	}

	return &MessageBuffer{Buf: bytes.NewBuffer(data)}

}

func WriteMessage(writer io.Writer, message MessageBuffer) {
	writer.Write([]byte{magicByte})

	len := int32(message.Buf.Len())
	binary.Write(writer, binary.BigEndian, len)

	writer.Write(message.Buf.Bytes())
}

// func (mb *MessageBuffer) SetReadyToGet() bool {
// 	return true
// }

// func (mb *MessageBuffer) SetReadyToWrite() bool {
// 	return true
// }
