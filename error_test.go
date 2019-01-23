package jopher

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCallOnPanic(t *testing.T) {

	Convey("Calls the supplied function with the panic value if there is a panic", t, func() {
		ch := make(chan interface{})
		handleError := func(val interface{}) {
			ch <- val
		}

		go func() {
			defer CallOnPanic(handleError)
			panic("nope")
		}()

		So(<-ch, ShouldEqual, "nope")
	})
}
