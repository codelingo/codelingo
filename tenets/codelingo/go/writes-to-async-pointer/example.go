package main

import (
	"fmt"
)

type thing struct {
	name string
}

func getNameFromPointer(t *thing) {
	fmt.Println(t.name)
}

func getName(t thing) {
	fmt.Println(t.name)
}

func setNameOnPointer(t *thing) {
	t.name = "car"
}

func newThingPointer(name string) *thing {
	return &thing{name}
}

func newThing(name string) thing {
	return thing{name}
}

func main() {

	t := newThingPointer("Pointer")

	getNameFromPointer(t)

	go getNameFromPointer(t)

	go setNameOnPointer(t)

	t2 := newThing("Thing")

	getName(t2)

	go getName(t2)

	go func(t *thing) {
		fmt.Println(t.name)
	}(t)

	go func(t thing) {
		fmt.Println(t.name)
	}(t2)

	go func(t *thing) {
		t.name = "car"
	}(t)

	go func(t thing) {
		t.name = "bus"
	}(t2)
}
