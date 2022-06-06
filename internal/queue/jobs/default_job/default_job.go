package default_job

import (
	"context"
)

const (
	DefaultMaxAttempts = 5
)

type Job interface {
	GetEntityType() string
	Process(ctx context.Context, action string, payload string) (string, error)
	MaxAttempts() int64
}

type DefaultJob struct{}

func (j *DefaultJob) MaxAttempts() int64 {
	return DefaultMaxAttempts
}
