package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"

	es "github.com/atitsbest/go_cqrs/eventsourcing"
)

// ErrConcurrency gibt an, dass die aktuelle Version eines EventSources nicht mit
// der aus dem zu speicherden übereinstimmt d.h. der zu speichernde ist nicht mehr
// aktuell.
var ErrConcurrency = errors.New("EventStore Concurrency Fehler!")

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
		version int not null,
		data blob);`
	_, err := db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AppendToStream speichert Events zu einem EvenSource/AggregateRoot.
func (store *SqlStore) AppendToStream(eventsourceID es.EventSourceId, events []es.Event, expectedVersion uint64) error {
	// Transaction starten.
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	// Aktuelle Version für diesen Evensource ermitteln.
	currentVersion, err := getCurrentEventSourceVersion(tx, eventsourceID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Sind wir zu spät dran?
	if currentVersion != expectedVersion {
		return ErrConcurrency
	}

	for i, e := range events {
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

		sql := "insert into events (id, eventsource_id, type, version, data) values(?, ?, ?, ?, ?)"

		newVersion := uint64(i+1) + expectedVersion

		// Ab in die DB.
		_, err = tx.Exec(sql, eventID.String(), eventsourceID.String(), eventType, newVersion, se)
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
func (store *SqlStore) LoadEventStream(id es.EventSourceId) ([]es.Event, uint64, error) {
	result := []es.Event{} // Event ist ein Interface, also brauchen wir keinen Pointer.
	currentVersion := uint64(0)

	rows, err := store.db.Query("select type, version, data from events where eventsource_id = ?", id.String())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		data := []byte{}
		eventType := ""
		var version uint64

		err = rows.Scan(&eventType, &version, &data)
		if err != nil {
			return nil, 0, err
		}

		// Event-Typ ermitteln.
		t, err := store.events.Get(eventType)
		if err != nil {
			return nil, 0, err
		}
		// Event-Instanz erstellen.
		eventValue := reflect.New(t)
		event := eventValue.Interface()

		// Aus JSON wieder ein Event machen.
		err = json.Unmarshal(data, event)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, event)
		currentVersion = version
	}
	return result, currentVersion, nil
}

// getCurrentEventSourceVersion liefert die aktuelle Version aus der DB für den
// angegebenen EventSource.
func getCurrentEventSourceVersion(db *sql.Tx, eventsourceID es.EventSourceId) (uint64, error) {
	var version uint64
	row := db.QueryRow("select version from events where eventsource_id = ?", eventsourceID.String())
	err := row.Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return version, err
}
