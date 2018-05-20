package lib

import (
	"reflect"
	"errors"
	"github.com/gopherjs/gopherjs/js"
)

// callback defines a callback function
type callback func(value interface{}) interface{}

// ReflectFunction converts a supplied interface into a reflect.Value
func ReflectFunction(function interface{}) reflect.Value {
	reflected := reflect.ValueOf(function)

	if reflected.Kind() != reflect.Func {
		panic(errors.New("please supply a function"))
	}

	return reflected
}

// CallAsync calls a function, calling the supplied resolve.reject methods afterwards as necessary
func CallAsync(resolve func(interface{}), reject func(interface{}), fn reflect.Value, args ...interface{}) {
	go func() {

		// Reflect all the arguments and call the function
		reflectedArgs := reflectAll(args...)
		results := fn.Call(reflectedArgs)

		// Determine if the function returns an error as the last return value
		hasError := hasLastError(fn.Type())

		// Split the results into a slice of interfaces and an error value
		value, err := splitResults(results, hasError)

		if err != nil {
			reject(err)
			return
		}

		resolve(value)
	}()
}

// newCallback converts a js function to a proper callback
func newCallback(jsFunction *js.Object) callback {

	// No function supplied, convert to nil
	if jsFunction == nil || jsFunction == js.Undefined {
		return nil
	}

	// Return a callback function
	return func(val interface{}) interface{} {
		return jsFunction.Invoke(val)
	}
}

// reflectAll converts the supplied arguments to reflect values
func reflectAll(args ...interface{}) []reflect.Value {
	reflected := make([]reflect.Value, len(args))

	for i := range args {
		reflected[i] = reflect.ValueOf(args[i])
	}

	return reflected
}

// unReflectAll converts the supplied reflect values to a slice of interfaces
func unReflectAll(results []reflect.Value) []interface{} {
	outs := make([]interface{}, len(results))

	for i := range results {
		outs[i] = results[i].Interface()
	}

	return outs
}

// splitResults splits a slice of results into an interface and an error. The interface could contain
// nil (if no value was returned), a single value (if a single value was returned), or a slice
// of interface{}s (if multiple values were returned).
func splitResults(results []reflect.Value, lastError bool) (interface{}, error) {
	count := len(results)

	// Fish out the error at the end
	var err error
	if lastError && count > 0 {
		var errorValue reflect.Value

		results, errorValue = results[:count-1], results[count-1]
		if errorValue.IsValid() && !errorValue.IsNil() {
			err = errorValue.Interface().(error)
		}
	}

	// Clean up the returned result
	actualResults := unReflectAll(results)
	switch len(actualResults) {
	case 0:
		return nil, err

	case 1:
		return actualResults[0], err

	default:
		return actualResults, err
	}
}

func hasLastError(t reflect.Type) bool {
	count := t.NumOut()
	if count == 0 {
		return false
	}

	return t.Out(count-1) == reflect.ValueOf((*error)(nil)).Type().Elem()
}
