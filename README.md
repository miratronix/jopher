# jopher [![CircleCI](https://circleci.com/gh/miratronix/jopher.svg?style=svg)](https://circleci.com/gh/miratronix/jopher) [![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/miratronix/jopher)

`jopher` provides utility functions for working with GopherJS. Various functions are exposed, from
promise wrappers to helpers for calling JS functions.

## Installation
Use dep to install:
```bash
dep ensure -add github.com/miratronix/jopher
```

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
