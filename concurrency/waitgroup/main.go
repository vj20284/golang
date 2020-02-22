package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func race(wg *sync.WaitGroup, name string, speed int) {
	for i := 0; i < 10; i += speed {
		time.Sleep(1 * time.Second)
	}
	fmt.Println(name + " Finished the race")
	wg.Done() // Counter decrement
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i < 4; i++ {
		wg.Add(1) // Counter increment
		go race(&wg, strconv.Itoa(i), i)
	}

	wg.Wait() // Wait for all go routines to finish
	fmt.Println("Race is over")
}
