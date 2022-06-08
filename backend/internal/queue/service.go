package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/projectiu7/backend/src/master/internal/queue/jobs/default_job"
)

const DefaultLockTime = 5 * time.Second

//go:generate mockery --name repositoryTicketEvent --structname repositoryTicketEventMock --output . --filename repository_ticket_event_mock_test.go --inpackage
type repository interface {
	EnqueueJob(ctx context.Context, uniID string, queue string, action string, entityType string, payload []byte) (int, error)
	DequeueJobs(ctx context.Context, lockedUntil time.Time) ([]Job, error)
	UpdateJobStatus(ctx context.Context, job Job, isOk bool) error

	GetMessageCount(ctx context.Context, queue string, uniID string) ([]int, error)
	GetData(ctx context.Context, queue string, uniID string, IDJwt int) (Job, error)
	Confirm(ctx context.Context, queue string, uniID string, IDJwt int) (int, error)
}

type Service struct {
	repository repository
	jobs       map[string]default_job.Job
}

func NewService(repository repository, newJobs []default_job.Job) *Service {
	jj := make(map[string]default_job.Job)

	for _, j := range newJobs {
		jj[j.GetEntityType()] = j
	}

	return &Service{repository: repository, jobs: jj}
}

func (s *Service) Consume(ctx context.Context) error {
	newJobs, err := s.repository.DequeueJobs(ctx, time.Now().Local().Add(DefaultLockTime))
	if err != nil {
		return fmt.Errorf("get jobs from queue error: %w", err)
	}

	for _, job := range newJobs {
		j, ok := s.jobs[job.EntityType]
		if !ok {
			job.Result = ""
			job.Error = "dont have job of this type"

			err = s.repository.UpdateJobStatus(ctx, job, false)
			if err != nil {
				return fmt.Errorf("delete jobs in queue error: %w", err)
			}
			continue
		}

		if job.Attempts >= j.MaxAttempts() {
			job.Result = ""
			job.Error = "max attempts reached"

			err = s.repository.UpdateJobStatus(ctx, job, false)
			if err != nil {
				return fmt.Errorf("delete jobs in queue error: %w", err)
			}
			continue
		}

		res, err := j.Process(ctx, job.Action, job.Payload)
		if err != nil {
			job.Result = ""
			job.Error = err.Error()

			err = s.repository.UpdateJobStatus(ctx, job, false)
			if err != nil {
				return fmt.Errorf("delete jobs in queue error: %w", err)
			}
			continue
		}

		job.Result = res
		job.Error = ""

		err = s.repository.UpdateJobStatus(ctx, job, true)
		if err != nil {
			return fmt.Errorf("delete jobs in queue error: %w", err)
		}
	}

	return nil
}

func (s *Service) Produce(ctx context.Context, uniID string, action string, entityType string, payload []byte) (int, error) {
	return s.repository.EnqueueJob(ctx, uniID, "service", action, entityType, payload)
}

func (s *Service) GetMessageCount(ctx context.Context, queue string, uniID string) (string, error) {
	IDs, err := s.repository.GetMessageCount(ctx, queue, uniID)

	var respCountStruct RespCount
	respCountStruct.Messages = len(IDs)
	respCountStruct.IDJwts = IDs

	resp, err := json.Marshal(respCountStruct)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	return string(resp), nil

}

func (s *Service) GetData(ctx context.Context, queue string, uniID string, IDJwt int) (Job, error) {
	return s.repository.GetData(ctx, queue, uniID, IDJwt)
}

func (s *Service) Confirm(ctx context.Context, queue string, uniID string, IDJwt int) (int, error) {
	return s.repository.Confirm(ctx, queue, uniID, IDJwt)
}

//// You can edit this code!
//// Click here and start typing.
//package main
//
//import (
//"encoding/base64"
//"fmt"
//)
//
//func main() {
//
//	fmt.Println(base64.StdEncoding.EncodeToString([]byte("{\"action\":\"add\", \"entityType\":\"SubdivisionOrg\", \"ogrn\":\"1111\", \"kpp\":\"2222\"}")) + "." + base64.StdEncoding.EncodeToString([]byte("<PackageData><SubdivisionOrg><UID>string</UID><Name>string</Name></SubdivisionOrg></PackageData>")))
//}
