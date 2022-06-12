package worker

import (
	"context"
	"sync"
)

type fnworkerPool struct {
	mux         sync.Mutex
	ctx         context.Context
	taskChannel chan TaskFunc
	workers     []*fnworker
}

func NewFnWorkerPool(ctx context.Context, taskChannel chan TaskFunc, numWorkers int) *fnworkerPool {
	pool := &fnworkerPool{
		workers:     make([]*fnworker, numWorkers),
		ctx:         ctx,
		taskChannel: taskChannel,
	}

	for i := 0; i < numWorkers; i++ {
		w := NewFnWorker(taskChannel)
		pool.workers[i] = w
		w.Start(ctx)
	}

	return pool
}

// Stop stops all workers in the fnworkerPool.
func (p *fnworkerPool) Stop() {
	p.mux.Lock()
	defer p.mux.Unlock()
	for i := 0; i <= len(p.workers); i++ {
		if len(p.workers) > 0 {
			w := p.workers[len(p.workers)-1]
			w.Stop()
			p.workers = p.workers[:len(p.workers)-1]
		}
	}
}

// Add a single task to process queue
func (p *fnworkerPool) AddTask(task TaskFunc) {
	p.taskChannel <- task
}

// Reurns the queue chanf of the workerPool
func (p *fnworkerPool) TaskChannel() chan TaskFunc {
	return p.taskChannel
}

// Returns the number of workers in the workerPool
func (p *fnworkerPool) WorkersCount() int {
	return len(p.workers)
}

// AddWorkers adds the given number of workers to the workerPool.
func (p *fnworkerPool) AddWorkers(count int) {
	p.mux.Lock()
	defer p.mux.Unlock()
	for i := 0; i < count; i++ {
		w := NewFnWorker(p.taskChannel)
		p.workers = append(p.workers, w)
		w.Start(p.ctx)
	}
}

// RemoveWorkers removes the given number of workers from the workerPool.
// The workers are stopped and removed from the workerPool.
// If the number of workers to remove is greater than the number of workers in the workerPool,
// all workers are removed.
func (p *fnworkerPool) RemoveWorkers(count int) {
	p.mux.Lock()
	defer p.mux.Unlock()
	for i := 0; i < count; i++ {
		if len(p.workers) > 0 {
			w := p.workers[len(p.workers)-1]
			w.Stop()
			p.workers = p.workers[:len(p.workers)-1]
		}
	}
}

// Language: go
