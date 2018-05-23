package lib

import "github.com/gopherjs/gopherjs/js"

// promise defines a promise structure
type promise struct {
	JS      *js.Object
	Resolve func(interface{})
	Reject  func(interface{})
}

// NewPromise constructs a new promise
func NewPromise() *promise {
	p := &promise{}

	p.JS = js.Global.Get("Promise").New(func(resolve *js.Object, reject *js.Object) {
		p.Resolve = func(val interface{}) { resolve.Invoke(val) }
		p.Reject = func(val interface{}) { reject.Invoke(val) }
	})

	return p
}
