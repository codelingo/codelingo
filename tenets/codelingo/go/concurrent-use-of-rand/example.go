package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {

	source := rand.New(rand.NewSource(43))

	var mu sync.Mutex

	go func(r *rand.Rand) {
		fmt.Println(r.Int31())
	}(source) // Issue

	go func() {
		safeSource := rand.New(rand.NewSource(11))
		fmt.Println(safeSource.Int31())
	}()

	go func() {
		fmt.Println(source.Int31()) //Issue
	}()

	go func(r *rand.Rand) {

		mu.Lock()
		fmt.Println(r.Int31())
		mu.Unlock()

	}(source)

	printRand(source)

	go printRand(source) // Issue

	printRandUniqueSource()

	go printRandUniqueSource()

	sourceTwo := getSource(31)
	fmt.Println(sourceTwo.Int31())

	printRand(sourceTwo)

	go printRand(sourceTwo) // Issue (Requires use of CLQL types to catch)

	go func(r *rand.Rand) {

		mu.Lock()
		fmt.Println(r.Int31())
		mu.Unlock()
	}(sourceTwo) // Issue (Requires callgraph to identify sourceTwo as it is returned by an function)
}

func printRand(r *rand.Rand) {

	fmt.Println(r.Int31())
}

func printRandUniqueSource() {

	source := rand.New(rand.NewSource(11))
	fmt.Println(source.Int31())
}

func getSource(seed int64) (*rand.Rand) {

	return rand.New(rand.NewSource(seed))
}
