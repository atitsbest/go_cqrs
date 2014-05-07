package storage

import (
	"database/sql"
	"encoding/json"
	"reflect"

	es "github.com/atitsbest/go_cqrs/eventsourcing"
)

// SqlStore speichert Events in einer Sql-DB.
type SqlStore struct {
	db     *sql.DB
	events *EventRegistration
}

// NewSqlStore ist der CTR für einen
func NewSqlStore(db *sql.DB, reg *EventRegistration) (*SqlStore, error) {
	result := &SqlStore{
		db:     db,
		events: reg,
	}

	// Tabelle erstellen, falls noch nicht vorhanden.
	schema := `create table if not exists events (
		id string not null primary key, 
		eventsource_id string not null, 
		type string not null, 
		data blob);`
	_, err := db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AppendToStream speichert Events zu einem EvenSource/AggregateRoot.
func (store *SqlStore) AppendToStream(eventsourceID es.EventSourceId, events []es.Event) error {
	// Transaction starten.
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	for _, e := range events {
		// Event serialisieren.
		se, err := json.Marshal(e)
		if err != nil {
			return err
		}

		// Achtung: e ist ein Ptr => .Elem() verwenden, sonst liefert Name()
		// 			immer leer.
		eventType := reflect.TypeOf(e).Elem().Name()
		// Wurde das Event registriert?
		if _, err := store.events.Get(eventType); err != nil {
			return err
		}

		// TODO: Hier brauchen wir eine EventId
		eventID := es.NewEventSourceId()

		sql := "insert into events (id, eventsource_id, type, data) values(?, ?, ?, ?)"

		_, err = tx.Exec(sql, eventID.String(), eventsourceID.String(), eventType, se)
		if err != nil {
			tx.Rollback() // TODO: Error von Rollback wird ignoriert. Korrekt?
			return err
		}
	}

	// Commit.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// LoadEventStream läd alle Events zu einem EventSource/AggregateRoot.
func (store *SqlStore) LoadEventStream(id es.EventSourceId) ([]es.Event, error) {
	result := []es.Event{} // Event ist ein Interface, also brauchen wir keinen Pointer.

	rows, err := store.db.Query("select type, data from events where eventsource_id = ?", id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		data := []byte{}
		eventType := ""
		err = rows.Scan(&eventType, &data)
		if err != nil {
			return nil, err
		}

		// Event-Typ ermitteln.
		t, err := store.events.Get(eventType)
		if err != nil {
			return nil, err
		}
		// Event-Instanz erstellen.
		eventValue := reflect.New(t)
		event := eventValue.Interface()

		// Aus JSON wieder ein Event machen.
		err = json.Unmarshal(data, event)
		if err != nil {
			return nil, err
		}

		result = append(result, event)
	}
	return result, nil
}
