package main

import "fmt"
import "math"

type wellPlacedGeometry interface {
	area() float64
	perim() float64
}

func measureFoo(g wellPlacedGeometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}
