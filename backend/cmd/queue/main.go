package main

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/queue"
	"bitbucket.org/projectiu7/backend/src/master/internal/queue/jobs/default_job"
	"bitbucket.org/projectiu7/backend/src/master/internal/queue/jobs/subdivision_org"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func main() {
	//dotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connStr := "postgres://mdb:mdb@localhost:5432/mdb"

	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	serviceQueue := queue.NewService(queue.NewRepository(dbpool), []default_job.Job{
		subdivision_org.NewJob(subdivision_org.NewRepository(dbpool)),
	})

	err = serviceQueue.Consume(ctx)
	if err != nil {
		return
	}
}
