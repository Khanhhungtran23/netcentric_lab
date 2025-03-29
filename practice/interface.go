package main

import "fmt"

// struct define
type Engine struct {
	Power int
}

// composition define
type Car struct {
	Engine // object ?
	Brand string
}

// define interface
type Animal interface {
	Speak() string // this is a func
}

// struct dog trien khai interface Animal
type Dog struct{}

func (d Dog) Speak() string { return "Woof!" }

// struct cat trien khai interface Animal
type Cat struct{}

func (c Cat) Speak() string { return "Meow!" }

func main() {
	e := Dog{}
	c := Cat{}
	fmt.Println(c.Speak())
	fmt.Println(e.Speak())

	// using composition
	car := Car{Engine{500}, "Toyota"}
	fmt.Println(car.Power)
}
