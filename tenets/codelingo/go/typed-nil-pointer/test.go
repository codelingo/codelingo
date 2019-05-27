package main

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func test1(arg int) interface{} {
	var result *Person = nil

	if arg > 0 {
		result = &Person{}
	}

	return result //ISSUE
}

func test2(arg int) interface{} {
	var result *Person = nil

	if arg > 0 {
		result = &Person{}
	}

	if result == nil {
		return nil
	}

	return result
}

func main() {
	fmt.Println("test1 func")
	if res := test1(-1); res != nil {
		fmt.Println("non nil result:", res)
		fmt.Println(res == nil)
		fmt.Printf("%T", res)
		fmt.Println()
	} else {
		fmt.Println("nil result", res)
	}
	fmt.Println()
	fmt.Println("test2 func")
	if res := test2(-1); res != nil {
		fmt.Println("non nil result:", res)
	} else {
		fmt.Println("nil result", res)
		fmt.Println(res == nil)
		fmt.Printf("%T", res)
	}
}
