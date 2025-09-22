package worker

// структура воркера
// в будущем можно добавить кол-во выполненных задач определенным воркером
// для статистики
type Worker struct {
	ID int
}

// Функция NewWorkers возвращает заполненный массив воркеров
func NewWorkers(numberOfWorkers int) []*Worker {
	workers := make([]*Worker, 0, numberOfWorkers)

	for workerID := range numberOfWorkers {
		workers = append(workers, &Worker{ID: workerID})
	}

	return workers
}
