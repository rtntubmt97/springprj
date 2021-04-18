package observer

import (
	"fmt"
	"os"
	"sort"
	"strings"

	connectorPkg "github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/define"
	"github.com/rtntubmt97/springprj/node"
	"github.com/rtntubmt97/springprj/protocol"
	"github.com/rtntubmt97/springprj/utils"
)

type Observer struct {
	id        int32
	connector connectorPkg.Connector
	snapShots map[int32]node.SnapShot
}

func (observer *Observer) Init() {
	observer.id = define.ObserverId
	observer.connector = connectorPkg.Connector{}
	observer.connector.Init(define.ObserverId)
	observer.connector.ParticipantType = connectorPkg.ObserverType

	observer.connector.SetHandleFunc(define.Input_Kill, observer.kill_handle)
	observer.connector.SetHandleFunc(define.Input_CollectState, observer.inputCollectState_whandle)
	observer.connector.SetHandleFunc(define.Input_PrintSnapshot, observer.inputPrintSnapshot_handle)

	observer.snapShots = make(map[int32]node.SnapShot)
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
	utils.LogI(fmt.Sprintf("Node %d Received kill signal", observer.id))
	os.Exit(0)
}

func (observer *Observer) inputPrintSnapshot_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputPrintSnapshot signal", observer.id))

	ret := strings.Builder{}

	peerIds := make([]int, 0, len(observer.snapShots))
	for k := range observer.snapShots {
		peerIds = append(peerIds, int(k))
	}
	sort.Ints(peerIds)

	ret.WriteString("---Node states\n")
	for _, peerId := range peerIds {
		snapshot := observer.snapShots[int32(peerId)]
		ret.WriteString(CreateNodeState(int32(peerId), snapshot.NodeMoney))
	}

	ret.WriteString("---Channel states\n")
	for _, sender := range peerIds {
		for _, receiver := range peerIds {
			if sender == receiver {
				continue
			}
			snapshot := observer.snapShots[int32(receiver)]
			channelMoney := snapshot.ChannelMoneys[int32(sender)]
			channelState := CreateChannelState(int32(sender), int32(receiver), channelMoney)
			ret.WriteString(channelState)
		}
	}
	utils.LogR(define.ProjectOutput(ret.String()))

}

func (observer *Observer) inputCollectState_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputCollectState signal", observer.id))

	for peerId := range observer.connector.ConnectedConns {
		if peerId == define.MasterId || peerId == define.ObserverId {
			continue
		}
		snapShot := observer.collectState_wcall(peerId)
		observer.snapShots[peerId] = snapShot
	}

	rspMsg := protocol.SimpleMessageBuffer{}
	rspMsg.Init(define.Rsp)
	rspMsg.WriteI32(int32(define.Input_CollectStateRsp))
	observer.connector.WriteTo(define.MasterId, &rspMsg)
}

func (observer *Observer) collectState_wcall(connId int32) node.SnapShot {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.CollectState)
	observer.connector.WriteTo(connId, &msg)

	utils.LogI(fmt.Sprintf("Node %d sent collectState to node %d", observer.id, connId))

	rspMsg := observer.connector.WaitRsp(connId)
	cmd := define.ConnectorCmd(rspMsg.ReadI32())
	if cmd != define.CollectStateRsp {
		utils.LogE("Wrong collectState response")
	}
	utils.LogI("Correct collectState response")
	money := rspMsg.ReadI64()
	channelsLen := rspMsg.ReadI32()
	channels := make(map[int32]int64)
	for i := int32(0); i < channelsLen; i++ {
		channelId := rspMsg.ReadI32()
		channelMoney := rspMsg.ReadI64()
		channels[channelId] = channelMoney
	}

	utils.LogI(fmt.Sprintf("Observer finished collectState_wcall from node %d", connId))
	fmt.Println(node.SnapShot{NodeMoney: money, ChannelMoneys: channels})

	return node.SnapShot{NodeMoney: money, ChannelMoneys: channels}
}

func CreateNodeState(nodeId int32, money int64) string {
	return fmt.Sprintf("node %d = %d\n", nodeId, money)
}

func CreateChannelState(sender int32, reciever int32, money int64) string {
	return fmt.Sprintf("channel (%d → %d) = %d\n", sender, reciever, money)
}
