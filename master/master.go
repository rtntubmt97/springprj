package master

import (
	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
)

type Master struct {
	id        int32
	connector connectorPkg.Connector
}

func (master *Master) Init() {
	master.id = define.MasterId
	master.connector = connectorPkg.Connector{}
	master.connector.Init(define.MasterId)
	master.connector.ParticipantType = connectorPkg.MasterType

	master.connector.SetHandleFunc(define.RequestInfo, master.requestInfo_whandle)
}

func (master *Master) Listen() {
	master.connector.Listen(int(define.MasterPort))
}

func (master *Master) Connect(id int32, port int32) {
	master.connector.Connect(id, port)
}

func (master *Master) KillAll() {
	for nodeId := range master.connector.ConnectedConns {
		master.InputKill_wcall(nodeId)
	}

}
