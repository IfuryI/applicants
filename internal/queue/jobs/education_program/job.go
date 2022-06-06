package education_program

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/queue/jobs/default_job"
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"encoding/xml"
	"fmt"
)

const (
	EntityType = "Campaign"
)

type RepositoryJob interface {
	Add(ctx context.Context, so PackageData) (string, error)
	Get(ctx context.Context, so PackageData) (PackageData, error)
}
type job struct {
	default_job.DefaultJob
	repositoryJob RepositoryJob
}

func NewJob(repositoryJob RepositoryJob) *job {
	return &job{
		repositoryJob: repositoryJob,
	}
}

func (r *job) GetEntityType() string {
	return EntityType
}

func (r *job) Process(ctx context.Context, action string, payload string) (string, error) {
	var pd PackageData

	err := xml.Unmarshal([]byte(payload), &pd)
	if err != nil {
		return "", err
	}

	switch action {
	case utils.ActionAdd:
		_, err := r.repositoryJob.Add(ctx, pd)
		if err != nil {
			return "", fmt.Errorf("add error: %w", err)
		}

		return "<PackageData></PackageData>", nil
	case utils.ActionGet:
		soReq, err := r.repositoryJob.Get(ctx, pd)
		if err != nil {
			return "", fmt.Errorf("get error: %w", err)
		}

		result, err := xml.MarshalIndent(soReq, " ", "  ")
		if err != nil {
			return "", fmt.Errorf("get error: %w", err)
		}

		return string(result), nil
	default:
		return "", fmt.Errorf("job with wrong action")
	}
}
