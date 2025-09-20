package pool

import (
	"fmt"
	"sync"
)

type Pool struct {
	tasks     chan func()
	hook      func(int)
	isStopped bool
	mu        *sync.Mutex
	wg        *sync.WaitGroup
}

func New(numWorkers int, queueSize int, hook func(workerId int)) *Pool {
	pool := &Pool{
		tasks: make(chan func(), queueSize),
		hook:  hook,
		mu:    &sync.Mutex{},
		wg:    &sync.WaitGroup{},
	}

	pool.Run(numWorkers)

	return pool
}

func (p *Pool) Run(numWorkers int) {
	for workerId := range numWorkers {
		go func() {
			fmt.Printf("worker %d started\n", workerId+1)

			for task := range p.tasks {
				p.doTask(workerId, task)
			}

			fmt.Printf("worker %d stopped\n", workerId+1)
		}()
	}
}

func (p *Pool) Submit(task func()) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isStopped {
		return fmt.Errorf("pool is stopped")
	}

	select {
	case p.tasks <- task:
		p.wg.Add(1)
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (p *Pool) Stop() error {
	fmt.Println("stopping pool...")

	p.mu.Lock()
	p.isStopped = true
	p.mu.Unlock()

	close(p.tasks)
	p.wg.Wait()

	fmt.Println("worker pool was stopped")
	return nil
}

func (p *Pool) doTask(workerId int, task func()) {
	defer p.wg.Done()

	task()
	p.hook(workerId)
}
