package queue

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"fmt"
	"time"
)

const (
	StatusAddedInQueue      = 0
	StatusRetryLimitReached = 1
	StatusSuccess           = 2
	StatusConfirm           = 3
)

type Repository struct {
	db utils.PgxPoolIface
}

func NewRepository(db utils.PgxPoolIface) *Repository {
	return &Repository{db: db}
}

func (r *Repository) EnqueueJob(ctx context.Context, uniID string, queue string, action string, entityType string, payload []byte) (int, error) {
	query := `insert into mdb.jobs (uni_id, queue, action, entity_type, payload) values ($1, $2, $3, $4, $5) returning id;`

	row := r.db.QueryRow(ctx, query, uniID, queue, action, entityType, string(payload))

	var IDJwt int

	err := row.Scan(&IDJwt)
	if err != nil {
		return 0, fmt.Errorf("scan error: %w", err)
	}

	return IDJwt, nil
}

func (r *Repository) DequeueJobs(ctx context.Context, lockedUntil time.Time) ([]Job, error) {
	query := `
		update
			jobs
		set attempts       = attempts + 1,
			reserved_until = $2
		where id IN (
			select id
			from jobs
			where status_id = $1 and queue = 'service' and reserved_until <= now()
			order by create_txtime
			for update skip locked
		)
		returning id, action, entity_type, payload, attempts;`

	rows, err := r.db.Query(ctx, query, StatusAddedInQueue, lockedUntil)
	if err != nil {
		return nil, fmt.Errorf("query dequeue ticket error: %w", err)
	}

	var jobs []Job

	for rows.Next() {
		var job Job

		err = rows.Scan(&job.ID, &job.Action, &job.EntityType, &job.Payload, &job.Attempts)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("got rows error: %w", err)
	}

	return jobs, nil
}

func (r *Repository) UpdateJobStatus(ctx context.Context, job Job, isOk bool) error {
	queryUpdateStatus := `
		update
			mdb.jobs
		set status_id = $1,
			result = $2,
			error = $3
		where id = $4`

	status := StatusRetryLimitReached
	if isOk {
		status = StatusSuccess
	}

	_, err := r.db.Exec(ctx, queryUpdateStatus, status, job.Result, job.Error, job.ID)
	if err != nil {
		return fmt.Errorf("query update success jobs status error: %w", err)
	}

	return nil
}

func (r *Repository) GetMessageCount(ctx context.Context, queue string, uniID string) ([]int, error) {
	query := `
			select id
			from mdb.jobs
			where status_id = $1 and queue = $2 and uni_id = $3;`

	rows, err := r.db.Query(ctx, query, StatusSuccess, queue, uniID)
	if err != nil {
		return nil, fmt.Errorf("query get msg count error: %w", err)
	}

	var IDs []int

	for rows.Next() {
		var ID int

		err = rows.Scan(&ID)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		IDs = append(IDs, ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("got rows error: %w", err)
	}

	return IDs, nil
}

func (r *Repository) GetData(ctx context.Context, queue string, uniID string, IDJwt int) (Job, error) {
	query := `
			select id, result, error, action, entity_type
			from mdb.jobs
			where queue = $1 and uni_id = $2 and status_id = $3 and id = $4`

	row := r.db.QueryRow(ctx, query, queue, uniID, StatusSuccess, IDJwt)

	var job Job

	err := row.Scan(&job.ID, &job.Result, &job.Error, &job.Action, &job.EntityType)
	if err != nil {
		return Job{}, fmt.Errorf("scan error: %w", err)
	}

	return job, nil
}

func (r *Repository) Confirm(ctx context.Context, queue string, uniID string, IDJwt int) (int, error) {
	query := `
		update
			mdb.jobs
		set status_id = $1
		where queue = $2 and uni_id = $3 and id = $4 and status_id <> $1
		returning id`

	row := r.db.QueryRow(ctx, query, StatusConfirm, queue, uniID, IDJwt)

	var IDJwtNew int

	err := row.Scan(&IDJwtNew)
	if err != nil {
		return 0, fmt.Errorf("scan error: %w", err)
	}

	return IDJwtNew, nil
}
