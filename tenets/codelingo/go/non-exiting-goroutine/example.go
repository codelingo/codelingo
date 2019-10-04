package main

func main() {
	go func() {
	l:
		for { // ISSUE - labeled simple infinite loop
		}
	}()

	go func() {
	l:
		for {
			return
		}
	}()

	go func() {
	l:
		for { // ISSUE - labeled loop with an out of scope return child
			func() {
				return
			}()
		}
	}()

	go func() {
	l:
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

}
