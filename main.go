package promise

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/promise/lib"
)

// Promise exposes the underlying promise structure
type Promise = lib.Promise

// Promisify converts a function to a version that returns a promise
func Promisify(function interface{}) interface{} {
	reflected := lib.ReflectFunction(function)

	return func(args ...interface{}) *js.Object {
		promise := lib.Promise{}
		lib.CallAsync(promise.Resolve, promise.Reject, reflected, args...)
		return promise.JS()
	}
}

// New creates a new promise
func New(function func(resolve func(interface{}), reject func(interface{}))) *js.Object {
	promise := lib.Promise{}
	go function(promise.Resolve, promise.Reject)
	return promise.JS()
}

// Resolve creates a new resolved promise
func Resolve(value interface{}) *js.Object {
	p := &lib.Promise{}
	p.Resolve(value)
	return p.JS()
}

// Reject creates a new rejected promise
func Reject(value interface{}) *js.Object {
	p := &lib.Promise{}
	p.Resolve(value)
	return p.JS()
}
