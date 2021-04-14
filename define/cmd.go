package define

type Command int32

const (
	SendInt32 int32 = iota
	SendInt64
	SendString
)
