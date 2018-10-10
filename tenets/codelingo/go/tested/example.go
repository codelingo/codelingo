package main

import (
	"testing"
)

type aStruct struct {
	b *bStruct
}

// No test function
func (a *aStruct) SayA() string {
	if a.b == nil {
		return "end"
	}
	return "I'm a " // + a.b.SayB()
}

type bStruct struct {
	a *aStruct
}

func (b *bStruct) SayB() string {
	if b.a == nil {
		return "end"
	}
	return "I'm b " // + b.a.SayA()
}
