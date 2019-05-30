package jarvisbase

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

type taskMgr struct {
	sync.RWMutex

	finishTask int
	sendTask   int
	totalTask  int
	mapTask    map[int]bool

	cancel context.CancelFunc
}

func (mgr *taskMgr) getOutputString() string {
	return fmt.Sprintf("taskMgr finishTask %v sendTask %v totalTask %v", mgr.finishTask, mgr.sendTask, mgr.totalTask)
}

func (mgr *taskMgr) send() {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.sendTask++
}

func (mgr *taskMgr) finish(index int) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.finishTask++

	_, ok := mgr.mapTask[index]
	if ok {
		mgr.mapTask[index] = false
	} else {
		mgr.mapTask[index] = true
	}
}

func (mgr *taskMgr) isok(index int) bool {
	mgr.Lock()
	defer mgr.Unlock()

	if index == -1 {
		if mgr.finishTask == mgr.sendTask && mgr.finishTask == mgr.totalTask {
			nums := 0
			for _, v := range mgr.mapTask {
				if !v {
					return false
				}

				nums++
			}

			return nums == mgr.finishTask
		}

		return false
	}

	if mgr.finishTask == mgr.sendTask && mgr.finishTask == mgr.totalTask {
		isf, isok := mgr.mapTask[index]
		if !isok {
			return false
		}

		if !isf {
			return false
		}

		return true
	}

	return false
}

type taskTest struct {
	index int
	mgr   *taskMgr
}

func (task *taskTest) Run(ctx context.Context) error {
	// time.Sleep(1 * time.Second)
	Debug("run ", zap.Int("index", task.index))
	task.mgr.finish(task.index)

	if task.mgr.isok(task.index) {
		task.mgr.cancel()

		return nil
	}

	return nil
}

func makeTask(mgr *taskMgr, pool RoutinePool, maxtasks int) {
	time.Sleep(3 * time.Second)

	mgr.totalTask = maxtasks

	for i := 0; i < maxtasks; i++ {
		task := &taskTest{
			index: i,
			mgr:   mgr,
		}

		pool.SendTask(task)
		mgr.send()
	}

	// time.Sleep(3 * time.Second)
	// Debug("end")
	// fmt.Print("end\n")
	// cancel()
}

func TestRountinePool(t *testing.T) {
	InitLogger(zap.InfoLevel, true, "", "")
	Debug("start...")
	// fmt.Print("haha\n")

	pool := NewRoutinePool()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mgr := &taskMgr{
		cancel:  cancel,
		mapTask: make(map[int]bool),
	}
	go makeTask(mgr, pool, 10000)

	pool.Start(ctx, 128)

	if !mgr.isok(-1) {
		Error("TestRountinePool",
			zap.String("output", mgr.getOutputString()),
			zap.String("pool", pool.GetStatus()))

		t.Fail()
	}
}
