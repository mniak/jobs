package jobs

import (
	"context"

	"go.uber.org/atomic"
)

type loopjob struct {
	action  func(ctx context.Context)
	stop    chan struct{}
	done    chan struct{}
	started atomic.Bool
}

func NewLoopJob(action func(ctx context.Context)) Job {
	return &loopjob{
		action: wrapActionWithPanicHandler(action),
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
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
		defer close(l.done)
		for {
			select {
			case <-l.stop:
				return
			case <-ctx.Done():
				return
			default:
				l.action(ctx)
			}
		}
	}()
	return nil
}

func (l *loopjob) Stop(ctx context.Context) error {
	if !l.started.Load() {
		return ErrNotYetStarted
	}
	close(l.stop)
	return nil
}

func (l *loopjob) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.done:
		return nil
	}
}

func wrapActionWithPanicHandler(action func(ctx context.Context)) func(ctx context.Context) {
	return func(ctx context.Context) {
		defer func() {
			recover()
		}()
		action(ctx)
	}
}
