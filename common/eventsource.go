package common

import (
	"reflect"
	"strings"
)

type EventSource interface {
	Id() EventSourceId
	ApplyChange(e Event)
	UncommitedChanges() []Event
}

type eventSource struct {
	id      EventSourceId
	changes []Event
	source  interface{}
}

// CTR
func NewEventSource(source interface{}) *eventSource {
	es := &eventSource{}
	es.source = source
	es.id = NewEventSourceId()
	return es
}

// Events von einem EventStream lesen.
func CreateFromEventStream(source interface{}, id EventSourceId, es []Event) *eventSource {
	result := NewEventSource(source)
	result.id = id
	for _, e := range es {
		result.handleChange(e)
	}

	return result
}

// UId für diesen EventSoruce
func (self *eventSource) Id() EventSourceId {
	return self.id
}

// Domain-Event anwenden (aber nicht persistieren)
func (self *eventSource) ApplyChange(e Event) {
	// Event-Handler aufrufen
	self.handleChange(e)

	// Event sichern
	self.changes = append(self.changes, e)
}

// Liste mit allen noch nicht gespeicherten Events.
func (self *eventSource) UncommitedChanges() []Event {
	return self.changes
}

// Event vom Entity verarbeiten lassen.
func (self *eventSource) handleChange(e Event) {
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
