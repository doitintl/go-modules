package worker

import (
	"context"
	"sync"
	"testing"
	"time"
)

type testTask struct {
	wg *sync.WaitGroup
	Id int
}

func (t *testTask) Run(ctx context.Context) {
	t.Id = t.Id + 1
	time.Sleep(time.Second * 1)
	t.wg.Done()
}

func TestWorker(t *testing.T) {
	ctx := context.Background()
	tasks := make(chan Task)
	w := NewWorker(tasks)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	task := &testTask{Id: 1, wg: wg}
	w.Start(ctx)
	tasks <- task
	w.Stop()
	if task.Id != 2 {
		t.Errorf("Expected task id to be 2, got %d", task.Id)
	}
}
