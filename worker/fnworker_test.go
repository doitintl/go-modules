package worker

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestFnWorker(t *testing.T) {
	ctx := context.Background()
	tasks := make(chan TaskFunc)
	w := NewFnWorker(tasks)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ok := false
	task := func(_ context.Context) {
		time.Sleep(time.Second * 1)
		ok = true
		wg.Done()
	}
	w.Start(ctx)
	tasks <- task
	w.Stop()
	if !ok {
		t.Errorf("Expected to task to be executed")
	}
}
