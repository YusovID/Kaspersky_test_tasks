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

// в будущем можно перенести в флаг или конфиг
const numberOfWorkers = 10

func main() {
	// в самом конце выводим автора кода (меня)
	defer fmt.Println("This code was written by Ilya Yusov")

	// контекст с отменой для генерации и работы программы в целом
	ctx, cancel := context.WithCancel(context.Background())
	// на случай непредвиденной остановки программы
	defer cancel()

	// WaitGroup для главного процесса,
	// чтобы корректно дождаться завершения работы программы
	wg := &sync.WaitGroup{}

	runMainProcess(ctx, wg)

	// канал для отслуживания завершения работы программы, например, по Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// блокируемся, пока не получим сигнал
	<-sigChan

	cancel()

	fmt.Println("\nContext was cancelled!")

	// используем WaitGroup по назначению - дожидаемся окончания выполнения runMainProcess
	wg.Wait()

	fmt.Println("Program was stopped.")
}

// Функция runMainProcess управляет жизненным циклом пула воркеров.
// Работает в горутине, поэтому нужна WaitGroup
func runMainProcess(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		// канал для получения задач
		tasksChan := generator.GenerateTasks(ctx)

		workerPool := worker_pool.NewWorkerPool(ctx, numberOfWorkers)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping program...")
				workerPool.StopWait() // NOTE здесь можно поменять логику завершения (Stop или StopWait)
				return

			case task, ok := <-tasksChan:
				if !ok {
					fmt.Println("Tasks channel was closed, stopping program...")
					workerPool.StopWait() // NOTE здесь можно поменять логику завершения (Stop или StopWait)
					return
				}
				workerPool.Submit(task) // NOTE здесь тоже можно поменять логику обработки задачи (Submit или SubmitWait)
			}
		}
	}()
}
