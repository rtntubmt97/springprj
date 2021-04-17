package observer

import (
	"fmt"
	"os"

	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/utils"
)

type Observer struct {
	id        int32
	connector connectorPkg.Connector
}

func (observer *Observer) Init() {
	observer.id = define.ObserverId
	observer.connector = connectorPkg.Connector{}
	observer.connector.Init(define.ObserverId)
	observer.connector.ParticipantType = connectorPkg.ObserverType

	observer.connector.SetHandleFunc(define.Input_Kill, observer.kill_handle)
}

func (observer *Observer) Listen() {
	observer.connector.Listen(int(define.ObserverPort))
}

func (observer *Observer) Connect(id int32, port int32) {
	observer.connector.Connect(id, port)
}

func (observer *Observer) ConnectMaster() {
	observer.Connect(define.MasterId, define.MasterPort)
}

func (observer *Observer) kill_handle(connId int32, msg define.MessageBuffer) {
	// fmt.Println("SendString_handle run!")
	utils.LogI(fmt.Sprintf("Node %d Received kill_handle signal", observer.id))
	os.Exit(0)
}
