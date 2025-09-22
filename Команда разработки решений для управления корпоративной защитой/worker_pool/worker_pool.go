package worker_pool

import (
	"context"
	"fmt"
	w "kaspersky/worker"
	"sync"
)

const queueSize int = 10

type WorkerPool struct {
	ctx     context.Context
	workers []*w.Worker
	pool    chan *w.Worker
	queue   chan func()
	wg      *sync.WaitGroup
}

func NewWorkerPool(ctx context.Context, numberOfWorkers int) *WorkerPool {
	wp := &WorkerPool{
		ctx:     ctx,
		workers: w.NewWorkers(numberOfWorkers),
		pool:    make(chan *w.Worker, numberOfWorkers),
		queue:   make(chan func(), queueSize),
		wg:      &sync.WaitGroup{},
	}

	wp.fillPool()
	wp.listenQueue()

	return wp
}

// Submit - добавить таску в воркер пул
func (wp *WorkerPool) Submit(task func()) {
	select {
	case w := <-wp.pool:
		fmt.Println("Got free worker, doing new task.")
		wp.doTaskAsync(w, task)
	case wp.queue <- task:
		fmt.Println("No free workers, sending task in queue.")
	}
}

// SubmitWait - добавить таску в воркер пул и дождаться окончания ее выполнения
func (wp *WorkerPool) SubmitWait(task func()) {
	w := <-wp.pool

	wp.doTaskSync(w, task)
}

// Stop - остановить воркер пул, дождаться выполнения только тех тасок, которые выполняются сейчас
func (wp *WorkerPool) Stop() {
	fmt.Println("Stopping worker pool...")

	close(wp.queue)

	wp.wg.Wait()

	for range wp.workers {
		<-wp.pool
	}

	fmt.Println("Worker pool was stopped.")
}

// StopWait - остановить воркер пул, дождаться выполнения всех тасок, даже тех, что не начали выполняться, но лежат в очереди
func (wp *WorkerPool) StopWait() {
	fmt.Println("Stopping worker pool...")

	close(wp.queue)

	wp.wg.Wait()

	for task := range wp.queue {
		w := <-wp.pool

		wp.doTaskAsync(w, task)
	}

	for range wp.workers {
		<-wp.pool
	}

	fmt.Println("Worker pool was stopped.")
}

func (wp *WorkerPool) fillPool() {
	for _, w := range wp.workers {
		wp.pool <- w
	}
}

func (wp *WorkerPool) doTaskSync(worker *w.Worker, task func()) {
	wp.wg.Add(1)
	defer wp.wg.Done()

	task()

	wp.pool <- worker
}

func (wp *WorkerPool) doTaskAsync(worker *w.Worker, task func()) {
	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()

		task()

		wp.pool <- worker
	}()
}

func (wp *WorkerPool) listenQueue() {
	go func() {
		for {
			task, ok := <-wp.queue
			if !ok {
				fmt.Println("Queue is closed!")
				return
			}
			fmt.Println("Got new task from queue.")

			w := <-wp.pool

			wp.doTaskAsync(w, task)
		}
	}()
}
