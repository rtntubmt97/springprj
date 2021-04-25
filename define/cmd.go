package define

// The connector command integer will be written in front of every message in order
// to help the receiver know what to do with that message

type ConnectorCmd int32

// iota is bad for converse int to name by visual
const (
	SendInt32     ConnectorCmd = 111
	SendInt64     ConnectorCmd = 112
	SendString    ConnectorCmd = 113
	Rsp           ConnectorCmd = 114
	Greeting      ConnectorCmd = 115
	RequestInfo   ConnectorCmd = 116
	CreateNode    ConnectorCmd = 117
	Send          ConnectorCmd = 118
	SendToken     ConnectorCmd = 119
	CollectState  ConnectorCmd = 1110
	BeginSnapshot ConnectorCmd = 1111

	GreetingRsp     ConnectorCmd = 125
	RequestInfoRsp  ConnectorCmd = 126
	SendRsp         ConnectorCmd = 128
	SendTokenRsp    ConnectorCmd = 129
	CollectStateRsp ConnectorCmd = 1210

	Input_Send          ConnectorCmd = 211
	Input_receive       ConnectorCmd = 212
	Input_receiveAll    ConnectorCmd = 213
	Input_Kill          ConnectorCmd = 214
	Input_BeginSnapshot ConnectorCmd = 215
	Input_CollectState  ConnectorCmd = 216
	Input_PrintSnapshot ConnectorCmd = 217

	Input_SendRsp          ConnectorCmd = 221
	Input_receiveRsp       ConnectorCmd = 222
	Input_receiveAllRsp    ConnectorCmd = 223
	Input_KillRsp          ConnectorCmd = 224
	Input_BeginSnapshotRsp ConnectorCmd = 225
	Input_CollectStateRsp  ConnectorCmd = 226
	Input_PrintSnapshotRsp ConnectorCmd = 227
)
