package master

import (
	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (master *Master) requestInfo_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI("requestInfo_whandle run")
	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Init(define.Rsp)
	rspMsg.WriteI32(int32(define.RequestInfoRsp))
	rspMsg.WriteI32(int32(len(master.connector.ConnectedConns)))
	for otherConnId := range master.connector.ConnectedConns {
		rspMsg.WriteI32(otherConnId)
		port := master.connector.OtherInfos[otherConnId].ListenPort
		rspMsg.WriteI32(int32(port))
	}
	conn := master.connector.ConnectedConns[connId]
	rspMsg.Write(conn)
}
