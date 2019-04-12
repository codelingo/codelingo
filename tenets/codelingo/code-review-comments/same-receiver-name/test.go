package main

import "fmt"

type T1 struct{}
type T2 struct{}

func (t T1) Method1() {
	fmt.Println("Hello from Method1")
}

func (this T1) Method2() {
	fmt.Println("Hello from Method2")
}

func (t T2) Method3() {
	fmt.Println("Hello from Method3")
}

func (t T2) Method4() {
	fmt.Println("Hello from Method4")
}

func main() {
	t1 := T1{}
	t1.Method1()
	t1.Method2()
	t2 := T2{}
	t2.Method3()
	t2.Method4()

}
