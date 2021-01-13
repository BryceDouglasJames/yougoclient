package yougoclient

import (
	"fmt"
)

func search(finished chan bool) {
	//fmt.Println("Worker: Started")
	NewClient()
	//fmt.Println("Worker: Finished")
	finished <- true
}

func adduser(finished chan bool, name string) {
	AmountOfUsers++
	fmt.Println("Adding new user..." + name)
	NewUser(name)
	fmt.Println("User added")
	finished <- true
}
