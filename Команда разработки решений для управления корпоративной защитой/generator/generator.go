package generator

import (
	"context"
	"fmt"
	"kaspersky/tasks"
	"math/rand/v2"
	"time"
)

func GenerateTasks(ctx context.Context) chan func() {
	out := make(chan func())

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping generator...")
				return

			case out <- tasks.ChooseRandom():
				// RandomSleep()
			}
		}
	}()

	return out
}

func RandomSleep() {
	timeToSleep := rand.IntN(1000)

	time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
}
