package generator

import (
	"context"
	"sandbox-team/task"
)

// Generate creates and returns a channel of task functions.
//
// It spawns a goroutine that continuously sends tasks to the channel
// until the provided context is done. The output channel is closed
// once the context is cancelled and the generator goroutine exits.
func Generate(ctx context.Context) chan func() {
	out := make(chan func())

	go func() {
		// defer ensures the channel is closed on exit, whether through the loop
		// break or a potential future panic.
		defer close(out)
	loop:
		for {
			select {
			case <-ctx.Done():
				// Context is cancelled, break the loop to exit.
				break loop
			case out <- task.Task:
				// Send a new task to the channel.
			}
		}
	}()

	return out
}
