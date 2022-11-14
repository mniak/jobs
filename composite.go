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
	hasError := false
	for _, job := range cj.Jobs {
		err := job.Start(ctx)
		statuses[job] = err
		if err != nil {
			hasError = true
		}
	}

	var result error
	if hasError {
		for job, err := range statuses {
			if !multierr.AppendInto(&result, err) {
				multierr.AppendInto(&result, job.Shutdown(ctx))
			}
		}
	}
	return result
}

func (cj *CompositeJob) Wait() error {
	var err error
	for _, srv := range cj.Jobs {
		w := srv.Wait()
		multierr.AppendInto(&err, w)
	}
	return err
}

func (cj *CompositeJob) Shutdown(ctx context.Context) error {
	var err error
	for _, job := range cj.Jobs {
		multierr.AppendInto(&err, job.Shutdown(ctx))
	}
	return err
}
