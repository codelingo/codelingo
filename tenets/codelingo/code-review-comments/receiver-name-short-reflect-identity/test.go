package main

import "fmt"

type TestStruct struct{}

func (t TestStruct) Method1() {
	fmt.Println("Hello from Method1")
}

func (ts TestStruct) Method2() {
	fmt.Println("Hello from Method2")
}

func (test TestStruct) Method3() { //ISSUE
	fmt.Println("Hello from Method3")
}

func (m TestStruct) Method4() { //ISSUE
	fmt.Println("Hello from Method4")
}

func (td TestStruct) Method5() { //ISSUE
	fmt.Println("Hello from Method5")
}

func main() {
	t := TestStruct{}
	t.Method1()
	t.Method2()
	t.Method3()
	t.Method4()
	t.Method5()

}
