package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
)

type jarvisMsgTask struct {
	msg          *pb.JarvisMsg
	stream       pb.JarvisCoreServ_ProcMsgServer
	mgr          *jarvisMsgMgr
	chanEnd      chan int
	funcOnResult FuncOnProcMsgResult
}

func (task *jarvisMsgTask) Run(ctx context.Context) error {
	err := task.mgr.node.OnMsg(ctx, task.msg, task.stream, task.funcOnResult)

	if task.chanEnd != nil {
		task.chanEnd <- 0
	}

	return err
}

// jarvisMsgMgr - jarvis msg mgr
type jarvisMsgMgr struct {
	pool jarvisbase.RoutinePool
	node JarvisNode
	// mgrClient2 *jarvisClient2
}

// newJarvisMsgMgr - new jarvisMsgMgr
func newJarvisMsgMgr(node JarvisNode) *jarvisMsgMgr {
	mgr := &jarvisMsgMgr{
		pool: jarvisbase.NewRoutinePool(),
		node: node,
		// mgrClient2: newClient2(node),
	}

	return mgr
}

// sendMsg - send a ctrl msg
func (mgr *jarvisMsgMgr) sendMsg(msg *pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgServer,
	chanEnd chan int, funcOnResult FuncOnProcMsgResult) {

	task := &jarvisMsgTask{
		msg:          msg,
		stream:       stream,
		mgr:          mgr,
		chanEnd:      chanEnd,
		funcOnResult: funcOnResult,
	}

	mgr.pool.SendTask(task)
}

// start - start goroutine to proc ctrl msg
func (mgr *jarvisMsgMgr) start(ctx context.Context) error {
	return mgr.pool.Start(ctx, 1)
}
