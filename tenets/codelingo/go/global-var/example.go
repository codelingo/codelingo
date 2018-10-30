package code

import (
	"errors"
)

const c1 = 0

const (
	c2 = 1
	c3 = 2
)

var v1 = 0
var v2, v3 = 1, 2
var _ = 0

func someCode() bool {
	v4 := true
	return v4
}

var errOne = errors.New("myErrVar")

var errTwo error = errors.New("declaredErr")

type customError struct{ e string }

func (e *customError) Error() string { return e.e }

// These should also be detected as global variables
var (
	v5                = 1
	errDeclared error = errors.New("errDeclared")
)
