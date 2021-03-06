// Sequenz für das Laden/Ändern/Speichern eines Entities
//
// AL->Repository: GetEntity()
// Repository-->AL: entity
// AL->Entity: ChangeName()
// Entity->Entity: Validieren
// Entity->EventSource: ApplyChange()
// EventSource->Entity: HandleNameChanged()
// AL->Repository: Save(entity)
// Repository->Entity: UncommitedChanges()
// Entity-->Repository: []Events
// Repository->EventStore: Persist([]Events)
//
// (graphische Darstellung unter: http://bramp.github.io/js-sequence-diagrams/)
package eventsourcing

import (
	"reflect"
	"strings"
)

// EventSource is die Grundlage für alle Entities/AggregateRoots die beim EventSourcing
// teilnehmen wollen.
// Das Interface wir am besten per "Vererbung" verwendet.
type EventSource interface {
	ID() EventSourceId
	Version() uint64
	ApplyChange(e Event)
	UncommitedChanges() []Event
}

type eventSource struct {
	id      EventSourceId
	version uint64
	changes []Event
	source  interface{}
}

// NewEventSource erstellt einen EventSource für den
// Typen source. Der Type ist wichtig, damit die Events
// richtig geroutet werden können (Handle...)
func NewEventSource(source interface{}) *eventSource {
	es := &eventSource{}
	es.source = source
	es.id = NewEventSourceId()
	return es
}

// CreateFromEventStream stellt ein Entity aus einem EventStream wieer her.
// Weitere Infos siehe: NewEventSource
func CreateFromEventStream(source interface{}, id EventSourceId, es []Event) *eventSource {
	result := NewEventSource(source)
	result.id = id
	for _, e := range es {
		result.handleChange(e)
	}

	return result
}

// UId für diesen EventSoruce
func (es *eventSource) ID() EventSourceId {
	return es.id
}

// Version des EventSource (wird für Concurrency benötigt).
func (es *eventSource) Version() uint64 {
	return es.version
}

// Domain-Event anwenden (aber nicht persistieren)
func (es *eventSource) ApplyChange(e Event) {
	// Event-Handler aufrufen
	es.handleChange(e)

	// Event sichern
	es.changes = append(es.changes, e)
}

// Liste mit allen noch nicht gespeicherten Events.
func (es *eventSource) UncommitedChanges() []Event {
	return es.changes
}

// Event vom Entity verarbeiten lassen.
func (es *eventSource) handleChange(e Event) {
	sourceType := reflect.TypeOf(es.source)
	mc := sourceType.NumMethod()
	for i := 0; i < mc; i++ {
		method := sourceType.Method(i)

		if strings.HasPrefix(method.Name, "Handle") {
			if method.Type.NumIn() == 2 {
				eventType := method.Type.In(1)
				if eventType == reflect.TypeOf(e) {
					sourceValue := reflect.ValueOf(es.source)
					eventValue := reflect.ValueOf(e)
					method.Func.Call([]reflect.Value{sourceValue, eventValue})
				}
			}
		}
	}
}
