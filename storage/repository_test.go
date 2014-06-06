package storage

import (
	"fmt"
	"testing"

	"github.com/atitsbest/go_cqrs/tests/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRepository(t *testing.T) {
	var (
		sut    *Repository
		store  *MemoryStore
		cart   *domain.Cart
		loaded *domain.Cart
		err    error
	)

	Convey("Given a repository", t, func() {
		store = NewMemoryStore()
		sut = NewRepository(store)

		Convey("When I save an entity", func() {
			cart = domain.NewCart()
			cart.SetName("TestCart")
			fmt.Printf("%d EVENTS", len(cart.UncommitedChanges()))
			err = sut.Save(cart)

			Convey("Then the unchanged events got persisted", func() {
			})
			Convey("And the domain events got dispatched", func() {
			})

			Convey("When I load the entity", func() {

				// Damit ein Entity/EventSource geladen werden kann, muss er zuerst
				// instanziert werden.
				// Ein einfaches Instanzieren wie &domain.Cart{} functioniert nicht,
				// weil dabei der EventSource-Teil des Entities (siehe Cart-CTR) nicht
				// initialisiert wird und wir beim ApplyChange einen Fehler bekommen.
				// Dafür brauchen wir noch eine Lösung!
				loaded = domain.NewCart()
				err = sut.Load(cart.ID(), loaded)

				So(loaded, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Then I get the last saved state of the entity", func() {
					So(loaded.Name(), ShouldEqual, "TestCart")
					So(loaded.Version(), ShouldEqual, cart.Version())
				})
			})
		})
	})
}

func must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}
