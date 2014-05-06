package storage

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEventRegistration(t *testing.T) {
	type (
		Event1 struct{ Name string }
		Event2 struct{ Count int64 }
	)
	var (
		sut *EventRegistration
		err error
	)

	Convey("Given an empty EventRegistration", t, func() {
		sut = NewEventRegistration()
		So(sut, ShouldNotBeNil)

		Convey("When I register an event", func() {
			err = sut.Register(Event1{})
			So(err, ShouldBeNil)

			Convey("Then I can get the Type from the event name", func() {
				t, err := sut.Get("Event1")
				So(err, ShouldBeNil)
				So(t, ShouldNotBeNil)
				So(t.Name(), ShouldEqual, "Event1")
			})

			Convey("And I register the same event again", func() {
				err = sut.Register(Event1{})
				Convey("Then I get an error", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Convey("And I try to get a not registered event", func() {
				_, err = sut.Get("NichtDa")
				Convey("Then I get an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})

	})
}
