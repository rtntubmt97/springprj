// Observer is another actor of this project. It receives the command signal from
// the master and then do its job. It receives CollectState or PrintSnapshot signal

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

// Observer structure contains its id, an connector, and nodeId->nodeSnapshot map
type Observer struct {
	id        int32
	connector connectorPkg.Connector
	snapShots map[int32]node.SnapShot
}

// Initilize the Observer
func (observer *Observer) Init() {
	observer.id = define.ObserverId
	observer.connector = connectorPkg.Connector{}
	observer.connector.Init(define.ObserverId)
	observer.connector.ParticipantType = connectorPkg.ObserverType

	observer.connector.SetHandleFunc(define.Input_Kill, observer.inputKillHandle)
	observer.connector.SetHandleFunc(define.Input_CollectState, observer.inputCollectState_whandle)
	observer.connector.SetHandleFunc(define.Input_PrintSnapshot, observer.inputPrintSnapshotHandle)

	observer.snapShots = make(map[int32]node.SnapShot)
}

// Start the listen operation.
func (observer *Observer) Listen() {
	observer.connector.Listen(int(define.ObserverPort))
}

// connect to a connector (master, observer or node).
func (observer *Observer) Connect(id int32, port int32) {
	observer.connector.Connect(id, port)
}

// Connect to master.
func (observer *Observer) ConnectMaster() {
	observer.Connect(define.MasterId, define.MasterPort)
}

// Handle the Kill signal from the master.
func (observer *Observer) inputKillHandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received kill signal", observer.id))
	observer.connector.SendAckRsp(define.MasterId, define.Input_KillRsp)
	os.Exit(0)
}

// Handle the CollectState signal from the master. Observer will send the collect state
// cmd to all nodes and get their snapshots, then save it to the observer snapshot map.
func (observer *Observer) inputCollectState_whandle(connId int32, msg define.MessageBuffer) {
	utils.LogI(fmt.Sprintf("Node %d Received inputCollectState signal", observer.id))

	for peerId := range observer.connector.ConnectedConns {
		if peerId == define.MasterId || peerId == define.ObserverId {
			continue
		}
		snapShot := observer.collectStateHandle(peerId)
		observer.snapShots[peerId] = snapShot
	}

	observer.connector.SendAckRsp(define.MasterId, define.Input_CollectStateRsp)
}

// Send the collect state cmd to a node specified by an id.
func (observer *Observer) collectStateHandle(connId int32) node.SnapShot {
	msg := protocol.BinaryProtocol{}
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

	return node.SnapShot{NodeMoney: money, ChannelMoneys: channels}
}

// Handle the PrintSnapshot from the master. Observer will format the snapshots which were
// saved in the observer snapshot map, then print it.
func (observer *Observer) inputPrintSnapshotHandle(connId int32, msg define.MessageBuffer) {
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

	observer.connector.SendAckRsp(define.MasterId, define.Input_PrintSnapshotRsp)
}

// Create NodeState as the specification required
func CreateNodeState(nodeId int32, money int64) string {
	return fmt.Sprintf("node %d = %d\n", nodeId, money)
}

// Create ChannelState as the specification required
func CreateChannelState(sender int32, receiver int32, money int64) string {
	return fmt.Sprintf("channel (%d → %d) = %d\n", sender, receiver, money)
}
