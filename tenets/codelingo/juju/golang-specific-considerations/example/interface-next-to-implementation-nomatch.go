package mypack

type projective_geometry interface {
    area() float64
    perim() float64
}

func (c triangle) area() float64 {
    return math.Pi * c.radius * c.radius
}

func (c square) perim() float64 {
    return 2 * math.Pi * c.radius
}