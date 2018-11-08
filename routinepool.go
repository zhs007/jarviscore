package jarviscore

import (
	"context"

	"github.com/zhs007/jarviscore/base"
)

// Task - task
type Task interface {
	Run(ctx context.Context) error
}

// routine - routine
type routine struct {
	routineID   int
	chanTask    chan Task
	chanRemove  chan *routine
	chanWaiting chan *routine
}

// sendTask - start a routine
func (r *routine) sendTask(task Task) {
	r.chanTask <- task
}

// start - start a routine
func (r *routine) start(ctx context.Context) error {
	jarvisbase.Debug("routine.Start...")

	for {
		isend := false
		select {
		case task, ok := <-r.chanTask:
			if !ok {
				isend = true
				break
			}

			jarvisbase.Debug("routine.Start:get new task")

			if task != nil {
				task.Run(ctx)
			}

			r.chanWaiting <- r
		case <-ctx.Done():
			jarvisbase.Debug("routine.Start:context done")
			isend = true
			break
		}

		if isend {
			break
		}
	}

	r.chanRemove <- r

	return nil
}

// RoutinePool - RoutinePool
type RoutinePool interface {
	// SendTask - send new task
	SendTask(task Task)
	// Start - start a routine pool
	Start(ctx context.Context, maxNums int) error
}

// routinePool - routinePool
type routinePool struct {
	lstRoutine  []*routine
	chanRemove  chan *routine
	chanWaiting chan *routine
	chanTask    chan Task
	lstTask     []Task
	lstWaiting  []*routine
	maxNums     int
}

// NewRoutinePool - new RoutinePool
func NewRoutinePool() RoutinePool {
	return &routinePool{
		chanRemove:  make(chan *routine, 128),
		chanWaiting: make(chan *routine, 128),
		chanTask:    make(chan Task, 256),
	}
}

// SendTask - send new task
func (pool *routinePool) SendTask(task Task) {
	pool.chanTask <- task
}

// Start - start a routine pool
func (pool *routinePool) Start(ctx context.Context, maxNums int) error {
	jarvisbase.Debug("RoutinePool.Start...")

	pool.maxNums = maxNums

	for {
		isend := false
		select {
		case task, ok := <-pool.chanTask:
			if !ok {
				isend = true
				break
			}

			jarvisbase.Debug("RoutinePool.Start:new task")

			pool.run(ctx, task)
		case r, ok := <-pool.chanWaiting:
			if !ok {
				isend = true
				break
			}

			jarvisbase.Debug("RoutinePool.Start:new waiting")

			pool.lstWaiting = append(pool.lstWaiting, r)

			pool.onNewWaiting()
		case r, ok := <-pool.chanRemove:
			if !ok {
				isend = true
				break
			}

			jarvisbase.Debug("RoutinePool.Start:new remove")

			pool.removeRoutine(r)
		case <-ctx.Done():
			jarvisbase.Debug("RoutinePool.Start:context done")

			isend = true
			break
		}

		if isend {
			break
		}
	}

	return nil
}

// onNewWaiting - on new waiting
func (pool *routinePool) onNewWaiting() error {
	if len(pool.lstTask) <= 0 {
		return nil
	}

	task := pool.lstTask[0]
	pool.lstTask = append(pool.lstTask[:0], pool.lstTask[1:]...)

	if len(pool.lstWaiting) > 0 {
		r := pool.lstWaiting[0]

		pool.lstWaiting = append(pool.lstWaiting[:0], pool.lstWaiting[1:]...)

		r.sendTask(task)

		return nil
	}

	return nil
}

// run - run with task
func (pool *routinePool) run(ctx context.Context, task Task) error {
	if len(pool.lstWaiting) > 0 {
		r := pool.lstWaiting[0]

		pool.lstWaiting = append(pool.lstWaiting[:0], pool.lstWaiting[1:]...)

		r.sendTask(task)

		return nil
	}

	if len(pool.lstRoutine) < pool.maxNums {
		r := &routine{
			routineID:   len(pool.lstRoutine) + 1,
			chanTask:    make(chan Task, 8),
			chanRemove:  pool.chanRemove,
			chanWaiting: pool.chanWaiting,
		}

		pool.lstRoutine = append(pool.lstRoutine, r)

		go pool.startRountine(ctx, r)

		r.sendTask(task)

		return nil
	}

	pool.lstTask = append(pool.lstTask, task)

	return nil
}

// startRountine - start a new routine
func (pool *routinePool) startRountine(ctx context.Context, r *routine) error {
	curctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return r.start(curctx)
}

// removeRoutine - remove a routine
func (pool *routinePool) removeRoutine(r *routine) {
	i := pool.findRoutine(r)
	if i >= 0 {
		pool.lstRoutine = append(pool.lstRoutine[:i], pool.lstRoutine[i+1:]...)
	}
}

// findRoutine - find a routine
func (pool *routinePool) findRoutine(r *routine) int {
	if len(pool.lstRoutine) <= 0 {
		return -1
	}

	for i := range pool.lstRoutine {
		if pool.lstRoutine[i].routineID == r.routineID {
			return i
		}
	}

	return -1
}
