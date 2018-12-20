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
	// EventOnIConnectNodeFail - onIConnectNodeFail
	EventOnIConnectNodeFail = "connectnodefail"

	// EventOnCtrl - OnCtrl
	EventOnCtrl = "onctrl"
	// EventOnCtrlResult - OnCtrlResult
	EventOnCtrlResult = "ctrlresult"
	// EventOnReplyRequestFile - OnReplyRequestFile
	EventOnReplyRequestFile = "replyrequestfile"

	// EventOnPrivateKey - OnPrivateKey
	EventOnPrivateKey = "privatekey"
)

// FuncNodeEvent - func event
type FuncNodeEvent func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error

// FuncMsgEvent - func event
type FuncMsgEvent func(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error

// FuncStateEvent - func event
type FuncStateEvent func(ctx context.Context, jarvisnode JarvisNode) error

// eventMgr event manager
type eventMgr struct {
	mapNodeEvent  map[string]([]FuncNodeEvent)
	mapMsgEvent   map[string]([]FuncMsgEvent)
	mapStateEvent map[string]([]FuncStateEvent)
	node          JarvisNode
}

func newEventMgr(node JarvisNode) *eventMgr {
	mgr := &eventMgr{
		mapNodeEvent:  make(map[string]([]FuncNodeEvent)),
		mapMsgEvent:   make(map[string]([]FuncMsgEvent)),
		mapStateEvent: make(map[string]([]FuncStateEvent)),
		node:          node,
	}

	return mgr
}

func (mgr *eventMgr) checkNodeEvent(event string) bool {
	return event == EventOnNodeConnected ||
		event == EventOnIConnectNode ||
		event == EventOnIConnectNodeFail
}

func (mgr *eventMgr) checkMsgEvent(event string) bool {
	return event == EventOnCtrl ||
		event == EventOnCtrlResult ||
		event == EventOnReplyRequestFile
}

func (mgr *eventMgr) checkStateEvent(event string) bool {
	return event == EventOnPrivateKey
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

func (mgr *eventMgr) regStateEventFunc(event string, eventfunc FuncStateEvent) error {
	if !mgr.checkStateEvent(event) {
		return ErrInvalidEvent
	}

	mgr.mapStateEvent[event] = append(mgr.mapStateEvent[event], eventfunc)

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

func (mgr *eventMgr) onStateEvent(ctx context.Context, event string) error {
	lst, ok := mgr.mapStateEvent[event]
	if !ok {
		return nil
	}

	for i := range lst {
		lst[i](ctx, mgr.node)
	}

	return nil
}
