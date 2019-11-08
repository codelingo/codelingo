package main

import "fmt"

func do(i interface{}) {
	switch v := i.(type) { // Non-issue
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

}

func dont(i interface{}) {
	switch v := i.(type) { // ISSUE
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	}

}

func main() {
	do(21)
	do("hello")
	do(true)

	dont(21)
	dont("hello")
	dont(true)
}

