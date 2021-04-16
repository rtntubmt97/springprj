package define

type MessageBuffer interface {
	WriteI32(i int32) MessageBuffer

	WriteI64(i int64) MessageBuffer

	WriteString(s string) MessageBuffer

	ReadI32() int32

	ReadI64() int64

	ReadString() string
}

type HandleFunc func(connId int32, msg MessageBuffer)
