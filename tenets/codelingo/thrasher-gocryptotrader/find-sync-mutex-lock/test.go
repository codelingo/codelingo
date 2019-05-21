package main

import (
	"sync"
)
var m sync.Mutex
var t = false
type test struct{
	b bool
}
func test1(){
	m.Lock() //ISSUE
	t = false
	m.Unlock()	
}

func test2(){
	m.Lock() //ISSUE
	defer m.Unlock()	
	t = true
}

func main() {
	go test1()
	go test2()
}
