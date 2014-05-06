package storage

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/atitsbest/go_cqrs/eventsourcing"
)

type EventRegistration struct {
	events map[string]reflect.Type
}

// EventRegistration CTR
func NewEventRegistration() *EventRegistration {
	return &EventRegistration{
		events: map[string]reflect.Type{},
	}
}

// Register erfasst ein Event. Dabei wird per Reflection der
// Name des Events, unabhängig vom Package verwendet.
// Ist der Name des Events bereits erfasst, wird ein Fehler
// zurückgegeben.
func (reg *EventRegistration) Register(e eventsourcing.Event) error {
	// Typ ermitteln.
	t := reflect.TypeOf(e)
	// Bereits registriert?
	_, ok := reg.events[t.Name()]
	if ok {
		return errors.New(fmt.Sprintf("Event %s, ist bereits erfasst!", t.Name()))
	}

	// Event erfassen.
	reg.events[t.Name()] = t

	return nil
}

// Get liefert ein erfasstes Event per Name.
// Ist das Event nicht erfasst, wird ein Fehler zurück gegeben.
func (reg *EventRegistration) Get(name string) (reflect.Type, error) {
	t, ok := reg.events[name]
	if ok {
		return t, nil
	}
	return nil, errors.New(fmt.Sprintf("Das Event %s ist nicht erfasst!", name))
}
