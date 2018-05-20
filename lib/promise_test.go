package lib

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResolve(t *testing.T) {

	Convey("When no callback is attached", t, func() {

		Convey("Sets the state to resolved and saves the value", func() {
			p := &Promise{}
			p.Resolve("hello")
			So(p.state, ShouldEqual, resolved)
			So(p.value, ShouldEqual, "hello")
		})
	})

	Convey("When a callback is attached", t, func() {

		Convey("Calls the callback with the resolved value", func() {
			p := &Promise{}

			p.resolved = func(val interface{}) interface{} {
				So(val, ShouldEqual, "hello")
				return nil
			}
			p.Resolve("hello")
		})

		Convey("Sets the sate to resolved and saves the result if the callback didn't panic", func() {
			p := &Promise{}

			p.resolved = func(val interface{}) interface{} { return "hello" }
			p.Resolve(nil)

			So(p.state, ShouldEqual, resolved)
			So(p.value, ShouldEqual, "hello")
		})

		Convey("Sets the state to rejected and saves the result if the callback panicked", func() {
			p := &Promise{}

			p.resolved = func(val interface{}) interface{} { panic("hello") }
			p.Resolve(nil)

			So(p.value, ShouldEqual, "hello")
			So(p.state, ShouldEqual, rejected)
		})
	})
}

func TestReject(t *testing.T) {

	Convey("When no callback is attached", t, func() {

		Convey("Sets the state to rejected and saves the value", func() {
			p := &Promise{}
			p.Reject("hello")
			So(p.state, ShouldEqual, rejected)
			So(p.value, ShouldEqual, "hello")
		})
	})

	Convey("When a callback is attached", t, func() {

		Convey("Calls the callback with the rejected value", func() {
			p := &Promise{}

			p.rejected = func(val interface{}) interface{} {
				So(val, ShouldEqual, "hello")
				return nil
			}
			p.Reject("hello")
		})

		Convey("Sets the state to resolved and saves the result if the callback didn't panic", func() {
			p := &Promise{}

			p.rejected = func(val interface{}) interface{} { return "hello" }
			p.Reject(nil)

			So(p.state, ShouldEqual, resolved)
			So(p.value, ShouldEqual, "hello")
		})

		Convey("Sets the state to rejected and saves the result if the callback panicked", func() {
			p := &Promise{}

			p.rejected = func(val interface{}) interface{} { panic("hello") }
			p.Reject(nil)

			So(p.value, ShouldEqual, "hello")
			So(p.state, ShouldEqual, rejected)
		})
	})
}

func TestThen(t *testing.T) {

	Convey("Saves the supplied resolve/reject callbacks", t, func() {
		p := &Promise{}

		resolve := func(val interface{}) interface{} { return nil }
		reject := func(val interface{}) interface{} { return nil }

		p.Then(resolve, reject)
		So(p.resolved, ShouldEqual, resolve)
		So(p.rejected, ShouldEqual, reject)
	})

	Convey("Creates, stores, and returns a new child promise", t, func() {
		p := &Promise{}
		c := p.Then(nil, nil)

		So(c, ShouldHaveSameTypeAs, p)
		So(p.child, ShouldEqual, c)
	})

	Convey("When the promise is pending", t, func() {

		Convey("Returns a pending child promise", func() {
			p := &Promise{}
			c := p.Then(nil, nil)
			So(c.state, ShouldEqual, pending)
		})
	})

	Convey("When the promise is already resolved", t, func() {

		Convey("Calls the resolve callback", func() {
			p := &Promise{state: resolved, value: "hello"}

			resolve := func(val interface{}) interface{} {
				So(val, ShouldEqual, "hello")
				return "goodbye"
			}
			p.Then(resolve, nil)

			So(p.state, ShouldEqual, resolved)
			So(p.value, ShouldEqual, "goodbye")
		})

		Convey("When the resolve callback doesn't panic", func() {

			Convey("Sets the child state to resolved and stores the value", func() {
				p := &Promise{state: resolved, value: "hello"}
				c := p.Then(func(val interface{}) interface{} { return "hello" }, nil)

				So(c.state, ShouldEqual, resolved)
				So(c.value, ShouldEqual, "hello")
			})
		})

		Convey("When the resolve callback panics", func() {

			Convey("Sets the child state to rejected and stores the value", func() {
				p := &Promise{state: resolved, value: "hello"}
				c := p.Then(func(val interface{}) interface{} { panic("nope") }, nil)

				So(c.state, ShouldEqual, rejected)
				So(c.value, ShouldEqual, "nope")
			})
		})
	})

	Convey("When the promise is already rejected", t, func() {

		Convey("Calls the reject callback", func() {
			p := &Promise{state: rejected, value: "hello"}

			reject := func(val interface{}) interface{} {
				So(val, ShouldEqual, "hello")
				return "goodbye"
			}
			p.Then(nil, reject)

			So(p.state, ShouldEqual, resolved)
			So(p.value, ShouldEqual, "goodbye")
		})

		Convey("When the reject callback doesn't panic", func() {

			Convey("Sets the child state to resolved and stores the value", func() {
				p := &Promise{state: rejected, value: "hello"}
				c := p.Then(nil, func(val interface{}) interface{} { return "hello" })

				So(c.state, ShouldEqual, resolved)
				So(c.value, ShouldEqual, "hello")
			})
		})

		Convey("When the resolve callback panics", func() {

			Convey("Sets the child state to rejected and stores the value", func() {
				p := &Promise{state: rejected, value: "hello"}
				c := p.Then(nil, func(val interface{}) interface{} { panic("nope") })

				So(c.state, ShouldEqual, rejected)
				So(c.value, ShouldEqual, "nope")
			})
		})
	})
}

func TestCatch(t *testing.T) {

	Convey("Saves the supplied reject callback", t, func() {
		p := &Promise{}

		reject := func(val interface{}) interface{} { return nil }
		p.Catch(reject)

		So(p.rejected, ShouldEqual, reject)
	})
}
