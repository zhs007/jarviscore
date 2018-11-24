package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/coredb/proto"
	pb "github.com/zhs007/jarviscore/proto"
)

var (
	// EventOnNodeConnected - onNodeConnected
	EventOnNodeConnected = "nodeconnected"
	// EventOnIConnectNode - onIConnectNode
	EventOnIConnectNode = "connectnode"

	// EventOnCtrl - EventOnCtrl
	EventOnCtrl = "onctrl"
	// EventOnCtrlResult - EventOnCtrlResult
	EventOnCtrlResult = "ctrlresult"
	// EventOnReplyRequestFile - EventOnCtrlResult
	EventOnReplyRequestFile = "replyrequestfile"
)

// FuncNodeEvent - func event
type FuncNodeEvent func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error

// FuncMsgEvent - func event
type FuncMsgEvent func(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error

// eventMgr event manager
type eventMgr struct {
	mapNodeEvent map[string]([]FuncNodeEvent)
	mapMsgEvent  map[string]([]FuncMsgEvent)
	node         JarvisNode
}

func newEventMgr(node JarvisNode) *eventMgr {
	mgr := &eventMgr{
		mapNodeEvent: make(map[string]([]FuncNodeEvent)),
		mapMsgEvent:  make(map[string]([]FuncMsgEvent)),
		node:         node,
	}

	return mgr
}

func (mgr *eventMgr) checkNodeEvent(event string) bool {
	return event == EventOnNodeConnected ||
		event == EventOnIConnectNode
}

func (mgr *eventMgr) checkMsgEvent(event string) bool {
	return event == EventOnCtrl ||
		event == EventOnCtrlResult ||
		event == EventOnReplyRequestFile
}

func (mgr *eventMgr) regNodeEventFunc(event string, eventfunc FuncNodeEvent) error {
	if !mgr.checkNodeEvent(event) {
		return ErrInvalidEvent
	}

	mgr.mapNodeEvent[event] = append(mgr.mapNodeEvent[event], eventfunc)

	return nil
}

func (mgr *eventMgr) regMsgEventFunc(event string, eventfunc FuncMsgEvent) error {
	if !mgr.checkMsgEvent(event) {
		return ErrInvalidEvent
	}

	mgr.mapMsgEvent[event] = append(mgr.mapMsgEvent[event], eventfunc)

	return nil
}

func (mgr *eventMgr) onNodeEvent(ctx context.Context, event string, node *coredbpb.NodeInfo) error {
	lst, ok := mgr.mapNodeEvent[event]
	if !ok {
		return nil
	}

	for i := range lst {
		lst[i](ctx, mgr.node, node)
	}

	return nil
}

func (mgr *eventMgr) onMsgEvent(ctx context.Context, event string, msg *pb.JarvisMsg) error {
	lst, ok := mgr.mapMsgEvent[event]
	if !ok {
		return nil
	}

	for i := range lst {
		lst[i](ctx, mgr.node, msg)
	}

	return nil
}
