package tasks

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// новый тип для возможного дальнейшего расширения функциональности программы
type Task func()

// tasks - массив для примера задач, которые работает разное время
var tasks []Task = []Task{
	func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("i'm task #1")
	},
	func() {
		time.Sleep(250 * time.Millisecond)
		fmt.Println("i'm task #2")
	},
	func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("i'm task #3")
	},
}

// Функция ChooseRandom возвращает случайную задачу из массива задач
func ChooseRandom() func() {
	return tasks[rand.IntN(len(tasks))]
}
