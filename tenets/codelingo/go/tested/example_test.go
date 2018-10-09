package main

import (
	"testing"
)

func TestSayB(t *testing.T) {
	x := &aStruct{
		b: &bStruct{
			a: &aStruct{
				b: &bStruct{},
			},
		},
	}

	// print(x.Say())
	print(x.b.SayB())
	print(x.b.a.b.SayB())
}
