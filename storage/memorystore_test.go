package storage

import (
	"testing"

	"github.com/atitsbest/go_cqrs/eventsourcing"
	"github.com/atitsbest/go_cqrs/tests/events"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMemoryStore(t *testing.T) {
	var (
		sut      *MemoryStore
		id       eventsourcing.EventSourceId
		original []eventsourcing.Event
		loaded   []eventsourcing.Event
		err      error
		version  uint64
	)

	Convey("When I create a new MemoryStore", t, func() {
		sut = NewMemoryStore()
		id = eventsourcing.NewEventSourceId()

		Convey("And save an event", func() {
			original = []eventsourcing.Event{events.CartNameChanged{Name: "Cart"}}
			sut.AppendToStream(id, original, 0)

			Convey("And load that event", func() {
				loaded, version, err = sut.LoadEventStream(id)
				So(err, ShouldBeNil)

				Convey("Then I get the saved event", func() {
					So(len(loaded), ShouldEqual, len(original))
					// Geladene Events müsse umgewandelt werden.
					originalEvent := original[0].(events.CartNameChanged)
					loadedEvent := loaded[0].(events.CartNameChanged)
					So(loadedEvent.Name, ShouldEqual, originalEvent.Name)
				})

				// Convey("With the current version", func() {
				// 	So(version, ShouldEqual, 1)
				// })
			})
		})

		// Convey("And save two events of an aggregate", func() {
		// 	Convey("When I load the aggregate twice", func() {
		// 		Convey("And make changes to both instances", func() {
		// 			Convey("AppendToStream should recognize concurrency problems", func() {
		//
		// 			})
		// 		})
		// 	})
		// })
	})
}
