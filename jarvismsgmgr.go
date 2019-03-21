package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/base"
	pb "github.com/zhs007/jarviscore/proto"
)

type jarvisMsgTask struct {
	// msg                *pb.JarvisMsg
	// msgstream          pb.JarvisCoreServ_ProcMsgServer
	// lstmsg             []*pb.JarvisMsg
	// lstmsgstreamstream pb.JarvisCoreServ_ProcMsgStreamServer
	mgr      *jarvisMsgMgr
	chanEnd  chan int
	taskinfo JarvisTask
	// funcOnResult       FuncOnProcMsgResult
}

func (task *jarvisMsgTask) Run(ctx context.Context) error {
	err := task.mgr.node.OnMsg(ctx, task.taskinfo)

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
		taskinfo: JarvisTask{
			Normal: &NormalTaskInfo{
				Msg:      msg,
				Stream:   stream,
				OnResult: funcOnResult,
			},
		},
		mgr:     mgr,
		chanEnd: chanEnd,
	}

	mgr.pool.SendTask(task)
}

// sendStreamMsg - send stream msg
func (mgr *jarvisMsgMgr) sendStreamMsg(msgs []*pb.JarvisMsg, stream pb.JarvisCoreServ_ProcMsgStreamServer,
	chanEnd chan int, funcOnResult FuncOnProcMsgResult) {

	task := &jarvisMsgTask{
		taskinfo: JarvisTask{
			Stream: &StreamTaskInfo{
				Msgs:     msgs,
				Stream:   stream,
				OnResult: funcOnResult,
			},
		},
		mgr:     mgr,
		chanEnd: chanEnd,
	}

	mgr.pool.SendTask(task)
}

// start - start goroutine to proc ctrl msg
func (mgr *jarvisMsgMgr) start(ctx context.Context) error {
	return mgr.pool.Start(ctx, 1)
}
