package main

import "fmt"

func squares(c chan int) {
	for i := 0; i < 4; i++ {
		num := <-c
		fmt.Println(num * num)
	}
}

func main() {
	c := make(chan int, 3)

	go squares(c)

	c <- 1
	c <- 2
	c <- 3
	c <- 4
	c <- 5

	go squares(c)

	c <- 6
	c <- 7
	c <- 8
}