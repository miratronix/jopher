package lib

import (
	"github.com/gopherjs/gopherjs/js"
	"errors"
	"reflect"
)

// CallWithResultCallback calls a function in the supplied JS object, appending a callback to the
// argument list and returning once that is called.
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

// CallWithErrorCallback calls a function in the supplied JS object, with the supplied arguments. It
// appends a callback to the end of the argument list that accepts an error.
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

// Require requires a node.js package
func Require(module string) *js.Object {
	return js.Global.Call("require", module)
}

// ToGoError translates a javascript error to a go error
func ToGoError(jsError *js.Error) error {
	if jsError == nil || jsError.Object == nil || jsError.Object == js.Undefined || jsError.String() == "null" {
		return nil
	}
	return errors.New(jsError.String())
}

// IsFunction determines if the supplied JS object is a function
func IsFunction(object *js.Object) bool {
	return reflect.TypeOf(object.Interface()).Kind() == reflect.Func
}

// IsArray determines if the supplied JS object is an array
func IsArray(object *js.Object) bool {
	return reflect.TypeOf(object.Interface()).Kind() == reflect.Slice
}

// ForEach iterates over the keys in a JS object
func ForEach(object *js.Object, iterator func(key string, value *js.Object)) {
	js.Global.Get("Object").Call("keys", object).Call("forEach", func(key string) {
		iterator(key, object.Get(key))
	})
}

// ToSlice converts a JS object to a slice
func ToSlice(array *js.Object) []interface{} {
	return array.Interface().([]interface{})
}

// ToMap converts a JS object to a map
func ToMap(object *js.Object) map[string]interface{} {
	return object.Interface().(map[string]interface{})
}

// HasKey determines if a object has a key
func HasKey(object *js.Object, key string) bool {
	return object.Call("hasOwnProperty", key).Bool()
}
