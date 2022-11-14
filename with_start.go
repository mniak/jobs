package jobs

import "context"

type jobWithPreStart struct {
	job       Job
	startFunc func(ctx context.Context) error
}

func WithPreStart(job Job, startFunc func(ctx context.Context) error) Job {
	return &jobWithPreStart{
		job:       job,
		startFunc: startFunc,
	}
}

func (j *jobWithPreStart) Start(ctx context.Context) error {
	err := j.startFunc(ctx)
	if err != nil {
		return err
	}
	return j.job.Start(ctx)
}

func (j *jobWithPreStart) Shutdown(ctx context.Context) error {
	return j.job.Shutdown(ctx)
}

func (j *jobWithPreStart) Wait() error {
	return j.job.Wait()
}
