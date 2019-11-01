package main

import "pointer-leaves-method/race"

type Bar struct {
	name *string
}

type Foo struct {
	id   int
	name string
	bar  *Bar
}

func main() {
	race.Race()
}

func NewBar(name string) *Bar {
	return &Bar{
		name: &name,
	}
}

func (b *Bar) NewFooFromBar() *Foo {
	f := &Foo{
		id:   1,
		name: *b.name,
		bar:  b, // ISSUE
	}
	return f
}
