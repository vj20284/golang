package main

import (
	"fmt"
	"time"
)

func Publish(text string, delay time.Duration) {
	// Go routines will be run in the same process address space
	// They are async
	go func() {
		time.Sleep(delay)
		fmt.Println("BREAKING NEWS:", text)
	}() // Note the parentheses. We must call the anonymous function.
}

func main() {
	Publish("Here is the news :-)", 5*time.Second)
	fmt.Println("Hope we get the news before program exits")

	time.Sleep(10 * time.Second)
	fmt.Println("Alright now - time to leave")
}
