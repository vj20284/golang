package main

import "fmt"

func main() {
	var fruitArr [2]string

	fruitArr[0] = "Apple"
	fruitArr[1] = "Orange"

	fmt.Println(" ", fruitArr[0], " " , fruitArr[1])

	animals := []string{"Lion", "Tiger", "Horse"}
	fmt.Println(animals)
	fmt.Println(animals[1:])
	fmt.Println(animals[:1])
}
