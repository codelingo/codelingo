package main

import (
	"fmt"
)

func test1() {
	x := true

	if x == true {
		fmt.Println("Hello from if")
		return
	} else { //ISSUE
		fmt.Println("Hello from else")
	}
}

func test2() {
	x := false

	if x == true {
		fmt.Println("Hello from if")
		return
	}
	fmt.Println("Hello after if")

}

func main() {
	test1()
	test2()
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		} else { //ISSUE
			fmt.Println(n)
		}
	}

	if true {
		goto end
	} else { //ISSUE
		fmt.Println("test")
	}
end:
	fmt.Println("Hello world!")

}
