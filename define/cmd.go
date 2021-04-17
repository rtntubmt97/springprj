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
	// StartMaster
	// KillAll
	CreateNode
	Send
	SendRsp
	// Receive
	// ReceiveAll
	// BeginSnapshot
	SendToken
	CollectState
	CollectStateRsp
	BeginSnapshot
	Input_Send
	Input_Recieve
	Input_RecieveAll
	Input_Kill
	Input_BeginSnapshot
	Input_CollectState
	Input_PrintSnapshot
)
