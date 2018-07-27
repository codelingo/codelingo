package main

func main() {
	doThisOrThat(false)
}

func doThisOrThat(flag bool) {
	if flag {
		doThis()
	} else {
		doThat()
	}
}

func doThis() {}

func doThat() {}