package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/coredb/proto"
)

var (
	// EventOnNodeConnected - onNodeConnected
	EventOnNodeConnected = "nodeconnected"
	// EventOnIConnectNode - onIConnectNode
	EventOnIConnectNode = "connectnode"
)

// FuncNodeEvent - func event
type FuncNodeEvent func(ctx context.Context, jarvisnode JarvisNode, node *coredbpb.NodeInfo) error

// eventMgr event manager
type eventMgr struct {
	mapNodeEvent map[string]([]FuncNodeEvent)
	node         JarvisNode
}

func newEventMgr(node JarvisNode) *eventMgr {
	mgr := &eventMgr{
		mapNodeEvent: make(map[string]([]FuncNodeEvent)),
		node:         node,
	}

	return mgr
}

func (mgr *eventMgr) checkEvent(event string) bool {
	return event == EventOnNodeConnected ||
		event == EventOnIConnectNode
}

func (mgr *eventMgr) regEventFunc(event string, eventfunc FuncNodeEvent) error {
	if !mgr.checkEvent(event) {
		return ErrInvalidEvent
	}

	mgr.mapNodeEvent[event] = append(mgr.mapNodeEvent[event], eventfunc)

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
