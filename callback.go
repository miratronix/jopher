package jopher

import (
	"github.com/gopherjs/gopherjs/js"
)

// CallWithResultCallback calls a function in the provided JS object, automatically attaching a
// callback parameter to the end of the argument list. Returns when the JS callback is called with
// the appropriate value or error.
func CallWithResultCallback(jsObject *js.Object, fn string, args ...interface{}) (*js.Object, error) {

	// Define a callback result structure
	type callbackResult struct {
		result *js.Object
		err    *js.Error
	}

	resultChannel := make(chan *callbackResult)

	// Add a callback at the end
	args = append(args, func(err *js.Error, result *js.Object) {
		resultChannel <- &callbackResult{result, err}
	})

	// Call the function
	jsObject.Call(fn, args...)

	// Await the callback
	result := <-resultChannel
	return result.result, ToGoError(result.err)
}

// CallWithErrorCallback calls a function in the supplied JS object with the supplied arguments,
// automatically attaching a callback to the end of the argument list that accepts an error.
func CallWithErrorCallback(jsObject *js.Object, fn string, args ...interface{}) error {
	errorChannel := make(chan *js.Error)

	// Add a callback at the end
	args = append(args, func(err *js.Error) {
		errorChannel <- err
	})

	// Call the function
	jsObject.Call(fn, args...)

	// Await the callback
	err := <-errorChannel
	return ToGoError(err)
}
