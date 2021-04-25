// In this project, master only send the signal to observer and nodes to do action.
// Master can send Kill signal to a node, an observer, or all nodes, to force them exit.
// Master can send Send signal to a node with reciver and money data (then that node
// will send the money to the receiver).
// Master can send Receive signal to a node to make it take money in the money channel
// and add that them to its current money.
// Similarly, master can send the ReceiveAll, BeginSnapShot, CollectState, PrintSnapshot
// signal to observer/nodes to make them do their own business.

// Master package includes 3 files: master.go, masterCallers.go and masterHandlers.go.
// Master.go file contains the master structure and its basic operation.
// masterCallers.go file contains its caller which will be used to send the signal
// to nodes/observer by sending specific message.
// masterHandlers.go file contains its handler which will be used to handle the incoming messages.

package master

import (
	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
)

// Master struct, it contain an id (which will be assign to an defined int32) and an connector
type Master struct {
	id        int32
	connector connectorPkg.Connector
}

// Initilize master
func (master *Master) Init() {
	master.id = define.MasterId
	master.connector = connectorPkg.Connector{}
	master.connector.Init(define.MasterId)
	master.connector.ParticipantType = connectorPkg.MasterType

	master.connector.SetHandleFunc(define.RequestInfo, master.requestInfoHandle)
}

// Start the listen operation of master
func (master *Master) Listen() {
	master.connector.Listen(int(define.MasterPort))
}

// Connect to a connector by id and port
func (master *Master) Connect(id int32, port int32) {
	master.connector.Connect(id, port)
}

// Send Kill signal to all connector it known
func (master *Master) KillAll() {
	for nodeId := range master.connector.ConnectedConns {
		master.InputKill(nodeId)
	}

}
