package main

import "fmt"

func main() {
	ints()
	floats()
}

func ints() {

	a := 3
	for a < 10 { // Issue
		fmt.Println(a)
	}

	b := 4
	for b < 5 {
		fmt.Println(b)
		b++
	}

	for c := 2; c < 8; c++ {
		fmt.Println(c)
	}

	d := 5
	for 6 > d { // Issue
		fmt.Println(d)
	}

	e := 7
	for e > 2 {
		fmt.Println(e)
		e--
	}

	for f := 6; f > 1; f-- {
		fmt.Println(f)
	}

	var g int
	g = 3
	for g < 8 { // Issue
		fmt.Println(g)
	}

	var h int = 2
	for h < 9 { // Issue
		fmt.Println(h)
	}

	var j int8 = 3
	for j < 8 { // Issue
		fmt.Println(j)
	}

	var k int32
	k = 2
	for k < 7 { // Issue
		fmt.Println(k)
	}

	var l uint = 2
	for l < 9 { // Issue
		fmt.Println(l)
	}
}

func floats() {

	pi := 3.1415
	for pi < 6 { // Issue
		fmt.Println(pi)
	}

	e := 2.7182
	for e < 6 {
		fmt.Println(e)
		e++
	}

	for a := 2.2; a < 5; a++ {
		fmt.Println(a)
	}

	var b float32 = 3.3
	for b < 8 { // Issue
		fmt.Println(b)
	}

	var c float64 = 4.4
	for c < 7 {
		fmt.Println(c)
		c++
	}
}
