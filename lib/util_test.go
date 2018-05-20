package lib

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	"reflect"
	"errors"
	"github.com/gopherjs/gopherjs/js"
)

func TestReflectFunction(t *testing.T) {

	Convey("Converts a function to a reflect.Value", t, func() {
		reflected := ReflectFunction(func() {})
		So(reflected, ShouldHaveSameTypeAs, reflect.Value{})
	})

	Convey("Panics if a non-function is supplied", t, func() {
		So(func() { ReflectFunction(0) }, ShouldPanic)
	})
}

func TestCallAsync(t *testing.T) {

	Convey("Calls the function with one supplied argument", t, func() {
		none := func(interface{}) {}
		reflected := ReflectFunction(func(ch chan bool) { ch <- true })

		ch := make(chan bool)
		CallAsync(none, none, reflected, ch)
		result := <-ch
		So(result, ShouldBeTrue)
	})

	Convey("Calls the function with two supplied arguments", t, func() {
		none := func(interface{}) {}
		reflected := ReflectFunction(func(ch1 chan bool, ch2 chan bool) {
			ch1 <- true
			ch2 <- true
		})

		ch1 := make(chan bool)
		ch2 := make(chan bool)
		CallAsync(none, none, reflected, ch1, ch2)
		result := <-ch1
		result = <-ch2
		So(result, ShouldBeTrue)
	})

	Convey("Calls resolve with nil when the function doesn't return anything", t, func() {
		none := func(interface{}) {}

		ch := make(chan interface{})
		resolve := func(value interface{}) { ch <- value }

		CallAsync(resolve, none, ReflectFunction(func() {}))
		result := <-ch
		So(result, ShouldBeNil)
	})

	Convey("Calls resolve with the value when the function returns a value", t, func() {
		none := func(interface{}) {}

		ch := make(chan interface{})
		resolve := func(value interface{}) { ch <- value }

		CallAsync(resolve, none, ReflectFunction(func() int { return 1 }))
		result := <-ch
		So(result, ShouldEqual, 1)
	})

	Convey("Calls resolve with a slice when the function returns multiple values", t, func() {
		none := func(interface{}) {}

		ch := make(chan interface{})
		resolve := func(value interface{}) { ch <- value }

		CallAsync(resolve, none, ReflectFunction(func() (int, int) { return 1, 2 }))
		result := <-ch
		So(result, ShouldResemble, []interface{}{1, 2})
	})

	Convey("Calls resolve with nil when the function returns a nil error", t, func() {
		none := func(interface{}) {}

		ch := make(chan interface{})
		resolve := func(value interface{}) { ch <- value }

		CallAsync(resolve, none, ReflectFunction(func() error { return nil }))
		result := <-ch
		So(result, ShouldBeNil)
	})

	Convey("Calls reject with the error when the function returns a error", t, func() {
		none := func(interface{}) {}

		ch := make(chan interface{})
		reject := func(value interface{}) { ch <- value }

		CallAsync(none, reject, ReflectFunction(func() error { return errors.New("nope") }))
		result := <-ch
		So(result, ShouldResemble, errors.New("nope"))
	})

	Convey("Calls reject with the error if the function returns an error as the last result", t, func() {
		none := func(interface{}) {}

		ch := make(chan interface{})
		reject := func(value interface{}) { ch <- value }

		CallAsync(none, reject, ReflectFunction(func() (int, error) { return 3, errors.New("nope") }))
		result := <-ch
		So(result, ShouldResemble, errors.New("nope"))
	})
}

func TestNewCallback(t *testing.T) {

	Convey("Returns nil is nil is supplied", t, func() {
		So(newCallback(nil), ShouldBeNil)
	})

	Convey("Returns nil if undefined is supplied", t, func() {
		So(newCallback(js.Undefined), ShouldBeNil)
	})

	Convey("Returns a function", t, func() {
		cb := newCallback(&js.Object{})
		var expected callback = func(interface{}) interface{} { return nil }
		So(cb, ShouldHaveSameTypeAs, expected)
	})
}

