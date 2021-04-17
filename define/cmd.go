package define

type ConnectorCmd int32

const (
	SendInt32 ConnectorCmd = iota
	SendInt64
	SendString
	Rsp
	Greeting
	GreetingRsp
	RequestInfo
	RequestInfoRsp
	StartMaster
	KillAll
	CreateNode
	Send
	Receive
	ReceiveAll
	BeginSnapshot
	CollectState
	PrintSnapshot
	Input_Send
	Input_Kill
)
