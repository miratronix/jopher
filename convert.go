package jopher

import (
	"errors"
	"github.com/gopherjs/gopherjs/js"
)

// ToGoError converts a javascript error object to a Go error.
func ToGoError(jsError *js.Error) error {
	if jsError == nil || jsError.Object == nil || jsError.Object == js.Undefined || jsError.String() == "null" {
		return nil
	}
	return errors.New(jsError.String())
}

// ToSlice converts a javascript object to a slice.
func ToSlice(array *js.Object) []interface{} {
	return array.Interface().([]interface{})
}

// ToMap converts a javascript object to a map.
func ToMap(object *js.Object) map[string]interface{} {
	return object.Interface().(map[string]interface{})
}

// ToString converts a javascript object to a string.
func ToString(object *js.Object) string {
	str := object.String()
	if str == "[object Object]" {
		str = js.Global.Get("JSON").Call("stringify", object).String()
	}
	return str
}
