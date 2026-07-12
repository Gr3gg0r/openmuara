package audit

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openmuara/openmuara/internal/store/migrations"
	_ "modernc.org/sqlite" // register sqlite driver with database/sql
)

// SQLiteStore persists audit events in a SQLite database.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (and migrates) the audit log table at path.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", fmt.Sprintf("%s?_busy_timeout=5000&_journal_mode=WAL", path))
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	store, err := NewSQLiteStoreFromDB(db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

// NewSQLiteStoreFromDB creates an audit SQLiteStore from an existing database connection.
func NewSQLiteStoreFromDB(db *sql.DB) (*SQLiteStore, error) {
	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("migrate audit logs: %w", err)
	}
	return store, nil
}

func (s *SQLiteStore) migrate() error {
	sqlText, err := migrations.Read("003_audit_logs.sql")
	if err != nil {
		return err
	}
	_, err = s.db.Exec(sqlText)
	return err
}

// Close closes the underlying database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// Save inserts an audit event.
func (s *SQLiteStore) Save(event Event) error {
	if event.ID == "" {
		event.ID = uuid.Must(uuid.NewRandom()).String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	_, err := s.db.Exec(`
		INSERT INTO audit_logs (id, timestamp, actor, action, resource_type, resource_id, payload, result)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, event.ID, event.Timestamp.UnixMilli(), event.Actor, event.Action, event.ResourceType, event.ResourceID, event.Payload, event.Result)
	if err != nil {
		return fmt.Errorf("insert audit event: %w", err)
	}
	return nil
}

// List returns the most recent audit events ordered by timestamp descending.
func (s *SQLiteStore) List(limit, offset int) ([]Event, error) {
	return s.ListSince(limit, offset, time.Time{})
}

// ListSince returns audit events at or after since, ordered by timestamp descending.
func (s *SQLiteStore) ListSince(limit, offset int, since time.Time) ([]Event, error) {
	if limit <= 0 {
		limit = 50
	}

	var rows *sql.Rows
	var err error
	if since.IsZero() {
		rows, err = s.db.Query(`
			SELECT id, timestamp, actor, action, resource_type, resource_id, payload, result
			FROM audit_logs
			ORDER BY timestamp DESC, rowid DESC
			LIMIT ? OFFSET ?
		`, limit, offset)
	} else {
		rows, err = s.db.Query(`
			SELECT id, timestamp, actor, action, resource_type, resource_id, payload, result
			FROM audit_logs
			WHERE timestamp >= ?
			ORDER BY timestamp DESC, rowid DESC
			LIMIT ? OFFSET ?
		`, since.UnixMilli(), limit, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("list audit events: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var events []Event
	for rows.Next() {
		var ev Event
		var tsMS int64
		if err := rows.Scan(&ev.ID, &tsMS, &ev.Actor, &ev.Action, &ev.ResourceType, &ev.ResourceID, &ev.Payload, &ev.Result); err != nil {
			return nil, fmt.Errorf("scan audit event: %w", err)
		}
		ev.Timestamp = time.UnixMilli(tsMS)
		events = append(events, ev)
	}
	return events, rows.Err()
}

// Clear removes all audit events from the store.
func (s *SQLiteStore) Clear() error {
	_, err := s.db.Exec(`DELETE FROM audit_logs`)
	if err != nil {
		return fmt.Errorf("clear audit logs: %w", err)
	}
	return nil
}
