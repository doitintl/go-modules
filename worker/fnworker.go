package worker

import (
	"context"
	"sync"
)

//go:generate mockery --name FnWorker
type FnWorker interface {
	Start(ctx context.Context)
	Stop()
	IsRunning() bool
}

// worker is a struct that contains input channel of tasks fn
type fnworker struct {
	tasks     chan TaskFunc
	quit      chan struct{}
	isRunning bool
	mut       sync.Mutex
}

func NewFnWorker(tasks chan TaskFunc) FnWorker {
	return &fnworker{
		tasks: tasks,
		quit:  make(chan struct{}),
	}
}

// Start starts the worker loop to consume and run tasks
func (w *fnworker) Start(ctx context.Context) {
	w.mut.Lock()
	defer w.mut.Unlock()

	if w.isRunning {
		return
	}
	w.isRunning = true

	go func(ctx context.Context) {
		for {
			select {
			case task, ok := <-w.tasks:
				if !ok {
					w.Stop()
				}

				// run the task
				task(ctx)

			case <-ctx.Done():
				w.Stop()

			case <-w.quit:
				return
			}
		}
	}(ctx)
}

// Stop stops the worker loop
func (w *fnworker) Stop() {
	w.mut.Lock()
	defer w.mut.Unlock()
	if !w.isRunning {
		return
	}
	w.quit <- struct{}{}
	w.isRunning = false
}

// IsRunning returns true if the worker is running
func (w *fnworker) IsRunning() bool {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.isRunning
}

// language: go
