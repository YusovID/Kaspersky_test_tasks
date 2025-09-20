package task

import (
	"fmt"
	"sync"
	"time"
)

var (
	TaskNum = 0
	mu      = &sync.Mutex{}
)

func Task() {
	mu.Lock()
	fmt.Printf("i'm task %d\n", TaskNum)
	TaskNum++
	mu.Unlock()

	time.Sleep(500 * time.Millisecond)
}
