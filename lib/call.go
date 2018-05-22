package lib

import (
	"github.com/gopherjs/gopherjs/js"
	"errors"
)

// CallWithResultCallback calls a function in the supplied JS object, appending a callback to the
// argument list and returning once that is called.
func CallWithResultCallback(jsObject *js.Object, fn string, args ...interface{}) (*js.Object, error) {

	// Define a callback result structure
	type callbackResult struct {
		result *js.Object
		err *js.Error
	}

	resultChannel := make(chan *callbackResult)

	// Add a callback at the end
	args = append(args, func(err *js.Error, result *js.Object) {
		resultChannel <- &callbackResult{result, err}
	})

	// Call the function
	jsObject.Call(fn, args...)

	result := <-resultChannel
	return result.result, ToGoError(result.err)
}

// CallWithErrorCallback calls a function in the supplied JS object, with the supplied arguments. It
// appends a callback to the end of the argument list that accepts an error.
func CallWithErrorCallback(jsObject *js.Object, fn string, args ...interface{}) error {
	errorChannel := make(chan *js.Error)

	// Add a callback at the end
	args = append(args, func(err *js.Error) {
		errorChannel <- err
	})

	jsObject.Call(fn, args...)

	err := <-errorChannel
	return ToGoError(err)
}

// ToGoError translates a javascript error to a go error
func ToGoError(jsError *js.Error) error {
	if jsError == nil || jsError.Object == nil || jsError.String() == "null" {
		return nil
	}
	return errors.New(jsError.String())
}
