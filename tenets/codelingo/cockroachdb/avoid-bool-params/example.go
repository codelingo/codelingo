//Package main is an example package
package main

func main() {
	doThisOrThat(false)
}

func doThisOrThat(flag bool) { // Issue
	if flag {
		doThis()
	} else {
		doThat()
	}
}

func returnsTrue() bool {
	return true
}

func doThis() {}

func doThat() {}
