Source code for playground explanation:
```go
package main

import (
	"fmt"
)

func main() {
	// If we take the address of a value in Go multiple times, we usually expect the same result
	a := "a"
	fmt.Println(&a)
	fmt.Println(&a)
	fmt.Println()

	// Likewise if we take the address of a value at a particular index of a slice
	things := []string{"b"}
	fmt.Println(&things[0])
	fmt.Println(&things[0])
	fmt.Println()

	// But slices reallocate memory under the hood when we append to them . . .
	previousPointer := &things[0]
	pointers := []*string{previousPointer}
	for i := 1; i <= 512; i++ {
		things = append(things, "c")

		newPointer := &things[0]
		if previousPointer != newPointer {
			// Aside: this shows exactly when go reallocates. Usually at powers of 2. Interestingly, the pattern changes after 2^10.
			// fmt.Printf(
			//	"Pointer at index 0 updated from %p to %p after %d appends.\n",
			//	previousPointer,
			//	newPointer,
			//	i,
			// )
			pointers = append(pointers, newPointer)
			previousPointer = newPointer
		}
	}

	// So each element will be given a new pointer every time the slice is reallocated
	fmt.Println(`the following all point to "things[0]":`)
	for _, p := range pointers {
		fmt.Printf("%p:%s\n", p, *p)
	}
}

```