package main

import (
	"context"
	"fmt"
	"kaspersky/generator"
	"kaspersky/worker_pool"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const numberOfWorkers = 10

func main() {
	defer fmt.Println("This code was written by Ilya Yusov")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go runMainProcess(ctx, wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()

	fmt.Println("\nContext was cancelled!")

	wg.Wait()

	fmt.Println("Program was stopped.")
}

func runMainProcess(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	tasksChan := generator.GenerateTasks(ctx)

	workerPool := worker_pool.NewWorkerPool(ctx, numberOfWorkers)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping program...")
			workerPool.StopWait()
			return

		case task := <-tasksChan:
			workerPool.Submit(task)
		}
	}
}
