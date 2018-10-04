package main

import "fmt"
import "math"

type goodGeometry interface {
	goodArea() float64
	goodPerim() float64
}

type goodRect struct {
	width, height float64
}
type goodCircle struct {
	radius float64
}

func (r goodRect) goodArea() float64 {
	return r.width * r.height
}
func (r goodRect) goodPerim() float64 {
	return 2*r.width + 2*r.height
}

func (c goodCircle) goodArea() float64 {
	return math.Pi * c.radius * c.radius
}
func (c goodCircle) goodPerim() float64 {
	return 2 * math.Pi * c.radius
}

func goodMeasure(g goodGeometry) {
	fmt.Println(g)
	fmt.Println(g.goodArea())
	fmt.Println(g.goodPerim())
}

//func main() {
//	r := goodRect{width: 3, height: 4}
//	c := goodCircle{radius: 5}
//
//	goodMeasure(r)
//	goodMeasure(c)
//}
