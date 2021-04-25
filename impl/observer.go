// Observer is another actor of this project. It receives the command signal from
// the master and then do its job. It receives CollectState or PrintSnapshot signal

package impl

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Observer structure contains its id, an connector, and nodeId->nodeSnapshot map
type Observer struct {
	id        int32
	connector Connector
	snapShots map[int32]SnapShot
}

// Initilize the Observer
func (observer *Observer) Init() {
	observer.id = ObserverId
	observer.connector = Connector{}
	observer.connector.Init(ObserverId)
	observer.connector.ParticipantType = ObserverType

	observer.connector.SetHandleFunc(Input_Kill, observer.inputKillHandle)
	observer.connector.SetHandleFunc(Input_CollectState, observer.inputCollectState_whandle)
	observer.connector.SetHandleFunc(Input_PrintSnapshot, observer.signalPrintSnapshotHandle)

	observer.snapShots = make(map[int32]SnapShot)
}

// Start the listen operation.
func (observer *Observer) Listen() {
	observer.connector.Listen(int(ObserverPort))
}

// connect to a connector (master, observer or node).
func (observer *Observer) Connect(id int32, port int32) {
	observer.connector.Connect(id, port)
}

// Connect to master.
func (observer *Observer) ConnectMaster() {
	observer.Connect(MasterId, MasterPort)
}

// Handle the Kill signal from the master.
func (observer *Observer) inputKillHandle(connId int32, msg MessageBuffer) {
	LogI(fmt.Sprintf("Node %d Received kill signal", observer.id))
	observer.connector.SendAckRsp(MasterId, Input_KillRsp)
	os.Exit(0)
}

// Handle the CollectState signal from the master. Observer will send the collect state
// cmd to all nodes and get their snapshots, then save it to the observer snapshot map.
func (observer *Observer) inputCollectState_whandle(connId int32, msg MessageBuffer) {
	LogI(fmt.Sprintf("Node %d Received inputCollectState signal", observer.id))

	for peerId := range observer.connector.ConnectedConns {
		if peerId == MasterId || peerId == ObserverId {
			continue
		}
		snapShot := observer.collectStateHandle(peerId)
		observer.snapShots[peerId] = snapShot
	}

	observer.connector.SendAckRsp(MasterId, Input_CollectStateRsp)
}

// Send the collect state cmd to a node specified by an id.
func (observer *Observer) collectStateHandle(connId int32) SnapShot {
	msg := BinaryProtocol{}
	msg.Init(CollectState)
	observer.connector.WriteTo(connId, &msg)

	LogI(fmt.Sprintf("Node %d sent collectState to node %d", observer.id, connId))

	rspMsg := observer.connector.WaitRsp(connId)
	cmd := ConnectorCmd(rspMsg.ReadI32())
	if cmd != CollectStateRsp {
		LogE("Wrong collectState response")
	}
	LogI("Correct collectState response")
	money := rspMsg.ReadI64()
	channelsLen := rspMsg.ReadI32()
	channels := make(map[int32]int64)
	for i := int32(0); i < channelsLen; i++ {
		channelId := rspMsg.ReadI32()
		channelMoney := rspMsg.ReadI64()
		channels[channelId] = channelMoney
	}

	LogI(fmt.Sprintf("Observer finished collectState_wcall from node %d", connId))

	return SnapShot{NodeMoney: money, ChannelMoneys: channels}
}

// Handle the PrintSnapshot from the master. Observer will format the snapshots which were
// saved in the observer snapshot map, then print it.
func (observer *Observer) signalPrintSnapshotHandle(connId int32, msg MessageBuffer) {
	LogI(fmt.Sprintf("Node %d Received inputPrintSnapshot signal", observer.id))

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
	LogR(ProjectOutput(ret.String()))

	observer.connector.SendAckRsp(MasterId, Input_PrintSnapshotRsp)
}

// Create NodeState as the specification required
func CreateNodeState(nodeId int32, money int64) string {
	return fmt.Sprintf("node %d = %d\n", nodeId, money)
}

// Create ChannelState as the specification required
func CreateChannelState(sender int32, receiver int32, money int64) string {
	return fmt.Sprintf("channel (%d → %d) = %d\n", sender, receiver, money)
}
