package main

func main() {
	go func() {
		for { // ISSUE - labeled simple infinite loop
		}
	}()

	go func() {
		for {
			return
		}
	}()

	go func() {
		for { // ISSUE - labeled loop with an out of scope return child
			func() {
				return
			}()
		}
	}()

	go func() {
		for {
			break
		}
	}()

	go func() { // ISSUE - non-labeled simple infinite loop
		for {
		}
	}()

	go func() {
		for {
			return
		}
	}()

	go func() { // ISSUE - non-labeled loop with an out of scope return child
		for {
			func() {
				return
			}()
		}
	}()

	go func() {
		for {
			break
		}
	}()

	go func() { // ISSUE - labeled loop with inapplicable break statement
	l:
		for {
			switch 1 {
			case 1:
				break
			}
			continue l
		}
	}()

	go func() {
	l:
		for {
			switch 1 {
			case 1:
				break l
			}
		}
	}()

	go func() { // ISSUE - labeled loop with inapplicable break statement
	l:
		for {
			select {
			case <-make(chan int):
				break
			}
			continue l
		}
	}()

	go func() {
	l:
		for {
			select {
			case <-make(chan int):
				break l
			}
		}
	}()

	go func() { // ISSUE - labeled loop with inapplicable break statement
	l:
		for {
			for {
				break
			}
			continue l
		}
	}()

	go func() { // ISSUE - labeled loop with inapplicable break statement
	l:
		for {
		m:
			for {
				break m
			}
			continue l
		}
	}()

	go func() {
	l:
		for {
			for {
				break l
			}
		}
	}()

}
