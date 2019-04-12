package main

import "sync"

type thing struct {
	name  string
	num   int
	mutex sync.Mutex
}

func (t thing) GetName() string { // ISSUE
	return t.name
}

func (t *thing) GetNum() int {
	return t.num
}

func main() {}
