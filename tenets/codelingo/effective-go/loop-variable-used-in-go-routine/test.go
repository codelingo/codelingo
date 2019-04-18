package main

import (
	"fmt"
)

func Test1(values []int) {
	for val := range values {
		go func() { //ISSUE
			fmt.Println(val)
		}()
	}
}

func Test2(values []int) {
	for val := range values {
		go func(val int) {
			fmt.Println(val)
		}(val)
	}
}

func Test3(values []int, n int) {
	for i, val := range values {
		go func(val int) { //ISSUE
			if i < n {
				fmt.Println(val)
			}
		}(val)
	}
}

func Test4(values []int, n int) {
	for i, val := range values {
		fmt.Println(i, val)
		go func() {
			fmt.Println(n)
		}()
	}
}

func Test5(values []int, n bool) {
	for val := range values {
		go func(val int, n bool) {
			if n {
				fmt.Println(val)
			}
		}(val, n)
	}
}

func Test6(values []int, n bool) {
	for val := range values {
		go func(n bool) { //ISSUE
			if n {
				fmt.Println(val)
			}
		}(n)
	}
}

func Test7(values []int) {
	for i := 0; i < 10; i++ {
		go func() { //ISSUE
			fmt.Println(values[i])
		}()
	}
}

func Test8(values []int) {
	for val := range values {
		val := val
		go func() {
			fmt.Println(val)
		}()
	}
}

func Test9(values []int, n int) {
	for i, val := range values {
		val := val
		go func() { //ISSUE
			if i < n {
				fmt.Println(val)
			}
		}()
	}
}

func Test10(values []int, n int) {
	for i, val := range values {
		val := val
		i := i
		go func() {
			if i < n {
				fmt.Println(val)
			}
		}()
	}
}

func main() {
	values := []int{2, 3, 5, 7, 11, 13}
	Test1(values)
	Test2(values)
	Test3(values, 4)
	Test4(values, 5)
	Test5(values, true)
	Test6(values, true)
	Test7(values)
	Test8(values)
	Test9(values, 4)
	Test10(values, 5)
}
