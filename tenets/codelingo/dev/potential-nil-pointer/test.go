package main

import (
	"fmt"
	"github.com/codelingo/rectangle"
	"log"
	"potentialNilPointer/types"
)

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

func main() {
	// type `Rect` is defined in the same file
	a := Rect{5.0, 4.0}

	modifiedA1 := a.Modify()
	fmt.Println("Width", modifiedA1.Width)  // ISSUE
	fmt.Println("Hight", modifiedA1.Height) // ISSUE

	modifiedA2 := a.Modify()
	if modifiedA2 != nil {
		fmt.Println("Width", modifiedA2.Width)
		fmt.Println("Hight", modifiedA2.Height)
	}

	modifiedA3 := a.Modify()
	if modifiedA3 == nil {
		log.Fatal("nil pointer panic")
	}
	fmt.Println("Width", modifiedA3.Width)
	fmt.Println("Hight", modifiedA3.Height)

	// type `Rect` is defined in a local package
	b := types.Rect{5.0, 4.0}

	modifiedB1 := b.Modify()
	fmt.Println("Width", modifiedB1.Width)  // ISSUE
	fmt.Println("Hight", modifiedB1.Height) // ISSUE

	modifiedB2 := b.Modify()
	if modifiedB2 != nil {
		fmt.Println("Width", modifiedB2.Width)
		fmt.Println("Hight", modifiedB2.Height)
	}
	modifiedB3 := b.Modify()
	if modifiedB3 == nil {
		log.Fatal("nil pointer panic")
	}
	fmt.Println("Width", modifiedB3.Width)
	fmt.Println("Hight", modifiedB3.Height)

	// type `Rect` is defined in a github package
	c := rectangle.Rect{3.0, 4.0}

	modifiedC1 := c.Modify()
	fmt.Println("Width", modifiedC1.Width)  // ISSUE
	fmt.Println("Hight", modifiedC1.Height) // ISSUE

	modifiedC2 := c.Modify()
	if modifiedC2 != nil {
		fmt.Println("Width", modifiedC2.Width)
		fmt.Println("Hight", modifiedC2.Height)
	}
	modifiedC3 := c.Modify()

	if modifiedC3 == nil {
		log.Fatal("nil pointer panic")
	}
	fmt.Println("Width", modifiedC3.Width)
	fmt.Println("Hight", modifiedC3.Height)
}
