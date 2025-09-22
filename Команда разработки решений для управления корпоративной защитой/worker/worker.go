package worker

type Worker struct {
	ID int
}

func NewWorkers(numberOfWorkers int) []*Worker {
	workers := make([]*Worker, 0, numberOfWorkers)

	for workerID := range numberOfWorkers {
		workers = append(workers, &Worker{ID: workerID})
	}

	return workers
}
