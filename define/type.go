package define

import "io"

// Interface of the MessageBuffer

type MessageBuffer interface {
	WriteI32(i int32) MessageBuffer

	WriteI64(i int64) MessageBuffer

	WriteString(s string) MessageBuffer

	ReadI32() int32

	ReadI64() int64

	ReadString() string
}

type Writeable interface {
	Write(writer io.Writer) error
}

type Readable interface {
	Read(reader io.Reader) error
}

type HandleFunc func(connId int32, msg MessageBuffer)
