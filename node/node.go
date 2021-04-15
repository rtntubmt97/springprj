package node

import (
	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
)

type Node struct {
	id        int32
	connector connectorPkg.Connector
	channels  map[int32](chan int32)
}

func (node *Node) Init(id int32) {
	node.id = id
	node.connector = connectorPkg.Connector{}
	node.connector.Init(id)

	node.channels = make(map[int32](chan int32))

	node.connector.SetHandleFunc(define.SendInt32, node.sendInt32_handle)
	node.connector.SetHandleFunc(define.SendInt64, node.sendInt64_handle)
	node.connector.SetHandleFunc(define.SendString, node.sendString_handle)
}

func (node *Node) Start() {

}

func (node *Node) Listen(port int) {
	node.connector.Listen(port)
}

func (node *Node) Connect(id int32, port int32) {
	node.connector.Connect(id, port)
}

func (node *Node) ConnectMaster() {

}
