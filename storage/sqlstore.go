package storage

import (
	"database/sql"
	"encoding/json"
	"reflect"

	es "github.com/atitsbest/go_cqrs/eventsourcing"
)

type SqlStore struct {
	db     *sql.DB
	events *EventRegistration
}

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

func (store *SqlStore) AppendToStream(eventsourceId es.EventSourceId, events []es.Event) error {
	for _, e := range events {
		// Event serialisieren.
		se, err := json.Marshal(e)
		if err != nil {
			return err
		}

		eventType := reflect.TypeOf(e).Name()
		// Wurde das Event registriert?
		if _, err := store.events.Get(eventType); err != nil {
			return err
		}

		// TODO: Hier brauchen wir eine EventId
		eventId := es.NewEventSourceId()

		sql := "insert into events (id, eventsource_id, type, data) values(?, ?, ?, ?)"

		_, err = store.db.Exec(sql, eventId.String(), eventsourceId.String(), eventType, se)
		if err != nil {
			return err
		}
	}
	return nil
}

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
		loadedEvent := reflect.New(t)

		err = json.Unmarshal(data, loadedEvent)
		if err != nil {
			return nil, err
		}

		result = append(result, loadedEvent)
	}
	return result, nil
}
