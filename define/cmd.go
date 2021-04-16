package define

type ConnectorCmd int32

const (
	SendInt32 ConnectorCmd = iota
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
