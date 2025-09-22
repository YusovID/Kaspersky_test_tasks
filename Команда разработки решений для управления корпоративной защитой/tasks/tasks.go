package tasks

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Task func()

var Tasks []Task = []Task{
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

func ChooseRandom() func() {
	return Tasks[rand.IntN(len(Tasks))]
}
