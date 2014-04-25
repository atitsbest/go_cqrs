package common

import (
	"reflect"
	"strings"
)

type EventSource struct {
	changes []Event
	source  interface{}
}

// CTR
func NewEventSource(source interface{}) *EventSource {
	es := new(EventSource)
	es.source = source
	return es
}

// Domain-Event anwenden (aber nicht persistieren)
func (self *EventSource) ApplyChange(e Event) {
	// Event-Handler aufrufen
	self.handleChange(e)

	// Event sichern
	self.changes = append(self.changes, e)
}

// Liste mit allen noch nicht gespeicherten Events.
func (self *EventSource) UncommitedChanges() []Event {
	return self.changes
}

// Events von einem EventStream lesen.
func (self *EventSource) CreateFromEventStream(es []Event) {
	for _, e := range es {
		self.handleChange(e)
	}
}

// Event vom Entity verarbeiten lassen.
func (self *EventSource) handleChange(e Event) {
	sourceType := reflect.TypeOf(self.source)
	mc := sourceType.NumMethod()
	for i := 0; i < mc; i += 1 {
		method := sourceType.Method(i)

		if strings.HasPrefix(method.Name, "Handle") {
			if method.Type.NumIn() == 2 {
				eventType := method.Type.In(1)
				if eventType == reflect.TypeOf(e) {
					sourceValue := reflect.ValueOf(self.source)
					eventValue := reflect.ValueOf(e)
					method.Func.Call([]reflect.Value{sourceValue, eventValue})
				}
			}
		}
	}
}
