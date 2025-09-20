package main

import (
	"context"
	"fmt"
	"sandbox-team/generator"
	pool "sandbox-team/worker_pool"
	"time"
)

const (
	NumWorkers = 10
	QueueSize  = 10
	WorkTime   = 3 // seconds
)

// main sets up and runs the worker pool demo.
//
// It creates a task generator, a worker pool, and submits tasks from
// the generator to the pool until a timeout is reached.
func main() {
	// Set an overall timeout for the application's execution.
	ctx, cancel := context.WithTimeout(context.Background(), WorkTime*time.Second)
	defer cancel()

	// Start the task generator. It will stop when the context is cancelled.
	tasksChan := generator.Generate(ctx)

	// Initialize the worker pool. defer pool.Stop() ensures that we wait for
	// all queued tasks to complete before exiting the program.
	workerPool := pool.New(NumWorkers, QueueSize, hook)
	defer workerPool.Stop()

	// Continuously read tasks from the generator and submit them to the pool.
	for task := range tasksChan {
		err := workerPool.Submit(task)
		if err != nil {
			fmt.Printf("failed to submit task: %v\n", err)

			// TODO: Implement a proper retry mechanism with backoff.
			// The current implementation loses the task on failure and introduces
			// a simple sleep, which is inefficient.
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("Task generator's channel was closed, main loop is finishing.")
}

// hook is a callback function executed after each task completion.
// It logs the ID of the worker that finished the task.
func hook(workerId int) {
	fmt.Printf("worker %d done task\n", workerId)
}
