package entities

import (
	"github.com/atitsbest/go_cqrs/cart/events"
	"github.com/atitsbest/go_cqrs/common"
)

type Cart struct {
	// Muss im CTR initialisiert werden.
	common.EventSource

	name string
}

// CTR
func NewCart() *Cart {
	cart := new(Cart)
	cart.EventSource = common.NewEventSource(cart)

	cart.EventSource.ApplyChange(events.CartCreated{})

	return cart
}

// CTR
func NewCartFromEventStream(es []common.Event) *Cart {
	cart := new(Cart)
	cart.EventSource = common.CreateFromEventStream(cart, es)

	return cart
}

// Name
func (self *Cart) Name() string {
	return self.name
}

// ------------------------- EVENTS -------------------------

func (self *Cart) HandleCartNameChanged(e events.CartNameChanged) {
	self.name = e.Name
}
