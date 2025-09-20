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
	WorkTime   = 3
)

type Pool interface {
	// Submit - добавить задачу в пул.
	// Если пул не имеет свободных воркеров, то задачу нужно добавить в очередь.
	Submit(task func()) error

	// Stop - остановить воркер пул, дождаться выполнения всех добавленных ранее в очередь задач.
	Stop() error
}

// Дополнительно:
// - Ограничить размер очереди. Если очередь переполнена то вернуть ошибку из метода Submit.
// - Реализовать поддержку хука, который будет вызываться после выполнения каждой задачи.

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), WorkTime*time.Second)
	defer cancel()

	out := generator.Generate(ctx)

	var pool Pool = pool.New(NumWorkers, QueueSize, hook)
	defer pool.Stop()

	for {
		task, isOpen := <-out
		if !isOpen {
			fmt.Println("channel was closed")
			return
		}

		err := pool.Submit(task)
		if err != nil {
			fmt.Printf("failed to submit task: %v\n", err)

			// TODO retry + DLQ

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func hook(workerId int) {
	fmt.Printf("worker %d done task\n", workerId)
}
