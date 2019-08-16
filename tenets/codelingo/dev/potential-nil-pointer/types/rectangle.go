package types

type Rect struct {
	Width  float64
	Height float64
}

func (r Rect) Area() float64 {
	return r.Width * r.Height
}

func (r Rect) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r *Rect) Modify() *Rect {
	if r.Area() < 20 {
		return nil
	}
	r.Height *= 2
	r.Width *= 2
	return r
}
