package generator

import (
	"context"
	"sandbox-team/task"
)

func Generate(ctx context.Context) chan func() {
	out := make(chan func())

	go func() {
		defer close(out)
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case out <- task.Task:
			}
		}
	}()

	return out
}
