package worker

import (
	"context"
	"sync"
)

//go:generate mockery --name Worker
type Worker interface {
	Start(ctx context.Context)
	Stop()
	IsRunning() bool
}

// worker is a struct that contains input channel of tasks
type worker struct {
	tasks     chan Task
	quit      chan struct{}
	isRunning bool
	mut       sync.Mutex
}

// NewWorker creates a new worker with input channel of tasks
// and returns a pointer to the worker
// The worker is not started yet
func NewWorker(tasks chan Task) Worker {
	return &worker{
		tasks: tasks,
		quit:  make(chan struct{}),
	}
}

// Start starts the worker loop to consume and run tasks
func (w *worker) Start(ctx context.Context) {
	w.mut.Lock()
	defer w.mut.Unlock()

	if w.isRunning {
		return
	}
	w.isRunning = true

	go func() {
		for {
			select {
			case task, ok := <-w.tasks:
				if !ok {
					w.Stop()
				}

				// run the task
				task.Run(ctx)

			case <-ctx.Done():
				w.Stop()

			case <-w.quit:
				return
			}
		}
	}()
}

// Stop stops the worker loop
func (w *worker) Stop() {
	w.mut.Lock()
	defer w.mut.Unlock()
	if !w.isRunning {
		return
	}
	w.quit <- struct{}{}
	w.isRunning = false
}

// IsRunning returns true if the worker is running
func (w *worker) IsRunning() bool {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.isRunning
}

// language: go
