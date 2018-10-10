package code

import (
	"errors"
)

const constant = 0

var myVar = 0

var myTestVar = 0

func someCode() bool {
	yourVar := true
	return yourVar
}

var theVar = true

var myVar = "global"

var theVar = 2

func someCode() bool {
	return true
}

var myVar = "global"

var myVar1, myVar2 = 1, 2

var _ = 0

var myVar = 1

var errFakeErrorUnexported = 1
var ErrFakeErrorExported = 1

var myErrVar = errors.New("myErrVar")
var myVarErr = errors.New("myVarErr")
var myVarError = errors.New("myVarErr")
var customErr = customError{"customErr"}

var errUnexported = errors.New("errUnexported")
var ErrExported = errors.New("ErrExported")
var errCustomUnexported = customError{"errCustomUnexported"}
var ErrCustomExported = customError{"ErrCustomExported"}

var declaredErr error = errors.New("declaredErr")
var errDeclared error = errors.New("errDeclared")

type customError struct{ e string }

func (e *customError) Error() string { return e.e }

// These should also be detected as global variables
var (
	// Those are not errors
	myVar = 1

	errFakeErrorUnexported = 1
	ErrFakeErrorExported   = 1

	myErrVar   = errors.New("myErrVar")
	myVarErr   = errors.New("myVarErr")
	myVarError = errors.New("myVarErr")
	customErr  = customError{"customErr"}

	errUnexported       = errors.New("errUnexported")
	ErrExported         = errors.New("ErrExported")
	errCustomUnexported = customError{"errCustomUnexported"}
	ErrCustomExported   = customError{"ErrCustomExported"}

	declaredErr error = errors.New("declaredErr")
	errDeclared error = errors.New("errDeclared")
)

type customError struct{ e string }

func (e *customError) Error() string { return e.e }
