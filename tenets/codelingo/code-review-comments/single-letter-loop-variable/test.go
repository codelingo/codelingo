package main

import "fmt"

func main() {

	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)

	sum = 0
	for index := 0; index < 10; index++ { //ISSUE
		sum += index
	}
	fmt.Println(sum)

	sum = 0
	for index := 0; ; index++ {
		if index > 10 {
			break
		}
		sum += index
	}
	fmt.Println(sum)

	sum = 0
	index := 0
	for ; ; index++ {
		if index > 10 {
			break
		}
		sum += index
	}
	fmt.Println(sum)

	sum = 0
	index = 0
	for {
		if index > 10 {
			break
		}
		sum += index
		index++
	}
	fmt.Println(sum)

	nums := []int{2, 3, 4}
	sum = 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println(sum)

	sum = 0
	for index, num := range nums { //ISSUE
		sum += index + num
	}
	fmt.Println(sum)
}
