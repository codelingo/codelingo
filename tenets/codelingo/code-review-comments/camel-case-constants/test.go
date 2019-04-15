package main

import "fmt"

const Pi = 3.14
const private_Pi = 3.14 //ISSUE
const Public_Pi = 3.14  //ISSUE
const PUBLIC_PI = 3.14  //ISSUE

func main() {
	const World = "World"
	fmt.Println("Hello", World)
	fmt.Println("Happy", Pi, "Day")

	const private_World = "World" //ISSUE
	fmt.Println("Hello", private_World)
	fmt.Println("Happy", private_Pi, "Day")

	const Public_World = "World" //ISSUE
	fmt.Println("Hello", Public_World)
	fmt.Println("Happy", Public_Pi, "Day")

	const PUBLIC_WORLD = "World" //ISSUE
	fmt.Println("Hello", PUBLIC_WORLD)
	fmt.Println("Happy", PUBLIC_PI, "Day")

}
