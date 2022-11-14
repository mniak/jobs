package jobs

import (
	"context"

	"go.uber.org/atomic"
)

type loopjob struct {
	action  func(ctx context.Context)
	stop    chan struct{}
	stopped chan error
	started atomic.Bool
}

func NewLoopJob(action func(ctx context.Context)) Job {
	return &loopjob{
		action:  wrapActionWithPanicHandler(action),
		stop:    make(chan struct{}),
		stopped: make(chan error),
	}
}

func StartLoop(ctx context.Context, action func(ctx context.Context)) (StartedJob, error) {
	job := NewLoopJob(action)
	err := job.Start(ctx)
	return job, err
}

func (l *loopjob) Start(ctx context.Context) error {
	if !l.started.CompareAndSwap(false, true) {
		return ErrAlreadyStarted
	}
	go func() {
		for {
			select {
			case <-l.stop:
				l.stopped <- nil
				return
			case <-ctx.Done():
				l.stopped <- ctx.Err()
				return
			default:
				l.action(ctx)
			}
		}
	}()
	return nil
}

func (l *loopjob) Shutdown(ctx context.Context) error {
	if !l.started.Load() {
		return ErrNotYetStarted
	}
	close(l.stop)
	return nil
}

func (l *loopjob) Wait() error {
	return <-l.stopped
}

func wrapActionWithPanicHandler(action func(ctx context.Context)) func(ctx context.Context) {
	return func(ctx context.Context) {
		defer func() {
			recover()
		}()
		action(ctx)
	}
}
