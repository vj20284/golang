package main

import (
	"fmt"
	"math"
)

type Shape interface {
	area() float64
}

type Circle struct {
	radius float64
}

type Rectangle struct {
	length, height float64
}

func (c Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (r Rectangle) area() float64 {
	return r.length * r.height
}

func getArea(s Shape) float64 {
	return s.area()
}

func main() {
	circle := Circle{5}
	rectangle := Rectangle{2,5}

	fmt.Printf("Area is %f\n", getArea(circle))
	fmt.Printf("Area is %f\n", getArea(rectangle))
}
