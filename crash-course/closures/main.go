package main

import "fmt"

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func counter(x int) func() int {
	count := x
	return func() int {
		count++
		return count
	}
}

func main() {
	/* sum := adder()
	for i := 0; i < 10; i++ {
		fmt.Println(sum(i))
	} */

	seq := counter(5)
	seq2 := counter(55)

	for i := 0; i < 10; i++ {
		fmt.Println(seq())
		fmt.Println(seq2())
	}
}
