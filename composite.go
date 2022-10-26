package jobs

import (
	"context"

	"go.uber.org/multierr"
)

type CompositeJob struct {
	Jobs []Job
}

func (cj *CompositeJob) Start(ctx context.Context) error {
	statuses := make(map[Job]error)
	for _, job := range cj.Jobs {
		statuses[job] = job.Start(ctx)
	}

	var result error
	for job, err := range statuses {
		if multierr.AppendInto(&result, err) {
			multierr.AppendInto(&result, job.Stop(ctx))
		}
	}
	return result
}

func (cj *CompositeJob) Wait(ctx context.Context) error {
	var err error
	for _, srv := range cj.Jobs {
		w := srv.Wait(ctx)
		multierr.AppendInto(&err, w)
	}
	return err
}

func (mcjs *CompositeJob) Stop(ctx context.Context) error {
	var err error
	for _, job := range mcjs.Jobs {
		multierr.AppendInto(&err, job.Stop(ctx))
	}
	return err
}
