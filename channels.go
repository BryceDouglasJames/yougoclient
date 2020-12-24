package yougoclient

import (
	"fmt"
)

func worker(finished chan bool) {
	fmt.Println("Worker: Started")
	NewClient()
	fmt.Println("Worker: Finished")
	finished <- true
}
