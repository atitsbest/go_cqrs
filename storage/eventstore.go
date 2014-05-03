// Package eventstore provides ...
package storage

import (
	"github.com/atitsbest/go_cqrs/eventsourcing"
)

type EventStorage interface {
	AppendToStream(id eventsourcing.EventSourceId, events []eventsourcing.Event) error
	LoadEventStream(id eventsourcing.EventSourceId) ([]eventsourcing.Event, error)
}
