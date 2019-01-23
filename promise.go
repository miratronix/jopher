package jopher

import "github.com/gopherjs/gopherjs/js"

// Promisify promisifies an existing function, returning the new version.
func Promisify(function interface{}) func(args ...interface{}) *js.Object {
	reflected := reflectFunction(function)

	return func(args ...interface{}) *js.Object {
		p := newPromise()
		go func() {

			// Reject on panic
			defer CallOnPanic(p.Reject)

			// Call the reflected function
			result, err := callReflected(reflected, args...)
			if err != nil {
				p.Reject(err.Error())
				return
			}

			p.Resolve(result)
		}()
		return p.JS
	}
}

// NewPromise constructs a new promise using a `(resolve, reject)` callback, similar to javascript.
func NewPromise(function func(resolve func(interface{}), reject func(interface{}))) *js.Object {
	p := newPromise()
	go func() {

		// Reject on panic
		defer CallOnPanic(p.Reject)

		// Call the function, allowing the user to resolve/reject
		function(p.Resolve, p.Reject)
	}()
	return p.JS
}

// Resolve returns a new promise that is resolved with the supplied value.
func Resolve(value interface{}) *js.Object {
	p := newPromise()
	p.Resolve(value)
	return p.JS
}

// Reject returns a new promise that is rejected with the supplied value.
func Reject(value interface{}) *js.Object {
	p := newPromise()
	p.Resolve(value)
	return p.JS
}

// promise defines a promise structure
type promise struct {
	JS      *js.Object
	Resolve func(interface{})
	Reject  func(interface{})
}

// newPromise constructs a new promise
func newPromise() *promise {
	p := &promise{}

	p.JS = js.Global.Get("Promise").New(func(resolve *js.Object, reject *js.Object) {
		p.Resolve = func(val interface{}) { resolve.Invoke(val) }
		p.Reject = func(val interface{}) { reject.Invoke(val) }
	})

	return p
}
