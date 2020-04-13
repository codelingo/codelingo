package sample

import "fmt"

func good() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	for i := 10; i >= 0; i-- {
		fmt.Println(i)
	}
}

func bad() {
	for i := 0; i > 10; i++ {
		fmt.Println(i)
	}
	for i := 0; i <= 10; i-- {
		fmt.Println(i)
	}
	for i := uint(10); i <= 10; i-- {
		// This is... correct, but confusing.
		fmt.Println(i)
	}
}
