// Package eventstore provides ...
package storage

import (
	"github.com/atitsbest/go_cqrs/eventsourcing"
)

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

// Save speichert ein oder mehrere Events unter der angegebenen id.
func (mem *MemoryStore) Save(id eventsourcing.EventSourceId, events []eventsourcing.Event) error {
	changes, ok := mem.items[id]

	if !ok {
		changes = []eventsourcing.Event{}
	}

	for _, e := range events {
		mem.items[id] = append(changes, e)
	}
	return nil
}

func (mem *MemoryStore) LoadEventStream(id eventsourcing.EventSourceId) ([]eventsourcing.Event, error) {
	changes, ok := mem.items[id]

	if !ok {
		return nil, nil
	}

	return changes, nil
}
