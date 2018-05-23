package jopher

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/jopher/lib"
)

// Promisify converts a function to a version that returns a promise
func Promisify(function interface{}) func(args ...interface{}) *js.Object {
	reflected := lib.ReflectFunction(function)

	return func(args ...interface{}) *js.Object {
		promise := lib.NewPromise()
		go func() {

			// Reject on panic
			defer lib.CallOnPanic(promise.Reject)

			// Call the reflected function
			result, err := lib.CallReflected(reflected, args...)
			if err != nil {
				promise.Reject(err.Error())
				return
			}

			promise.Resolve(result)
		}()
		return promise.JS
	}
}

// NewPromise creates a new promise
func NewPromise(function func(resolve func(interface{}), reject func(interface{}))) *js.Object {
	promise := lib.NewPromise()
	go func() {

		// Reject on panic
		defer lib.CallOnPanic(promise.Reject)

		// Call the function, allowing the user to resolve/reject
		function(promise.Resolve, promise.Reject)
	}()
	return promise.JS
}

// Resolve creates a new resolved promise
func Resolve(value interface{}) *js.Object {
	p := lib.NewPromise()
	p.Resolve(value)
	return p.JS
}

// Reject creates a new rejected promise
func Reject(value interface{}) *js.Object {
	p := lib.NewPromise()
	p.Resolve(value)
	return p.JS
}

// CallWithResultCallback calls a function in the supplied JS object, appending a callback to the
// argument list and returning once that is called.
var CallWithResultCallback = lib.CallWithResultCallback

// CallWithErrorCallback calls a function in the supplied JS object, with the supplied arguments. It
// appends a callback to the end of the argument list that accepts an error.
var CallWithErrorCallback = lib.CallWithErrorCallback

// ToGoError translates a javascript error to a go error
var ToGoError = lib.ToGoError
