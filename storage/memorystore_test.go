package storage

import (
	"testing"

	"github.com/atitsbest/go_cqrs/eventsourcing"
	"github.com/atitsbest/go_cqrs/tests/events"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	var (
		sut      *MemoryStore
		id       eventsourcing.EventSourceId
		original []eventsourcing.Event
		loaded   []eventsourcing.Event
		err      error
	)

	Convey("When I create a new MemoryStore", t, func() {
		sut = NewMemoryStore()
		id = eventsourcing.NewEventSourceId()

		Convey("And save an event", func() {
			original = []eventsourcing.Event{events.CartNameChanged{Name: "Cart"}}
			sut.Save(id, original)

			Convey("And load that event", func() {
				loaded, err = sut.LoadEventStream(id)
				So(err, ShouldBeNil)

				Convey("Then I get the saved event", func() {
					So(len(loaded), ShouldEqual, len(original))
					// Geladene Events müsse umgewandelt werden.
					originalEvent := original[0].(events.CartNameChanged)
					loadedEvent := loaded[0].(events.CartNameChanged)
					So(loadedEvent.Name, ShouldEqual, originalEvent.Name)
				})
			})
		})
	})
}
