package main

import "fmt"

func print(pi *int) { fmt.Println(*pi) }

func main() {

	for i := 0; i < 10; i++ {
		fmt.Println(i)
		defer fmt.Println(i)
		func() { fmt.Println(i) }()
		defer func() { fmt.Println(i) }()
		defer func(i int) { fmt.Println(i) }(i)
		defer print(&i)
		print(&i)
		go fmt.Println(i)
		go func() { fmt.Println(i) }()
	}

	// TODO: cover range..
	// nums := []int{2, 3, 4}

	// for _, num := range nums {
	// 	fmt.Println(num)
	// 	defer fmt.Println(num)
	// 	func() { fmt.Println(num) }()
	// 	defer func() { fmt.Println(num) }()
	// 	defer func(i int) { fmt.Println(num) }(num)
	// 	defer print(&num)
	// 	print(&num)
	// 	go fmt.Println(num)
	// 	go func() { fmt.Println(num) }()
	// }

	// for i, num := range nums {
	// 	fmt.Println(num)
	// 	fmt.Println(i)
	// 	defer fmt.Println(i)
	// 	func() { fmt.Println(i) }()
	// 	defer func() { fmt.Println(i) }()
	// 	defer func(i int) { fmt.Println(i) }(i)
	// 	defer print(&i)
	// 	print(&i)
	// 	go fmt.Println(i)
	// 	go func() { fmt.Println(i) }()
	// }
}
