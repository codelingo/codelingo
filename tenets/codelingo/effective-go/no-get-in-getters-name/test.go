package main

import (
	"fmt"
	"time"
)

type rect struct {
	width, height int
}

func (r *rect) area() int {
	return r.width * r.height
}

func (r *rect) GetWidth() int { //ISSUE
	return r.width
}

func (r *rect) GetHeight() int { //ISSUE
	return r.height
}

func (r rect) perim() int {
	return 2*r.width + 2*r.height
}

func GetTime() {
	fmt.Println("The time is :", time.Now())
}

func main() {

	r := rect{width: 10, height: 5}

	fmt.Println("width ", r.GetWidth())
	fmt.Println("height ", r.GetHeight())

	fmt.Println("area: ", r.area())
	fmt.Println("perim:", r.perim())

	GetTime()
}
