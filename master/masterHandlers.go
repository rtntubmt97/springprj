package master

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

func (master *Master) requestInfo_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI("requestInfo_whandle run")
	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Init(define.RequestInfoRsp)
	rspMsg.WriteI32(int32(len(master.connector.ConnectedConns)))
	for otherConnId, otherConn := range master.connector.ConnectedConns {
		rspMsg.WriteI32(otherConnId)
		portStr := strings.Split(otherConn.RemoteAddr().String(), ":")[0]
		port, _ := strconv.Atoi(portStr)
		rspMsg.WriteI32(int32(port))
	}
	conn := master.connector.ConnectedConns[connId]
	rspMsg.WriteMessage(conn)
	utils.LogI(fmt.Sprintf("Master responsed to node %d address %s", connId, conn.RemoteAddr().String()))
}
