package task

import (
	"fmt"
	"sync"
	"time"
)

var (
	// TaskNum is a global counter to demonstrate task uniqueness.
	// NOTE: This global state with a mutex creates a bottleneck and is generally
	// an anti-pattern for concurrent task processing as it serializes execution.
	// It is used here for demonstration purposes only.
	TaskNum = 0
	mu      = &sync.Mutex{}
)

// Task represents a simple, stateful unit of work.
//
// It prints its unique number, increments a global counter, and sleeps to
// simulate work. The use of a global mutex makes this function safe for
// concurrent use but inefficient.
func Task() {
	mu.Lock()
	fmt.Printf("i'm task %d\n", TaskNum)
	TaskNum++
	mu.Unlock()

	time.Sleep(500 * time.Millisecond)
}
