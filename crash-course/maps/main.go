package main

import "fmt"

func main() {
	emails := make(map[string]string)

	emails["Bob"] = "bob@gmail.com"
	emails["Alice"] = "alice@gmail.com"

	fmt.Println(emails)

	emails["Sharon"] = "sharon@gmail.com"

	fmt.Println(emails)

	delete(emails, "Sharon")
	fmt.Println(emails)
}
