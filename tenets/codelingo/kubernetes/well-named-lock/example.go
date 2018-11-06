package main

import "fmt"

type anyStruct struct {
	lock sync.Mutex
}

type badStruct struct {
	badLock sync.Mutex
}

type MockAddresses struct {
	lockA sync.Mutex
	lockB sync.Mutex
}
