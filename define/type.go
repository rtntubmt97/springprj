package define

import "github.com/rtntubmt97/springprj/protocol"

type HandleFunc func(connId int32, msg protocol.MessageBuffer)
