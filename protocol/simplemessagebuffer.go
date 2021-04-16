package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/utils"
)

type SimpleMessageBuffer struct {
	Buf *bytes.Buffer
}

func (mb *SimpleMessageBuffer) InitEmpty() {
	mb.Buf = new(bytes.Buffer)
}

func (mb *SimpleMessageBuffer) Init(command define.ConnectorCmd) {
	mb.InitEmpty()
	mb.WriteI32(int32(command))
}

func (mb SimpleMessageBuffer) WriteI32(i int32) define.MessageBuffer {
	binary.Write(mb.Buf, binary.BigEndian, i)
	return mb
}

func (mb SimpleMessageBuffer) WriteI64(i int64) define.MessageBuffer {
	binary.Write(mb.Buf, binary.BigEndian, i)
	return mb
}

func (mb SimpleMessageBuffer) WriteString(s string) define.MessageBuffer {
	sLen := int32(len(s))
	mb.WriteI32(sLen)
	mb.Buf.WriteString(s)
	return mb
}

func (mb SimpleMessageBuffer) ReadI32() int32 {
	var out int32
	binary.Read(mb.Buf, binary.BigEndian, &out)
	return out
}

func (mb SimpleMessageBuffer) ReadI64() int64 {
	var out int64
	binary.Read(mb.Buf, binary.BigEndian, &out)
	return out
}

func (mb SimpleMessageBuffer) ReadString() string {
	sLen := int(mb.ReadI32())
	return string(mb.Buf.Next(sLen))
}

var magicBytes = []byte("xxDDxx")

func (message *SimpleMessageBuffer) ReadMessage(reader io.Reader) error {
	// bufReader := bufio.NewReader(reader)
	initBytes := make([]byte, len(magicBytes))
	_, err := reader.Read(initBytes)
	if err != nil {
		utils.LogE(err.Error())
		return err
	}

	if !reflect.DeepEqual(magicBytes, initBytes) {
		utils.LogE("wrong initBytes")
		return err
	}

	var len int32
	binary.Read(reader, binary.BigEndian, &len)
	// data, err := ioutil.ReadAll(io.LimitReader(reader, int64(len)))
	data := make([]byte, len)
	n, err := io.ReadFull(reader, data)
	if err != nil {
		utils.LogE(err.Error())
		return err
	}
	if int32(n) != len {
		utils.LogE("wrong len")
		return err
	}

	message.Buf = bytes.NewBuffer(data)
	return nil
}

func (message *SimpleMessageBuffer) WriteMessage(writer io.Writer) error {
	writer.Write(magicBytes)

	len := int32(message.Buf.Len())
	binary.Write(writer, binary.BigEndian, len)

	writer.Write(message.Buf.Bytes())
	return nil
}

// func (mb *MessageBuffer) SetReadyToGet() bool {
// 	return true
// }

// func (mb *MessageBuffer) SetReadyToWrite() bool {
// 	return true
// }