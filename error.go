package jopher

import (
	"errors"
	"github.com/gopherjs/gopherjs/js"
)

// CatchPanic catches a panic and puts the error into the supplied error. This should be called in a
// defer.
func CatchPanic(returnedError *error, failureMessage string) {

	// Catch the panic
	err := recover()
	if err == nil {
		return
	}

	// Try to make it a JS error
	jsErr, ok := err.(*js.Error)
	if ok {
		*returnedError = jsErr
	}

	// Try to make it a go error
	goErr, ok := err.(error)
	if ok {
		*returnedError = goErr
	}

	// Not a go or JS error
	*returnedError = errors.New(failureMessage)
}

// CallOnPanic calls the supplied function when a panic is recovered. This should be called in a defer.
func CallOnPanic(reject func(interface{})) {
	err := recover()
	if err != nil {
		reject(err)
	}
}

// Throw throws the supplied javascript object.
func Throw(object *js.Object) {
	panic(object)
}

// ThrowOnError throws when the supplied error is not nil.
func ThrowOnError(err error) {
	if err != nil {
		panic(NewError(err.Error()))
	}
}

// NewError creates a new JS error.
func NewError(message string) *js.Object {
	return js.Global.Get("Error").New(message)
}
