package main

import "fmt"
import "math"

type goodGoodGeometry interface {
	goodGoodArea() float64
	goodGoodPerim() float64
}

type goodGoodRect struct {
	width, height float64
}
type goodGoodCircle struct {
	radius float64
}

func (r goodGoodRect) goodGoodArea() float64 {
	return r.width * r.height
}
func (r goodGoodRect) goodGoodPerim() float64 {
	return 2*r.width + 2*r.height
}

func (c goodGoodCircle) goodGoodArea() float64 {
	return math.Pi * c.radius * c.radius
}
func (c goodGoodCircle) goodGoodPerim() float64 {
	return 2 * math.Pi * c.radius
}

func goodGoodMeasure(g goodGoodGeometry) {
	fmt.Println(g)
	fmt.Println(g.goodGoodArea())
	fmt.Println(g.goodGoodPerim())
}

//func main() {
//	r := goodGoodRect{width: 3, height: 4}
//	c := goodGoodCircle{radius: 5}
//
//	goodGoodMeasure(r)
//	goodGoodMeasure(c)
//}
