package jarviscore

import (
	"context"
	"time"
)

// FuncOnTimer - on timer
// 			- If it returns false, the timer will end
type FuncOnTimer func(ctx context.Context, timer *Timer) bool

// Timer - timer
type Timer struct {
	timer   int
	ontimer FuncOnTimer
}

// NewTimer - new timer
func NewTimer(timer int, ontimer FuncOnTimer) *Timer {
	return &Timer{
		timer:   timer,
		ontimer: ontimer,
	}
}

// StartTimer - start timer
func StartTimer(ctx context.Context, timer int, ontimer FuncOnTimer) (*Timer, error) {
	if timer <= 0 {
		return nil, ErrInvalidTimer
	}

	if ontimer == nil {
		return nil, ErrInvalidTimerFunc
	}

	t := NewTimer(timer, ontimer)

	go t.Start(ctx)

	return t, nil
}

// Start - start
func (t *Timer) Start(ctx context.Context) {
	ct := time.Second * time.Duration(t.timer)
	timer := time.NewTimer(ct)

	for {
		timer.Reset(ct)

		select {
		case <-timer.C:
			if !t.ontimer(ctx, t) {
				break
			}
		case <-ctx.Done():
			break
		}
	}
}
