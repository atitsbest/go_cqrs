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

	cart.ApplyChange(events.CartCreated{})

	return cart
}

// CTR
func CreateCartFromEventStream(id common.EventSourceId, es []common.Event) *Cart {
	cart := new(Cart)
	cart.EventSource = common.CreateFromEventStream(cart, id, es)

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
