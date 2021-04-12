package test

import "fmt"

type Person struct {
	Name string
}

func hello(person Person) {
	fmt.Printf("hello, %s\n", person.Name)
}
