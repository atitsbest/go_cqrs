package domain

import "github.com/atitsbest/go_cqrs/eventsourcing"
import "github.com/atitsbest/go_cqrs/tests/events"

type Cart struct {
	// Muss im CTR initialisiert werden.
	eventsourcing.EventSource

	name string
}

// CTR
func NewCart() *Cart {
	cart := new(Cart)
	cart.EventSource = eventsourcing.NewEventSource(cart)

	cart.ApplyChange(events.CartCreated{})

	return cart
}

// CTR
func CreateCartFromEventStream(id eventsourcing.EventSourceId, es []eventsourcing.Event) *Cart {
	cart := new(Cart)
	cart.EventSource = eventsourcing.CreateFromEventStream(cart, id, es)

	return cart
}

// Name
func (self *Cart) Name() string {
	return self.name
}

// Name ändern.
func (self *Cart) SetName(name string) {
	self.ApplyChange(events.CartNameChanged{Name: name})
}

// ------------------------- EVENTS -------------------------

func (self *Cart) HandleCartNameChanged(e events.CartNameChanged) {
	self.name = e.Name
}
