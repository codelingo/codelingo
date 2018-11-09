package main

import (
	"log"
)

func main() {
	correct1()
	correct2()
	incorrect()
}

func correct1() {
	for i := 0; i < 10; i++ {
		log.Printf("correct1 %d", i)
	}
}

func correct2() {
	for i := 0; i < 10; i++ {
		func() {
			defer log.Printf("correct2 %d", i) // FALSE POSITIVE
		}()
	}
}

func incorrect() {
	for i := 0; i < 10; i++ {
		defer log.Printf("incorrect %d", i) // ISSUE
	}
}
