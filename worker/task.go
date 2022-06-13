package worker

import "context"

// Task is a task to be executed by a worker.
//go:generate mockery --name Task
type Task interface {
	Run(ctx context.Context)
}

// TaskFunc is a function that can run in a fnworker.
type TaskFunc func(ctx context.Context)
