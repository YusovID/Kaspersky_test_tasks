// Пакет worker_pool используется для работы с пулом воркеров,
// включает в себя методы:
//
// - NewWorkerPool - создание нового пула воркеров
//
// - Submit - асинхронное выполнение задачи
//
// - SubmitWait - синхронное выполнение задачи
//
// - Stop - остановить пул воркеров, не выполняя задач из очереди
//
// - StopWait - остановить пул воркеров, выполняя задачи из очереди
package worker_pool

import (
	"context"
	"fmt"
	w "kaspersky/worker"
	"sync"
)

// можно в будущем сделать через флаг или брать из конфига
const QueueSize int = 10

type WorkerPool struct {
	ctx     context.Context
	workers []*w.Worker // слайс воркеров, который заполняется при запуске программы,
	// и из которого потом заполняется пул
	pool chan *w.Worker // cам пул воркеров, из которого берутся воркеры, когда появляется новая задача
	// и возвращаются, когда задача отработала
	queueCtx    context.Context
	queueCancel context.CancelFunc
	queue       chan func()     // очередь, куда добавляются задачи, если нет свободных воркеров
	wg          *sync.WaitGroup // WaitGroup нужна для отслеживания активных задач,
	// чтобы при выключении данные никуда не пропали
}

// Функция NewWorkerPool создает новый объект WorkerPool,
// заполняет его и запускает прослушивание очереди
func NewWorkerPool(ctx context.Context, numberOfWorkers int) *WorkerPool {
	queueCtx, queueCancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		ctx:         ctx,
		workers:     w.NewWorkers(numberOfWorkers),
		pool:        make(chan *w.Worker, numberOfWorkers),
		queue:       make(chan func(), QueueSize),
		queueCtx:    queueCtx,
		queueCancel: queueCancel,
		wg:          &sync.WaitGroup{},
	}

	wp.fillPool()
	wp.listenQueue()

	return wp
}

// Submit - добавить таску в воркер пул
func (wp *WorkerPool) Submit(task func()) {
	// мы используем блокирующий select,
	// потому что мы должны дождаться двух случаев:
	// 1. воркер освобождается и начинает асинхронно выполнять работу
	// 2. если свободных воркеров нет, мы кладем задачу в очередь
	// в итоге первым отработает то, что первым освободится: либо воркер, либо место в очереди
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
	// здесь используется блокирующий select, чтобы программа не зависла при отмене контекста
	select {
	case <-wp.ctx.Done():
		return
	case w := <-wp.pool:
		// здесь мы не добавляем задачу в очередь и не выполняем ее асинхронно,
		// потому что мы должны дождаться окончания ее выполнения
		wp.doTaskSync(w, task)
	}
}

// Stop - остановить воркер пул, дождаться выполнения только тех тасок, которые выполняются сейчас
func (wp *WorkerPool) Stop() {
	fmt.Println("Stopping worker pool...")

	// отменяем контекст очереди для немедленного завершения ее работы
	wp.queueCancel()

	// используем WaitGroup по назначению - ожидаем, пока задачи завершатся
	wp.wg.Wait()

	// здесь мы не вычитываем задачи из очереди,
	// потому что в этой функции они нас не интересуют

	wp.cleanPool()

	fmt.Println("Worker pool was stopped.")
}

// StopWait - остановить воркер пул, дождаться выполнения всех тасок, даже тех, что не начали выполняться, но лежат в очереди
func (wp *WorkerPool) StopWait() {
	fmt.Println("Stopping worker pool...")

	// закрываем очередь, чтобы новые задачи в нее не поступали,
	// а функция listenQueue завершила свое выполнение
	close(wp.queue)

	// используем WaitGroup по назначению - ожидаем, пока задачи завершатся
	wp.wg.Wait()

	// завершаем контекст очереди для освобождения ресурсов
	wp.queueCancel()

	wp.cleanPool()

	fmt.Println("Worker pool was stopped.")
}

// Функция fillPool берет воркеров из слайса и добавляет их в пул
func (wp *WorkerPool) fillPool() {
	for _, w := range wp.workers {
		wp.pool <- w
	}
}

// Функция cleanPool вычитывает воркеров из пула, чтобы они не обрабатывали новые задачи
func (wp *WorkerPool) cleanPool() {
	for range wp.workers {
		<-wp.pool
	}
}

// Функция doTaskSync принимает в себя воркера и задачу
// и выполняет ее в синхронном режиме
func (wp *WorkerPool) doTaskSync(worker *w.Worker, task func()) {
	wp.wg.Add(1)
	defer wp.wg.Done()

	task()

	wp.pool <- worker
}

// Функция doTaskAsync принимает в себя воркера и задачу
// и выполняет ее в синхронном режиме
func (wp *WorkerPool) doTaskAsync(worker *w.Worker, task func()) {
	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()

		task()

		wp.pool <- worker
	}()
}

// Функция listenQueue запускает горутину, которая бесконечно читает из очереди
func (wp *WorkerPool) listenQueue() {
	go func() {
		for {
			select {
			case <-wp.queueCtx.Done():
				fmt.Println("Stopping queue listening due to queue context cancellation")
				return

			case task, ok := <-wp.queue:
				if !ok {
					fmt.Println("Queue is closed!")
					return
				}

				fmt.Println("Got new task from queue.")

				// блокируемся до тех пор, пока не освободится какой-нибудь воркер
				w := <-wp.pool

				wp.doTaskAsync(w, task)
			}
		}
	}()
}
