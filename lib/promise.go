package lib

import "github.com/gopherjs/gopherjs/js"

type state int

const (
	pending state = iota
	resolved
	rejected
)

// Promise defines a promise structure
type Promise struct {
	state              state
	value              interface{}
	resolved, rejected callback
	child              *Promise
}

// Resolves resolves the promise, calling the resolved callback
func (p *Promise) Resolve(value interface{}) {
	p.setState(resolved, value)

	// No resolved callback to make
	if p.resolved == nil {
		return
	}

	// If the callback panicked, set the state to rejected and save the error
	defer p.handlePanic()

	// Call the callback and save the value
	p.setState(resolved, p.resolved(value))
}

// Reject rejects the promise, calling the rejected callback
func (p *Promise) Reject(err interface{}) {
	p.setState(rejected, err)

	// No rejected callback to make
	if p.rejected == nil {
		return
	}

	// If the callback panicked, set the state to rejected and save the error
	defer p.handlePanic()

	// Call the callback and save the value
	p.setState(resolved, p.rejected(err))
}

// Then attaches success and failure callbacks to the promise. Both are optional.
func (p *Promise) Then(success callback, failure callback) *Promise {
	p.resolved = success
	p.rejected = failure

	// If we already resolved or rejected, call resolve/reject again to fire the callbacks
	if p.state == resolved {
		p.Resolve(p.value)
	} else if p.state == rejected {
		p.Reject(p.value)
	}

	// Create a child promise that we'll return
	p.child = &Promise{}

	// If our state is resolved or rejected after calling the callbacks, call things in the child
	if p.state == resolved {
		p.child.Resolve(p.value)
	} else if p.state == rejected {
		p.child.Reject(p.value)
	}

	return p.child
}

// Catch attaches a error handler on the promise
func (p *Promise) Catch(failure callback) *Promise {
	return p.Then(nil, failure)
}

// JS creates a JS wrapper object for this promise that includes the 'then', 'catch', 'resolve', and
// 'reject' methods
func (p *Promise) JS() *js.Object {
	wrapper := js.MakeWrapper(p)

	wrapper.Set("then", func(success *js.Object, failure *js.Object) *js.Object {
		return p.Then(newCallback(success), newCallback(failure)).JS()
	})

	wrapper.Set("catch", func(failure *js.Object) *js.Object {
		return p.Catch(newCallback(failure)).JS()
	})

	return wrapper
}

// handlePanic recovers from a panic and set the promise state/value appropriately
func (p *Promise) handlePanic() {
	if err := recover(); err != nil {
		p.setState(rejected, err)
	}
}

// setState sets the promise state
func (p *Promise) setState(s state, value interface{}) {
	p.state = s
	p.value = value
}
