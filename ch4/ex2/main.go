package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	hunger = 3
)

// variables - philosophers
var philosophers = []string{"A", "B", "C", "D", "E"}
var wg sync.WaitGroup
var sleepTime = 1 * time.Second
var eatTime = 2 * time.Second
var thinkTime = 1 * time.Second
var orderFinished []string
var orderMutex sync.Mutex

func diningProblem(philosopher string, leftFork, rightFork *sync.Mutex) {
	defer wg.Done()

	// print a message
	fmt.Println(philosopher, "is seated.")
	time.Sleep(sleepTime)

	for i := 0; i < hunger; i++ {
		fmt.Println(philosopher, "is hungry.")
		time.Sleep(sleepTime)

		// lock both forks
		leftFork.Lock()
		fmt.Printf("\t%s picked up the fork to his left.\n", philosopher)
		rightFork.Lock()
		fmt.Printf("\t%s picked up the fork to his right.\n", philosopher)

		// print a message
		fmt.Println(philosopher, "has both forks, and is eating.")
		time.Sleep(eatTime)

		// give the philosopher some time to think
		fmt.Println(philosopher, "is thinking")
		time.Sleep(thinkTime)

		// unlock the mutexes
		rightFork.Unlock()
		fmt.Printf("\t%s put down the fork on his right.\n", philosopher)
		leftFork.Unlock()
		fmt.Printf("\t%s put down the fork on his left.\n", philosopher)
		time.Sleep(sleepTime)
	}

	// print out done message
	fmt.Println(philosopher, "is satisfied.")
	time.Sleep(sleepTime)

	fmt.Println(philosopher, "has left the table.")

	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher)
	orderMutex.Unlock()
}

func main() {
	// print info
	fmt.Println("The Dining Philosophers Problem")
	fmt.Println("-------------------------------")

	wg.Add(len(philosophers))

	forkLeft := &sync.Mutex{}

	// spawn one goroutine for each philosopher
	for _, v := range philosophers {
		// create a mutex for the right fork
		forkRight := &sync.Mutex{}
		// call a goroutine
		go diningProblem(v, forkLeft, forkRight)

		forkLeft = forkRight
	}

	wg.Wait()

	fmt.Println("The table is empty.")
	fmt.Println("-------------------")
	fmt.Printf("Order finished: %s\n", strings.Join(orderFinished, ", "))
}
