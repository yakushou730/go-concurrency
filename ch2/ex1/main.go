package main

import (
	"fmt"
	"time"
)

func printSomething(s string) {
	fmt.Println(s)
}

func main() {
	go printSomething("this is the first thing to printed!")

	time.Sleep(1 * time.Second)

	printSomething("this is the second thing to printed!")
}
