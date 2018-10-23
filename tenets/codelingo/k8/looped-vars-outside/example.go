package main

// Heuristic:
// ‾‾‾‾‾‾‾‾‾‾
// '&'i inside for i. also, func(){} without (i)

import "fmt"

func print(pi *int) { fmt.Println(*pi) }

func main() {
	for i := 0; i < 10; i++ { // TODO capture i (in clql)
		defer fmt.Println(i)                    // 1) OK; prints 9 ... 0
		defer func() { fmt.Println(i) }()       // 2) WRNG; prints "10" 10 times
		func() { fmt.Println(i) }()             // 2.1) WRNG; prints "10" 10 times
		func() { fmt.Println(i) }               // 2.2) WRNG; prints "10" 10 times
		defer func(i int) { fmt.Println(i) }(i) // 3) OK
		defer print(&i)                         // 4) WRONG; prints "10" 10 times
		go fmt.Println(i)                       // 5) OK; prints 0 ... 9 in unpredictable order
		go func() { fmt.Println(i) }()          // 6) WRONG; totally unpredictable.
	}

	// for key, value := range myMap {
	// 	// Same for key & value as i!
	// }
}
