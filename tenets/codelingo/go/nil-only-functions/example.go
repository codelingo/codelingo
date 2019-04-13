package main

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var err error
	err = onlyReturnsNil()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = returnsMixed()
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = returnMultipleOnlyNil()
	if err != nil {
		log.Fatalf(err.Error())
	}
	_, err = returnMultipleMixed()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

type example struct {
	value int
}

func returnNum() int {
	return 3
}

func returnPointerAndNil() (*example, error) {
	return &example{46}, nil
}

func returnCallAndNil() (int, error) {
	return returnNum(), nil
}

func onlyReturnsNil() error {
	a := rand.Intn(10) + 1
	if a <= 5 {
		log.Println("error: wanted a number higher than 5")
	}
	log.Println("success: got a number higher than 5")
	return nil
}

func returnNonNil() (int, error) {
	return 3, nil
}

func returnsMixed() error {
	a := rand.Intn(10) + 1
	if a <= 5 {
		return errors.New("wanted a number higher than 5")
	}
	return nil
}

func returnMultipleOnlyNil() (*example, error) {
	ex := &example{
		value: rand.Intn(15),
	}
	if ex.value <= 5 {
		log.Println("error: wanted a number higher than 5")
	} else if ex.value >= 6 && ex.value <= 10 {
		ex.value = ex.value * 2
	}
	return nil, nil
}

func returnMultipleMixed() (*example, error) {
	ex := &example{
		value: rand.Intn(15),
	}
	if ex.value <= 5 {
		return nil, errors.New("error: wanted a number higher than 5")
	} else if ex.value >= 6 && ex.value <= 10 {
		ex.value = ex.value * 2
		return ex, nil
	}
	//return ex, errors.New("warning: uncaught value")
	return ex, nil
}

// intFunc satisfies interface X.
func (example) intFunc() error {
	return nil
}

// This is a random comment, func doesn't satisify interface.
func returnNilHasComment() error {
	return nil
}
