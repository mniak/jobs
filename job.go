package jobs

import (
	"context"
	"errors"
)

type StartedJob interface {
	Stop(ctx context.Context) error
	Wait(ctx context.Context) error
}

//go:generate mockgen --package=jobs --destination=job_mock_test.go --source=job.go Job
type Job interface {
	Start(ctx context.Context) error
	StartedJob
}

var (
	ErrAlreadyStarted = errors.New("the job has already been started previously")
	ErrNotYetStarted  = errors.New("the job has never started yet")
)
