package main

import "fmt"
import "math"

type badGeometry interface {
	badArea() float64
	badPerim() float64
}

type badRect struct {
	width, height float64
}
type badCircle struct {
	radius float64
}

func (r badRect) badArea() float64 {
	return r.width * r.height
}
func (r badRect) badPerim() float64 {
	return 2*r.width + 2*r.height
}

func (c badCircle) badArea() float64 {
	return math.Pi * c.radius * c.radius
}
func (c badCircle) badPerim() float64 {
	return 2 * math.Pi * c.radius
}

func badMeasure(g badGeometry) {
	fmt.Println(g)
	fmt.Println(g.badArea())
	fmt.Println(g.badPerim())
}

//func main() {
//	r := badRect{width: 3, height: 4}
//	c := badCircle{radius: 5}
//
//	badMeasure(r)
//	badMeasure(c)
//}
