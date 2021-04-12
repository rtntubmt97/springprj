package protocol

type ByteMessage struct {
	buf []byte
}

func (*ByteMessage) Init(command int, ver int) {

}

func (*ByteMessage) AppendI32(i int32) {

}
