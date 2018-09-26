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

// Those are not errors
var myVar = 1

// Those are fake errors which are currently not detected
// because they start with 'err' or 'Err' and we don't
// check if such a variable implements the error interface.
var errFakeErrorUnexported = 1
var ErrFakeErrorExported = 1

// Those errors are not named correctly
var myErrVar = errors.New("myErrVar")
var myVarErr = errors.New("myVarErr")
var myVarError = errors.New("myVarErr")
var customErr = customError{"customErr"}

// Those are actual errors which should be ignored
var errUnexported = errors.New("errUnexported")
var ErrExported = errors.New("ErrExported")
var errCustomUnexported = customError{"errCustomUnexported"}
var ErrCustomExported = customError{"ErrCustomExported"}

// Those actual errors have a declared error type
var declaredErr error = errors.New("declaredErr")
var errDeclared error = errors.New("errDeclared")

type customError struct{ e string }

func (e *customError) Error() string { return e.e }

var (
	// Those are not errors
	myVar = 1

	// Those are fake errors which are currently not detected
	// because they start with 'err' or 'Err' and we don't
	// check if such a variable implements the error interface.
	errFakeErrorUnexported = 1
	ErrFakeErrorExported   = 1

	// Those errors are not named correctly
	myErrVar   = errors.New("myErrVar")
	myVarErr   = errors.New("myVarErr")
	myVarError = errors.New("myVarErr")
	customErr  = customError{"customErr"}

	// Those are actual errors which should be ignored
	errUnexported       = errors.New("errUnexported")
	ErrExported         = errors.New("ErrExported")
	errCustomUnexported = customError{"errCustomUnexported"}
	ErrCustomExported   = customError{"ErrCustomExported"}

	// Those actual errors have a declared error type
	declaredErr error = errors.New("declaredErr")
	errDeclared error = errors.New("errDeclared")
)

type customError struct{ e string }

func (e *customError) Error() string { return e.e }
