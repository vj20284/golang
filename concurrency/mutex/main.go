package main

import (
	"fmt"
	"sync"
	"time"
)

var i int
var consumed = false

func producer(wg *sync.WaitGroup, m *sync.Mutex) {
	for index := 0; index < 10; index++ {
		m.Lock()
		i = i + 1
		fmt.Println("Produced ", i)
		consumed = false
		m.Unlock()
		for !consumed {
			time.Sleep(1 * time.Second)
		}
	}
	wg.Done()
}

func consumer(wg *sync.WaitGroup, m *sync.Mutex) {
	for index := 0; index < 10; index++ {
		for consumed {
			time.Sleep(1 * time.Second)
		}
		m.Lock()
		fmt.Println("Consumed ", i)
		consumed = true
		m.Unlock()
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex

	go producer(&wg, &m)
	wg.Add(1)
	go consumer(&wg, &m)
	wg.Add(1)

	wg.Wait()
}
