package main

import (
	"fmt"
)

func test(x int) {
	fmt.Println(x)
	{
		x := 2 //ISSUE
		fmt.Println(x)
	}
	fmt.Println(x)
}

func main() {
	x := 1
	test(x)
}
