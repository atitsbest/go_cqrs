// Package storage kümmert sich um die Speicherung der EventSourcing Events.
package storage

import (
	"fmt"

	"github.com/atitsbest/go_cqrs/eventsourcing"
)

// MemoryStore speichert EventStream in einer Map im Speichern.
type MemoryStore struct {
	items map[eventsourcing.EventSourceId][]eventsourcing.Event
}

// NewMemoryStore erstellt einen neuen EventStore im Speicher.
// Die Indizierung der EventSourceIds passiert über eine Map.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: map[eventsourcing.EventSourceId][]eventsourcing.Event{},
	}
}

// AppendToStream speichert ein oder mehrere Events unter der angegebenen id.
func (mem *MemoryStore) AppendToStream(id eventsourcing.EventSourceId, events []eventsourcing.Event, expectedVersion uint64) error {
	changes, ok := mem.items[id]

	if !ok {
		changes = []eventsourcing.Event{}
	}

	for _, e := range events {
		mem.items[id] = append(changes, e)
	}

	fmt.Printf("%d CHANGES", len(mem.items[id]))
	return nil
}

// LoadEventStream lädt alle Events zu einer EventSourceId (gleichzusetzen mit Aggregate).
func (mem *MemoryStore) LoadEventStream(id eventsourcing.EventSourceId) ([]eventsourcing.Event, uint64, error) {
	changes, ok := mem.items[id]

	if !ok {
		return nil, 0, nil
	}

	return changes, 0, nil
}
