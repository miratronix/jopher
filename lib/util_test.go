package lib

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"errors"
	"reflect"
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

func TestCallReflected(t *testing.T) {

	Convey("Calls the function with one supplied argument", t, func() {
		reflected := ReflectFunction(func(input *int) { *input++ })

		input := 1
		CallReflected(reflected, &input)
		So(input, ShouldEqual, 2)
	})

	Convey("Calls the function with two supplied arguments", t, func() {
		reflected := ReflectFunction(func(input1 *int, input2 *int) {
			*input1++
			*input2++
		})

		input1 := 1
		input2 := 2
		CallReflected(reflected, &input1, &input2)
		So(input1, ShouldEqual, 2)
		So(input2, ShouldEqual, 3)
	})

	Convey("Returns nil when the function doesn't return anything", t, func() {
		result, err := CallReflected(ReflectFunction(func() {}))
		So(result, ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Returns the value when the function returns a value", t, func() {
		result, err := CallReflected(ReflectFunction(func() int { return 1 }))
		So(result, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})

	Convey("Returns a slice when the function returns multiple values", t, func() {
		result, err := CallReflected(ReflectFunction(func() (int, int) { return 1, 2 }))
		So(result, ShouldResemble, []interface{}{1, 2})
		So(err, ShouldBeNil)
	})

	Convey("Returns a nil error when the function returns a nil error", t, func() {
		result, err := CallReflected(ReflectFunction(func() error { return nil }))
		So(result, ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Returns the error when the function returns a error", t, func() {
		result, err := CallReflected(ReflectFunction(func() error { return errors.New("nope") }))
		So(result, ShouldBeNil)
		So(err, ShouldResemble, errors.New("nope"))
	})

	Convey("Returns the error if the function returns an error as the last result", t, func() {
		result, err := CallReflected(ReflectFunction(func() (int, error) { return 3, errors.New("nope") }))
		So(result, ShouldBeNil)
		So(err, ShouldResemble, errors.New("nope"))
	})
}

func TestCallOnPanic(t *testing.T) {

	Convey("Calls the supplied function with the panic value if there is a panic", t, func() {
		ch := make(chan interface{})
		handleError := func(val interface{}) {
			ch <- val
		}

		go func(){
			defer CallOnPanic(handleError)
			panic("nope")
		}()

		So(<-ch, ShouldEqual, "nope")
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
		none := reflect.ValueOf(func() {})
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
