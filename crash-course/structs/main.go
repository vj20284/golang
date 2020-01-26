package main

import "fmt"

type Person struct {
	firstName string
	lastName  string
	age		  int
}

func (p Person) greet() string {
	return "Hello " + p.firstName + "!!"
}

func (p* Person) changeAge() {
	p.age++
}


func main() {
	person1 := Person{firstName:"Philip", lastName:"Lahm", age:35}
	fmt.Println(person1)

	person2 := Person{"Bastian", "Schweinsteiger", 33}
	fmt.Println(person2)

	fmt.Println(person1.greet())
	fmt.Println(person2.greet())

	person1.changeAge()
	fmt.Println(person1)
}
