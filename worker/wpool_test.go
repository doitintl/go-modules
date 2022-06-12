package worker

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	ctx := context.Background()
	countOfTasks := 20
	tasks := make([]*testTask, countOfTasks)
	tasksCH := make(chan Task)
	wp := NewWorkerPool(ctx, tasksCH, 2)
	wp.AddWorkers(8)
	if wp.WorkersCount() != 10 {
		t.Errorf("Expected 10 workers, got %d", wp.WorkersCount())
	}

	wg := &sync.WaitGroup{}
	start := time.Now()
	for i := 0; i < countOfTasks; i++ {
		wg.Add(1)
		task := &testTask{Id: i, wg: wg}
		tasks[i] = task
		wp.AddTask(task)
	}
	wg.Wait()

	if time.Since(start) > time.Second*3 {
		t.Errorf("Expected to finish in less than 3 seconds, took %s", time.Since(start))
	}

	for i := 0; i < countOfTasks; i++ {
		if tasks[i].Id != i+1 {
			t.Errorf("Expected task id to be %d, got %d", i+1, tasks[i].Id)
		}
	}

	wp.RemoveWorkers(8)
	if wp.WorkersCount() != 2 {
		t.Errorf("Expected 2 workers, got %d", wp.WorkersCount())
	}

	wp.Stop()

	if wp.WorkersCount() != 0 {
		t.Errorf("Expected 0 workers, got %d", wp.WorkersCount())
	}
}
