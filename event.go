package jarviscore

import (
	"github.com/zhs007/jarviscore/coredb/proto"
)

var (
	// EventOnNodeConnected - onNodeConnected
	EventOnNodeConnected = "nodeconnected"
)

// FuncNodeEvent - func event
type FuncNodeEvent func(node *coredbpb.NodeInfo) error

// eventMgr event manager
type eventMgr struct {
	mapNodeEvent map[string]([]FuncNodeEvent)
}

func newEventMgr() *eventMgr {
	mgr := &eventMgr{
		mapNodeEvent: make(map[string]([]FuncNodeEvent)),
	}

	return mgr
}

func (mgr *eventMgr) checkEvent(event string) bool {
	return event == EventOnNodeConnected
}

func (mgr *eventMgr) regEventFunc(event string, eventfunc FuncNodeEvent) error {
	if !mgr.checkEvent(event) {
		return ErrInvalidEvent
	}

	mgr.mapNodeEvent[event] = append(mgr.mapNodeEvent[event], eventfunc)

	return nil
}

func (mgr *eventMgr) onNodeEvent(event string, node *coredbpb.NodeInfo) error {
	lst, ok := mgr.mapNodeEvent[event]
	if !ok {
		return nil
	}

	for i := range lst {
		lst[i](node)
	}

	return nil
}
