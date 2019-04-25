package main

import (
	"fmt"
)

func main() {
	a := 1
	if a == 1 {
		fmt.Println(a)
	}
	if a > 1 { // ISSUE
		fmt.Println(a)
	}

	// Non issue as == and >= are not mutually exclusive
	b := 1
	if b == 1 {
		fmt.Println(b)
	}
	if b >= 1 {
		fmt.Println(b)
	}

	// Non issue as rhs variables are different
	c := 1
	if c == 1 {
		fmt.Println(c)
	}
	if c > 12 {
		fmt.Println(c)
	}

	// Non issue as lhs variables are different
	d := 1
	if b == 1 {
		fmt.Println(d)
	}
	if d > 1 {
		fmt.Println(d)
	}

	// Non issues as the second if is already nested in an else
	e := 1
	if e == 1 {
		fmt.Println(e)
	} else if e > 1 {
		fmt.Println(e)
	}
}