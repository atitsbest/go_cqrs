package test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	. "github.com/atitsbest/go_cqrs/cart/entities"
	"github.com/atitsbest/go_cqrs/common"
)

func TestSpec(t *testing.T) {
	var (
		sut     *Cart
		sut2    *Cart
		stream  []common.Event
		stream2 []common.Event
		id      common.EventSourceId
		id2     common.EventSourceId
	)

	Convey("When I create a new cart", t, func() {
		sut = NewCart()

		Convey("Then it has an uid", func() {
			id := sut.Id()
			So(sut.Id(), ShouldNotBeNil)
			So(sut.Id(), ShouldEqual, id) // Sicherstellen, dass auch immer dieselbe Id zurück kommt.

			Convey("When I create another cart", func() {
				sut2 = NewCart()

				Convey("Then it has a different uid than the pervious cart", func() {
					So(sut2.Id(), ShouldNotEqual, sut.Id())
				})
			})
		})

		Convey("Then the change is recorded as an event", func() {
			So(len(sut.EventSource.UncommitedChanges()), ShouldEqual, 1)
		})
	})

	Convey("Given an eventstream", t, func() {
		tmp := NewCart()
		id = tmp.Id()
		tmp.SetName("Yannick")
		stream = tmp.UncommitedChanges()

		Convey("When I create a cart from the eventstream", func() {
			sut = CreateCartFromEventStream(id, stream)

			Convey("Then I get a fully restored cart", func() {
				So(sut.Name(), ShouldEqual, "Yannick")
				So(len(sut.EventSource.UncommitedChanges()), ShouldEqual, 0)
			})

			Convey("With the same uid as before", func() {
				So(id, ShouldEqual, sut.Id())
			})
		})

		Convey("And an eventstream from another enitity", func() {
			tmp = NewCart()
			id2 = tmp.Id()
			tmp.SetName("Ederer")
			tmp.SetName("Meißner")
			stream2 = tmp.UncommitedChanges()

			Convey("When I create both carts from each eventstreams", func() {
				sut = CreateCartFromEventStream(id, stream)
				sut2 = CreateCartFromEventStream(id2, stream2)

				Convey("Then I get two differen carts", func() {
					So(sut, ShouldNotEqual, sut2)
					So(sut.Name(), ShouldEqual, "Yannick")
					So(sut2.Name(), ShouldEqual, "Meißner")
					So(sut.Id(), ShouldNotEqual, sut2.Id())
				})
			})
		})
	})

}
