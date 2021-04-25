// Package protocol contains serializing/deserializing data mechanics. It provides
// the structures which can be used to convert between in-program-data (string, int32, int64)
// and wired-format-data (bytes) in order exchange data through the network.
// It can be implemented in numerous formats (raw binary, json, xml ...) with
// different features (compressed, encrypted, readable, ...) as long as it satisfied
// the MessageBuffer, Writeable and Readable interfaces in the define package.

package impl

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

// +---------+---------+-----+...+-----+--------+...+--------+--------+...+--------+---------+-----+...+-----+
// |int32 i1 |int32 i2 |      ....     |    string length    |    string bytes     |int32 i2 |      ....     |
// +---------+---------+-----+...+-----+--------+...+--------+--------+...+--------+---------+-----+...+-----+

// For simplicity, the BinaryProtocol is the only implemented structure for transfering data.
// It contains an underlying Buf (bytes.Buffer) to store the writed data which later will be read. The reader
// must know the order of variable the writer writed data. BinaryProtocol behaves as an first-in-first-out
// queue, thus the read data will be removed from the buffer.
type BinaryProtocol struct {
	Buf *bytes.Buffer
}

// Initialize SimpleMessageBuffer with an empty Buf
func (mb *BinaryProtocol) InitEmpty() {
	mb.Buf = new(bytes.Buffer)
}

// Initialize SimpleMessageBuffer with an empty Buf, write an int32 to specify the
// the kind of message it stored. (This method is useful in current project but
// putting the method here is not a good pattern.)
func (mb *BinaryProtocol) Init(command ConnectorCmd) {
	mb.InitEmpty()
	mb.WriteI32(int32(command))
}

// Write int32 to the SimpleMessageBuffer
func (mb BinaryProtocol) WriteI32(i int32) MessageBuffer {
	binary.Write(mb.Buf, binary.BigEndian, i)
	return mb
}

// Write int64 to the SimpleMessageBuffer
func (mb BinaryProtocol) WriteI64(i int64) MessageBuffer {
	binary.Write(mb.Buf, binary.BigEndian, i)
	return mb
}

// Write string to the SimpleMessageBuffer
func (mb BinaryProtocol) WriteString(s string) MessageBuffer {
	sLen := int32(len(s))
	mb.WriteI32(sLen)
	mb.Buf.WriteString(s)
	return mb
}

// Read int32 from the SimpleMessageBuffer
func (mb BinaryProtocol) ReadI32() int32 {
	var out int32
	binary.Read(mb.Buf, binary.BigEndian, &out)
	return out
}

// Read int64 from the SimpleMessageBuffer
func (mb BinaryProtocol) ReadI64() int64 {
	var out int64
	binary.Read(mb.Buf, binary.BigEndian, &out)
	return out
}

// Read string from the SimpleMessageBuffer
func (mb BinaryProtocol) ReadString() string {
	sLen := int(mb.ReadI32())
	return string(mb.Buf.Next(sLen))
}

// While developing, reading data from a reader, or writing data to a writer, may get
// some mistakes or errors. Writing and reading magic bytes first will help to recognize
// them.
var magicBytes = []byte("xxDDxx")

// Read data from an io.Reader into the SimpleMessageBuffer, the data of the underlying Buf
// will be changed
func (message *BinaryProtocol) Read(reader io.Reader) error {
	initBytes := make([]byte, len(magicBytes))
	_, err := reader.Read(initBytes)
	if err != nil {
		LogE(err.Error())
		return err
	}

	if !reflect.DeepEqual(magicBytes, initBytes) {
		LogE("wrong initBytes")
		return ErrWrongInitBytes
	}

	var len int32
	binary.Read(reader, binary.BigEndian, &len)
	data := make([]byte, len)
	n, err := io.ReadFull(reader, data)
	if err != nil {
		LogE(err.Error())
		return err
	}
	if int32(n) != len {
		LogE("wrong len")
		return err
	}

	message.Buf = bytes.NewBuffer(data)
	return nil
}

// Write data from an io.Writer to the SimpleMessageBuffer, the data of the underlying Buf
// will not be changed
func (message *BinaryProtocol) Write(writer io.Writer) error {
	writer.Write(magicBytes)

	len := int32(message.Buf.Len())
	binary.Write(writer, binary.BigEndian, len)

	writer.Write(message.Buf.Bytes())
	return nil
}
