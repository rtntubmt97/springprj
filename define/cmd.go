package define

const (
	SendInt32 int32 = iota
	SendInt64
	SendString
	Greeting
	GreetingRsp
	StartMaster
	KillAll
	CreateNode
	Send
	Receive
	ReceiveAll
	BeginSnapshot
	CollectState
	PrintSnapshot
)
