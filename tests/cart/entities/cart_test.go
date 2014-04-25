package test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	. "github.com/atitsbest/go_cqrs/cart/entities"
	"github.com/atitsbest/go_cqrs/cart/events"
	"github.com/atitsbest/go_cqrs/common"
)

func TestSpec(t *testing.T) {
	var sut *Cart
	var stream []common.Event

	Convey("When I create a new cart", t, func() {
		sut = NewCart()

		Convey("Then the change is recorded as an event", func() {
			So(len(sut.EventSource.UncommitedChanges()), ShouldEqual, 1)
		})
	})

	Convey("Given an eventstream", t, func() {
		stream = []common.Event{
			events.CartCreated{},
			events.CartNameChanged{Name: "Yannick"},
		}
		Convey("When I create a cart from the eventstream", func() {
			sut = NewCartFromEventStream(stream)

			Convey("Then I get a fully restored cart", func() {
				So(sut.Name(), ShouldEqual, "Yannick")
				So(len(sut.EventSource.UncommitedChanges()), ShouldEqual, 0)
			})
		})
	})
}
