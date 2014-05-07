// Package eventstore provides ...
package storage

import (
	"github.com/atitsbest/go_cqrs/eventsourcing"
)

// EventStorage kümmert sich darum Events zu speichern und zu laden.
// Dabei wird auch die Gleichläufigkeit berücksichtigt.
type EventStorage interface {
	AppendToStream(id eventsourcing.EventSourceId, events []eventsourcing.Event, expectedVersion uint64) error
	LoadEventStream(id eventsourcing.EventSourceId) ([]eventsourcing.Event, uint64, error)
}
