package main

import "fmt"

func main() {
	x := 5
	y := 7

	if x < y {
		fmt.Printf("%d is less than %d\n", x, y)
	} else {
		fmt.Printf("%d is less than %d\n", y, x)
	}

	color := "green"

	if color == "red" {
		fmt.Println("Color is red")
	} else if color == "blue" {
		fmt.Println("Color is blue")
	} else {
		fmt.Println("Color is not blue or red")
	}

	switch color {
	case "red": fmt.Println("Color is red")
	case "blue" : fmt.Println("Color is blue")
	default : fmt.Println("Color is not blue or red")
	}
}
