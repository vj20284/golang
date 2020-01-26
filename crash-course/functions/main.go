package main

import "fmt"

func greeting(s string) string {
	return "Hello " + s
}

func sum(i int, j int) int {
	return i + j
}

func multiply(i, j int) int {
	return i * j
}

func main() {
	fmt.Println(greeting("Vivek"))
	fmt.Println(sum(2, 2))
	fmt.Println(multiply(2, 10))
}