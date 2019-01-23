package jopher

import (
	"github.com/gopherjs/gopherjs/js"
)

// Stringify is a shortcut function for using the native JS JSON stringify()
func Stringify(data map[string]interface{}) (_ []byte, returnedError error) {

	// Handle errors
	defer CatchPanic(&returnedError, "failed to stringify value")

	result := js.Global.Get("JSON").Call("stringify", data).String()
	return []byte(result), nil
}

// Parse is a shortcut function for using the native JS JSON parse()
func Parse(bytes []byte) (_ map[string]interface{}, returnedError error) {
	stringBytes := string(bytes)

	// Handle errors
	defer CatchPanic(&returnedError, "failed to parse value " + stringBytes)

	result := ToMap(js.Global.Get("JSON").Call("parse", stringBytes))
	return result, nil
}
