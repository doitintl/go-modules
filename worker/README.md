# Worker

```go get github.com/doitintl/go-modules/worker```

Worker pool module to run parallel tasks

You sould implement the ```Task``` interface.

```
   type Task interface {
      Run(ctx context.Context)
    }
```
OR implement a task func
```
    type TaskFunc func(ctx context.Context)
```


You can make simple or buffered Task channel.

Then You can delegate the task to worker pool with method or just send into the
channel.

```
    tasksCH := make(chan worker.Task)
    wp := worker.NewWorkerPool(ctx, tasksCH, countOfWorkers)
    
    wp.AddTask(task)
    // or
    taskCH <- task
```
OR
```
    tasksCH := make(chan TaskFunc)
    wp := worker.NewFnWorkerPool(ctx, tasksCH, countOfWorkers)
    wp.AddTask(func(ctx context.Context){
      ...
    })
```


You can add or remove workers from the pool:
```
    wp.AddWorkers(8)
    wp.RemoveWorkers(8)
```

You can get count of workers from the pool:
```
    wp.WorkersCount()
```

You should stop the pool for graceful shutdown.

```
    wp.Stop()
```
or with context cancellation




You can start again the pool.
```
    wp.Start(ctx)
```

