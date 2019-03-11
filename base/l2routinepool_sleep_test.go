package jarvisbase

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

type sleepl2taskMgr struct {
	sync.RWMutex

	finishTask int
	sendTask   int
	totalTask  int
	mapIndex   map[int]int

	cancel context.CancelFunc
}

func (mgr *sleepl2taskMgr) init(pid int) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.mapIndex[pid] = 0
}

func (mgr *sleepl2taskMgr) send() {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.sendTask++
}

func (mgr *sleepl2taskMgr) finish(pid int) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.mapIndex[pid]++
	mgr.finishTask++
}

func (mgr *sleepl2taskMgr) check(pid int, index int) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	if mgr.mapIndex[pid] == index {
		return true
	}

	Debug("sleepl2taskMgr.check",
		zap.Int("pid", pid),
		zap.Int("mgr.mapIndex", mgr.mapIndex[pid]),
		zap.Int("index", index))

	return false
}

func (mgr *sleepl2taskMgr) isok() bool {
	mgr.RLock()
	defer mgr.RUnlock()

	Debug("sleepl2taskMgr:isok",
		zap.Int("total", mgr.totalTask),
		zap.Int("send", mgr.sendTask),
		zap.Int("finish", mgr.finishTask))

	return mgr.finishTask == mgr.sendTask && mgr.finishTask == mgr.totalTask
}

type sleepl2taskTest struct {
	index int
	j     int
	mgr   *sleepl2taskMgr
	pid   string
}

func (task *sleepl2taskTest) Run(ctx context.Context) error {
	// time.Sleep(1 * time.Second)

	if !task.mgr.check(task.j, task.index) {
		task.mgr.cancel()

		return nil
	}

	time.Sleep(time.Microsecond * 10)
	Debug("run ", zap.Int("index", task.index), zap.Int("j", task.j))
	task.mgr.finish(task.j)

	if task.mgr.isok() {
		task.mgr.cancel()

		return nil
	}

	return nil
}

// GetParentID - get parentID
func (task *sleepl2taskTest) GetParentID() string {
	return task.pid
}

func sleepl2makeTask(mgr *sleepl2taskMgr, pool L2RoutinePool, maxpid int, maxc int) {
	time.Sleep(3 * time.Second)

	mgr.totalTask = maxpid * maxc

	for i := 0; i < maxc; i++ {
		for j := 0; j < maxpid; j++ {
			if i == 0 {
				mgr.init(j)
			}

			task := &sleepl2taskTest{
				index: i,
				j:     j,
				mgr:   mgr,
				pid:   fmt.Sprintf("pid%v", j),
			}

			pool.SendTask(task)
			mgr.send()
		}
	}

	// time.Sleep(3 * time.Second)
	// Debug("end")
	// fmt.Print("end\n")
	// cancel()
}

func TestL2RountinePoolSleep128(t *testing.T) {
	InitLogger(zap.InfoLevel, true, "")
	Debug("start...")
	// fmt.Print("haha\n")

	pool := NewL2RoutinePool()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mgr := &sleepl2taskMgr{
		cancel:   cancel,
		mapIndex: make(map[int]int),
	}
	go sleepl2makeTask(mgr, pool, 100, 10)

	pool.Start(ctx, 128)

	if !mgr.isok() {
		t.Fail()
	}
}

func TestL2RountinePoolSleep1(t *testing.T) {
	InitLogger(zap.InfoLevel, true, "")
	Debug("start...")
	// fmt.Print("haha\n")

	pool := NewL2RoutinePool()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mgr := &sleepl2taskMgr{
		cancel:   cancel,
		mapIndex: make(map[int]int),
	}
	go sleepl2makeTask(mgr, pool, 100, 10)

	pool.Start(ctx, 1)

	if !mgr.isok() {
		t.Fail()
	}
}

func TestL2RountinePoolSleep2(t *testing.T) {
	InitLogger(zap.InfoLevel, true, "")
	Debug("start...")
	// fmt.Print("haha\n")

	pool := NewL2RoutinePool()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mgr := &sleepl2taskMgr{
		cancel:   cancel,
		mapIndex: make(map[int]int),
	}
	go sleepl2makeTask(mgr, pool, 100, 10)

	pool.Start(ctx, 2)

	if !mgr.isok() {
		t.Fail()
	}
}
