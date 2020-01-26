package main

import "fmt"

func main() {
	var fruitArr [2]string

	fruitArr[0] = "Apple"
	fruitArr[1] = "Orange"

	for _, fruit := range fruitArr {
		fmt.Println(fruit)
	}

	animals := []string{"Lion", "Tiger", "Horse"}
	for _, animal := range animals {
		fmt.Println(animal)
	}

	emails := map[string]string {"Bob": "bob@gmail.com", "Alice": "alice@gmail.com"}
	for k, v := range emails {
		fmt.Println(k + " : " + v)
	}
}
