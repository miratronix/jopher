# jopher [![CircleCI](https://circleci.com/gh/miratronix/jopher.svg?style=svg)](https://circleci.com/gh/miratronix/jopher) [![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/miratronix/jopher)

`jopher` provides utility functions for working with GopherJS. Various functions are exposed, from
promise wrappers to helpers for calling JS functions.

## Installation
Use dep to install:
```bash
dep ensure -add github.com/miratronix/jopher
```

## Usage

### Exposed Functions
This package exposes several utility functions:

#### Promise-related
* `Promisify(function interface{}) interface{}` - 
    Promisifies an existing function, returning the new version.
* `NewPromise(function func(resolve func(interface{}), reject func(interface{}))) *js.Object` - 
    Constructs a new promise using a `(resolve, reject)` callback, similar to javascript.
* `Resolve(value interface{}) *js.Object` - 
    Returns a new promise that is resolved with the supplied value.
* `Reject(value interface{}) *js.Object` - 
    Returns a new promise that is rejected with the supplied value.

#### Other Utilities
* `CallWithResultCallback(jsObject *js.Object, fn string, args ...interface{}) (*js.Object, error)` -
    Calls a function in the provided JS object, automatically attaching a callback parameter to the
    end of the argument list. Returns when the JS callback is called with the appropriate value
    or error.
* `CallWithErrorCallback(jsObject *js.Object, fn string, args ...interface{}) error` -
    Calls a function in the supplied JS object with the supplied arguments, automatically attaching
    a callback to the end of the argument list that accepts an error.
* `ToGoError(jsError *js.Error) error` -
    Converts a javascript error object to a Go error.

### Promises
This package exposes utility methods for dealing with promises in gopherJS. Many of them return a JS
wrapped instance of `Promise`. A `Promise` has the following the methods:
* `Then(success, failure func(interface{}) interface{}})` - 
    Adds success/failure callbacks to the promise and returns a new promise. Either of these callbacks 
    can be nil. If the success callback returns a value, it will be returned to child promises. If the 
    success callback panics, child promises will be rejected. If the failure callback panics, the 
    failure will be propagated to child promises. However, if the reject callback does not panic, it's 
    return value will be returned to child promises. This function is exposed to javascript as `then()`.
* `Catch(failure func(interface{}) interface{}})` - 
    Shortcut method for specifying only the failure callback in a `then`. This is exposed to javascript 
    as `catch()`.
* `Resolve(value interface{})` - 
    Allows for manual resolving of the promise.
* `Reject(value interface{})` - 
    Allows for manual rejecting of the promise.

## Examples

### Promisify
The simplest way to create promises is to wrap existing functions with `Promisify`:
```go
import "github.com/miratronix/jopher"

func main() {

	// As part of a global
	js.Global.Set("httpCall", jopher.Promisify(httpCall))

	// or as part of a structured object:
	js.Global.Set("api", map[string]interface{}{
		"httpCall": jopher.Promisify(httpCall),
	})
}

// This is a blocking function -- it doesn't return until the http call completes or fails.
func httpCall() (SomeResponse, error) {
	response, err := http.Get("/someAPI")
	if err != nil {
		return nil, err
	}
	return response, nil
}
```

Promisify allows JS to call the underlying function via reflection and automatically detects an 
'error' return type, using the following rules, in order:
* If the function panics, the promise is rejected with the panic value.
* If the last return is of type 'error', then the promise is rejected if the returned error is non-nil.
* The promise is resolved with the remaining return values, according to how many there are:
    * 0:  resolved with nil
    * 1:  resolved with that value
    * 2+: resolved with a slice of the values

### New Promise
You can also construct a new promise and manage the resolve/reject yourself:
```go
import "github.com/miratronix/jopher"

func main() {
	js.Global.Set("httpCall", jopher.NewPromise(httpCall))
}

// A blocking function, as before
func httpCall(resolve, reject func(interface{})) {
	response, err := http.Get("/someAPI")
	if err != nil {
		reject(err)
	}
	resolve(response)
}
```

### Resolve/Reject
For small methods that don't block, it can be useful to quickly return a promise:
```go
import "github.com/miratronix/jopher"

func main() {
	js.Global.Set("httpCall", httpCall)
}

func httpCall() *js.Object {

	// Return an immediately resolved promise
	return jopher.Resolve(1)

	// Or a rejected one
	return jopher.Reject(2)
}
```

### Manual Promise Construction
Finally, you can also construct and manage the promise manually. For example:
```go
import "github.com/miratronix/jopher"

func main() {
	js.Global.Set("httpCall", httpCall)
}

func httpCall() *js.Object {
	p := jopher.Promise{}

	go func() {
		response, err := http.Get("/someAPI")
		if err != nil {
			p.Reject(err)
		}
		p.Resolve(response)
	}()

	return p.JS()
}
```

### Calling Javascript Functions
`jopher` supplies 2 functions for calling javascript functions:
```go
import "github.com/miratronix/jopher"

// Some JS object with functions that accept a callback
var jsObject *js.Object

// Attaches a callback to the supplied argument list and returns when it's called
result, err := jopher.CallWithResultCallback(jsObject, "someFunction", "someArgument")

// Attaches a callback to the supplied argument list and returns once the callback is called
err := jopher.CallWithErrorCallback(jsObject, "someOtherFunction", "someArgument")
```
