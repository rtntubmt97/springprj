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
	observer.connector.SetHandleFunc(define.Input_CollectState, observer.inputCollectState_handle)
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

	ret.WriteString("---Node states")
	for _, peerId := range peerIds {
		snapshot := observer.snapShots[int32(peerId)]
		ret.WriteString(CreateNodeState(int32(peerId), snapshot.Money))
	}

	ret.WriteString("---Channel states")
	for _, peerId := range peerIds {
		snapshot := observer.snapShots[int32(peerId)]
		channels := snapshot.Channels
		channelIds := make([]int, 0, len(channels))
		for k := range channels {
			channelIds = append(channelIds, int(k))
		}
		sort.Ints(channelIds)
		for _, channelId := range channelIds {
			channel := channels[int32(channelId)]
			for _, money := range channel {
				channelState := CreateChannelState(int32(channelId), int32(peerId), money)
				ret.WriteString(channelState)
			}
		}
	}
	utils.LogR(define.ProjectOutput(ret.String()))

}

func (observer *Observer) inputCollectState_handle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputCollectState signal", observer.id))

	for peerId := range observer.connector.ConnectedConns {
		if peerId == define.MasterId || peerId == define.ObserverId {
			continue
		}
		snapShot := observer.collectState_wcall(peerId)
		observer.snapShots[peerId] = snapShot
	}

}

func (observer *Observer) collectState_wcall(connId int32) node.SnapShot {
	msg := protocol.SimpleMessageBuffer{}
	msg.Init(define.CollectState)
	observer.connector.WriteTo(connId, &msg)

	utils.LogI(fmt.Sprintf("Node %d sent collectState to node %d", observer.id, connId))

	rspMsg := observer.connector.WaitRsp(connId)
	cmd := define.ConnectorCmd(rspMsg.ReadI32())
	if cmd != define.CollectStateRsp {
		utils.LogE("Wrong response")
	}
	utils.LogI("Correct response")
	money := rspMsg.ReadI64()
	channelsLen := rspMsg.ReadI32()
	channels := make(map[int32][]int32)
	for i := int32(0); i < channelsLen; i++ {
		channelId := rspMsg.ReadI32()
		channelLen := rspMsg.ReadI32()
		channel := make([]int32, 0)
		for j := int32(0); j < channelLen; j++ {
			channel = append(channel, rspMsg.ReadI32())
		}
		channels[channelId] = channel
	}

	utils.LogI("done collectState_wcall")

	return node.SnapShot{Money: money, Channels: channels}
}

func CreateNodeState(nodeId int32, money int64) string {
	return fmt.Sprintf("node %d = %d\n", nodeId, money)
}

func CreateChannelState(sender int32, reciever int32, money int32) string {
	return fmt.Sprintf("channel (%d → %d) = %d", sender, reciever, money)
}
