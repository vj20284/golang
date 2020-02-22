package main

import "fmt"

// roc is a Unidirectional channel so it
// can only receive from chan
func greet(roc <-chan string) {
	fmt.Println("Hello " + <-roc)
}

func main() {
	c := make(chan string)

	go greet(c)

	c <- "Vivek"
}