func TestReflectAll(t *testing.T) {

	Convey("Converts all supplied values to a slice of reflect.Value", t, func() {
		values := reflectAll(1, 2, 3)
		So(len(values), ShouldEqual, 3)
		So(values, ShouldHaveSameTypeAs, []reflect.Value{})
		So(values[0], ShouldHaveSameTypeAs, reflect.ValueOf(1))
		So(values[1], ShouldHaveSameTypeAs, reflect.ValueOf(2))
		So(values[2], ShouldHaveSameTypeAs, reflect.ValueOf(3))
	})
}

func TestUnReflectAll(t *testing.T) {

	Convey("Converts the supplied slice of reflect.Values to a slice of empty interfaces", t, func() {
		values := unReflectAll(reflectAll(1, 2, 3))
		So(len(values), ShouldEqual, 3)
		So(values, ShouldHaveSameTypeAs, []interface{}{})
		So(values[0], ShouldEqual, 1)
		So(values[1], ShouldEqual, 2)
		So(values[2], ShouldEqual, 3)
	})
}

func TestSplitResults(t *testing.T) {

	Convey("When there isn't an error at the end of the function", t, func() {

		Convey("Returns nil if no values were returned", func() {
			result, err := splitResults([]reflect.Value{}, false)
			So(result, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Returns a single value if 1 value was returned", func() {
			result, err := splitResults([]reflect.Value{reflect.ValueOf(1)}, false)
			So(result, ShouldEqual, 1)
			So(err, ShouldBeNil)
		})

		Convey("Returns a slice of values if multiple values were returned", func() {
			result, err := splitResults([]reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)}, false)
			So(result, ShouldHaveSameTypeAs, []interface{}{})
			So(err, ShouldBeNil)

			slice := result.([]interface{})
			So(len(slice), ShouldEqual, 2)
			So(slice[0], ShouldEqual, 1)
			So(slice[1], ShouldEqual, 2)
		})
	})

	Convey("When there is an error at the end of the function", t, func() {

		Convey("Returns nil if no values were returned", func() {
			result, err := splitResults([]reflect.Value{}, true)
			So(result, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Returns the error if there was only 1 error", func() {
			result, err := splitResults([]reflect.Value{reflect.ValueOf(errors.New("nope"))}, true)
			So(result, ShouldBeNil)
			So(err, ShouldResemble, errors.New("nope"))
		})

		Convey("Returns the result and error if there was 1 result", func() {
			result, err := splitResults([]reflect.Value{reflect.ValueOf(1), reflect.ValueOf(errors.New("nope"))}, true)
			So(result, ShouldEqual, 1)
			So(err, ShouldResemble, errors.New("nope"))
		})

		Convey("Returns a slice of results and final error if there were multiple results", func() {
			result, err := splitResults([]reflect.Value{
				reflect.ValueOf(1),
				reflect.ValueOf(2),
				reflect.ValueOf(errors.New("nope")),
			}, true)

			So(result, ShouldHaveSameTypeAs, []interface{}{})
			So(err, ShouldResemble, errors.New("nope"))

			slice := result.([]interface{})
			So(len(slice), ShouldEqual, 2)
			So(slice[0], ShouldEqual, 1)
			So(slice[1], ShouldEqual, 2)
		})
	})
}

func TestHasLastError(t *testing.T) {

	Convey("Returns false if the function returns nothing", t, func() {
		none := reflect.ValueOf(func(){})
		So(hasLastError(none.Type()), ShouldBeFalse)
	})

	Convey("Returns false if the function returns non-error values", t, func() {
		none := reflect.ValueOf(func() int { return 0 })
		So(hasLastError(none.Type()), ShouldBeFalse)
	})

	Convey("Returns true if the function returns only an error", t, func() {
		none := reflect.ValueOf(func() error { return nil })
		So(hasLastError(none.Type()), ShouldBeTrue)
	})

	Convey("Returns true if the function returns values and an error", t, func() {
		none := reflect.ValueOf(func() (int, error) { return 1, nil })
		So(hasLastError(none.Type()), ShouldBeTrue)
	})
}
