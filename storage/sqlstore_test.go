package storage

import (
	"database/sql"
	"os"
	"testing"

	es "github.com/atitsbest/go_cqrs/eventsourcing"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSqlStore(t *testing.T) {
	type (
		Event1 struct{ Name string }
		Event2 struct{ Count int64 }
	)
	var (
		sut          *SqlStore
		id           es.EventSourceId
		e1           *Event1
		e2           *Event2
		loadedEvents []es.Event
		err          error
		db           *sql.DB
		er           *EventRegistration
	)

	// DB init.
	dbName := "./sqlstore_test.db"
	os.Remove(dbName)
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	if db.Ping() != nil {
		panic(err)
	}

	// Events registrieren.
	er = NewEventRegistration()
	er.Register(Event1{})
	er.Register(Event2{})

	Convey("Given a SqlStore", t, func() {
		sut, err = NewSqlStore(db, er)
		So(sut, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("And two different events", func() {
			e1 = &Event1{Name: "Neu"}
			e2 = &Event2{Count: 17}

			Convey("When I append them to the EventStream", func() {
				err = sut.AppendToStream(id, []es.Event{e1, e2})
				So(err, ShouldBeNil)

				Convey("And load them from the EventStream", func() {
					loadedEvents, err = sut.LoadEventStream(id)
					So(err, ShouldBeNil)

					Convey("Then all events should be loaded", func() {
						So(len(loadedEvents), ShouldEqual, 2)

						// Inhalt der geladenen Events überprüfen.
						So(loadedEvents[0].(*Event1).Name, ShouldEqual, "Neu")
						So(loadedEvents[1].(*Event2).Count, ShouldEqual, 17)
					})
				})
			})
		})
	})

}
