package pool

import (
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrQueueFull is returned when a task is submitted to a pool with a full task queue.
	ErrQueueFull = errors.New("queue is full")
	// ErrPoolStopped is returned when a task is submitted to a pool that has been stopped.
	ErrPoolStopped = errors.New("pool is stopped")
)

// Pool manages a collection of worker goroutines to execute tasks concurrently.
type Pool struct {
	tasks     chan func()     // A buffered channel that acts as the task queue.
	hook      func(int)       // A function called by a worker after completing a task.
	isStopped bool            // A flag indicating the pool is in the process of stopping.
	mu        *sync.Mutex     // Protects access to the isStopped flag.
	wg        *sync.WaitGroup // Synchronizes the completion of all tasks upon stopping.
}

// New creates, initializes, and starts a new worker pool.
//
// numWorkers specifies the number of worker goroutines to run.
// queueSize specifies the maximum number of tasks that can be waiting in the queue.
// hook is an optional function that is called after each task is completed.
func New(numWorkers int, queueSize int, hook func(workerId int)) *Pool {
	pool := &Pool{
		tasks: make(chan func(), queueSize),
		hook:  hook,
		mu:    &sync.Mutex{},
		wg:    &sync.WaitGroup{},
	}

	pool.run(numWorkers)

	return pool
}

// run starts the specified number of worker goroutines.
func (p *Pool) run(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		// Pass the worker ID as a parameter to the goroutine to avoid the classic
		// loop variable closure issue. Each goroutine gets its own copy of workerID.
		go func(workerID int) {
			fmt.Printf("worker %d started\n", workerID+1)
			for task := range p.tasks {
				p.doTask(workerID, task)
			}
			fmt.Printf("worker %d stopped\n", workerID+1)
		}(i + 1) // Worker IDs are 1-based for clarity.
	}
}

// Submit adds a task to the pool's queue for execution.
//
// It is non-blocking. If the queue is full, it returns ErrQueueFull immediately.
// If the pool has been stopped via Stop(), it returns ErrPoolStopped.
func (p *Pool) Submit(task func()) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isStopped {
		return ErrPoolStopped
	}

	// This select is non-blocking. It tries to send the task, but if the channel
	// buffer is full, it immediately executes the default case.
	select {
	case p.tasks <- task:
		// We only increment the WaitGroup counter *after* successfully
		// adding the task to the queue.
		p.wg.Add(1)
		return nil
	default:
		return ErrQueueFull
	}
}

// Stop gracefully shuts down the worker pool.
//
// It stops accepting new tasks and waits for all previously submitted tasks
// to be completed by the workers.
func (p *Pool) Stop() {
	fmt.Println("Stopping pool...")

	p.mu.Lock()
	// Prevent new tasks from being submitted.
	p.isStopped = true
	p.mu.Unlock()

	// Closing the channel signals to the worker goroutines (in their range loop)
	// that no more tasks will be sent.
	close(p.tasks)

	// Wait for all tasks that were already in the queue to be processed.
	p.wg.Wait()

	fmt.Println("Worker pool was stopped.")
}

// doTask is an internal helper that executes a single task.
// It ensures that the WaitGroup counter is decremented and the post-execution
// hook is called after the task function completes.
func (p *Pool) doTask(workerId int, task func()) {
	// Decrement the WaitGroup counter once the task is done.
	defer p.wg.Done()

	task()

	if p.hook != nil {
		p.hook(workerId)
	}
}
