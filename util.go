package jopher

import (
	"github.com/gopherjs/gopherjs/js"
	"reflect"
)

// Require requires a module (only works in node or if a `require` polyfill is supplied).
func Require(module string) *js.Object {
	return js.Global.Call("require", module)
}

// IsFunction determines if the supplied javascript object is a function.
func IsFunction(object *js.Object) bool {
	return reflect.TypeOf(object.Interface()).Kind() == reflect.Func
}

// IsArray determines if the supplied javascript object is an array.
func IsArray(object *js.Object) bool {
	return reflect.TypeOf(object.Interface()).Kind() == reflect.Slice
}

// ForEach iterates over the keys in a javascript object.
func ForEach(object *js.Object, iterator func(key string, value *js.Object)) {
	js.Global.Get("Object").Call("keys", object).Call("forEach", func(key string) {
		iterator(key, object.Get(key))
	})
}

// HasKey determines if the supplied javascript object has the specified key.
func HasKey(object *js.Object, key string) bool {
	return object.Call("hasOwnProperty", key).Bool()
}

// NewObject create a new javascript object.
func NewObject() *js.Object {
	return js.Global.Get("Object").New()
}

// NewISODate creates a new RFC3339 date string using a JS native call.
func NewISODate() string {
	return js.Global.Get("Date").New().Call("toISOString").String()
}
