package main

import (
	"fmt"
)

type Ball struct {
	Name string
}

func (b *Ball) Ping() {
	b.Name = "111"
	fmt.Println("ping")
}

func (b Ball) Pong() {
	fmt.Println("pong")
}

func MethodMethod() {
	v := Ball{}
	v.Ping()
	v.Pong()

	Ball.Pong(v)
	// Ball.Ping(v)

	(*Ball).Ping(&v)
	(*Ball).Pong(&v)
}

func (cat *Cat) Rename(name string) {
	cat.name = name
}

func (cat Cat) WhoAmI() {
	fmt.Printf("I am %v\n", cat.name)
}

type Cat struct {
	name string
}

type Animal struct {
	*Cat
}

func main() {
	cat := Cat{"Tomcat"}
	cat.WhoAmI()
	(Cat).WhoAmI(cat)

	cat.Rename("Jerry")
	(*Cat).Rename(&cat, "Abinusi")
	(*Cat).WhoAmI(&cat)

	animal := &Animal{&cat}
	(*Animal).Rename(animal, "Woooh")
	Animal.WhoAmI(*animal)

	Cat.Rename(cat, "Harry")
}
