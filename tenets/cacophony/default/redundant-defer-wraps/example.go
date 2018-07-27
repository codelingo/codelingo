package main

import (
	"log"
)

func main() {
	correct()
	incorrect()
}

type Foo struct {
	bar int
}

func (f *Foo) Print() {
	log.Printf("bar: %d", f.bar)
}

func (f *Foo) Set(bar int) {
	f.bar = bar
}

func correct() {
	f := &Foo{4}
	defer log.Println("correct single line")
	defer func() {
		log.Println("correct multiline 1")
		log.Println("correct multiline 2")
	}()
	defer func() {
		f.bar = 24
	}()
	defer func() {
		f.bar = 12
		log.Println("correct assign multiline")
	}()
	defer f.Print()
	defer f.Set(8)
}

func incorrect() {
	f := &Foo{4}
	defer func() {
		log.Println("incorrect single line") // ISSUE
	}()
	defer func() {
		f.Print() // ISSUE
	}()
}
