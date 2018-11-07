package main

func testFunc() (err string) {
	if true {
		result, err := funcThatReturnsErr()
		if err != nil || result != "it worked" {
			err = "ERROR!!!"
		}
	}
	// Will return the err from the return param, not the internal err
	// which could be confusing. Internal 'err' should be renamed to avoid confusion.
	return err
}

func anotherTestFunc(strA string) string {
	if true {
		strB := "hello"
		if true {
			strA := "hello hello"
			if true {
				strA := "goodbye"
			}
		}
	}

	return strA
}

func goodFunc(strA string) string {
	if true {
		strB := "hello"
		if true {
			strC := "hello hello"
			if true {
				strD := "goodbye"
			}
		}
	}

	return strA
}
