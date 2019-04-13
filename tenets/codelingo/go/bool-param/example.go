//Package main is an example package
package main

func main() {
	doThisOrThat(false)
}

func doThisOrThat(flag bool) bool { // Issue
	if flag {
		doThis()
	} else {
		doThat()
	}

	return false
}

func doThisOrThatNoArg() bool { // Issue
	if flag {
		doThis()
	} else {
		doThat()
	}

	return false
}

func returnsTrue() bool {
	return true
}

func doThis() (flag bool) {}

func doThat() {}

func someFuncA(a string, b bool)           {}
func someFuncB(b bool, a string)           {}
func someFuncC(a string, b bool, c string) {}
func someFuncD(b bool, c string, d bool)   {}
