package main

import ("fmt")

func main() {
	a := []string{"hello", "bye"}
	_ = thing(a)
	b := []int{1, 2}
	another(b)

}

func thing(a []string) []string {

	return a
}

func another(b []int) {
	fmt.Println(b)
}
