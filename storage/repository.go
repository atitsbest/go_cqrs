package storage

import (
	"fmt"

	"github.com/atitsbest/go_cqrs/eventsourcing"
)

type (
	Repository struct {
		storage EventStorage
	}
)

// CTR
func NewRepository(store EventStorage) *Repository {
	return &Repository{
		storage: store,
	}
}

func (self *Repository) Save(entity eventsourcing.EventSource) error {
	return self.storage.AppendToStream(entity.ID(), entity.UncommitedChanges(), entity.Version())
}

func (self *Repository) Load(id eventsourcing.EventSourceId, source eventsourcing.EventSource) error {
	events, _, err := self.storage.LoadEventStream(id)
	if err != nil {
		return err
	}

	fmt.Printf("%d Events anwenden...", len(events))
	for _, e := range events {
		fmt.Printf("#v", e)
		source.ApplyChange(e)
	}

	return nil
}
