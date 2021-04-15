package master

import (
	connectorPkg "github.com/rtntubmt97/springprj/connector"
)

type Master struct {
	id        int32
	connector connectorPkg.Connector
}

func (master *Master) Listen(port int) {
	master.connector.Listen(port)
}

func (master *Master) Connect(id int32, port int32) {
	master.connector.Connect(id, port)
}
