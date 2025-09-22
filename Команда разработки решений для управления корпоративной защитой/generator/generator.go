// Пакет generator используется для генерации задач
package generator

import (
	"context"
	"fmt"
	"kaspersky/tasks"
	"math/rand/v2"
	"time"
)

const maxTimeToSleep = 1000 // миллисекунд

// Функция GenerateTasks используется для генерации задач.
// Контекст используется для остановки работы горутины
func GenerateTasks(ctx context.Context) chan func() {
	out := make(chan func())

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping generator...")
				return

			case out <- tasks.ChooseRandom():
				// randomSleep()
			}
		}
	}()

	return out
}

// Функция randomSleep используется для симуляции задержки между поступлениями задачи
func randomSleep() {
	timeToSleep := rand.IntN(maxTimeToSleep)

	time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
}
