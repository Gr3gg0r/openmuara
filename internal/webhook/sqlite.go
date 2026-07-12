package webhook

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite" // register sqlite driver with database/sql
)

// SQLiteStore is a file-backed AttemptStore using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (and migrates) the SQLite webhook store at path.
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

// NewSQLiteStoreFromDB creates a SQLiteStore from an existing database connection.
func NewSQLiteStoreFromDB(db *sql.DB) (*SQLiteStore, error) {
	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("migrate sqlite: %w", err)
	}
	return store, nil
}

// Close closes the underlying database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// Clear removes all webhook attempts from the SQLite store.
func (s *SQLiteStore) Clear() error {
	_, err := s.db.Exec(`DELETE FROM webhook_attempts`)
	if err != nil {
		return fmt.Errorf("clear webhook attempts: %w", err)
	}
	return nil
}

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS webhook_attempts (
			id TEXT PRIMARY KEY,
			ref TEXT NOT NULL UNIQUE,
			provider_name TEXT NOT NULL,
			url TEXT NOT NULL,
			status TEXT NOT NULL,
			attempts INTEGER NOT NULL DEFAULT 0,
			last_error TEXT,
			payload BLOB,
			headers TEXT,
			trace_id TEXT,
			signature_valid INTEGER,
			history TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_webhook_attempts_ref ON webhook_attempts(ref);
	`)
	return err
}

// Save stores or updates an attempt.
func (s *SQLiteStore) Save(a *Attempt) error {
	now := time.Now()
	if a.ID == "" {
		a.ID = uuid.Must(uuid.NewRandom()).String()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	headersJSON, err := json.Marshal(a.Headers)
	if err != nil {
		return fmt.Errorf("marshal headers: %w", err)
	}
	historyJSON, err := json.Marshal(a.History)
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}

	sigValid := sql.NullInt32{}
	if a.SignatureValid != nil {
		sigValid.Valid = true
		if *a.SignatureValid {
			sigValid.Int32 = 1
		}
	}

	_, err = s.db.Exec(`
		INSERT INTO webhook_attempts (id, ref, provider_name, url, status, attempts, last_error, payload, headers, trace_id, signature_valid, history, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(ref) DO UPDATE SET
			provider_name = excluded.provider_name,
			url = excluded.url,
			status = excluded.status,
			attempts = excluded.attempts,
			last_error = excluded.last_error,
			payload = excluded.payload,
			headers = excluded.headers,
			trace_id = excluded.trace_id,
			signature_valid = excluded.signature_valid,
			history = excluded.history,
			updated_at = excluded.updated_at
	`, a.ID, a.Ref, a.ProviderName, a.URL, a.Status, a.Attempts, a.LastError, a.Payload, string(headersJSON), a.TraceID, sigValid, string(historyJSON), a.CreatedAt.UnixMilli(), a.UpdatedAt.UnixMilli())
	if err != nil {
		return fmt.Errorf("insert webhook attempt: %w", err)
	}
	return nil
}

// Get retrieves an attempt by reference.
func (s *SQLiteStore) Get(ref string) (*Attempt, error) {
	row := s.db.QueryRow(`
		SELECT id, ref, provider_name, url, status, attempts, last_error, payload, headers, trace_id, signature_valid, history, created_at, updated_at
		FROM webhook_attempts
		WHERE ref = ?
	`, ref)
	return scanAttempt(row)
}

// List returns the most recent attempts up to limit, skipping offset results.
func (s *SQLiteStore) List(limit, offset int) ([]*Attempt, error) {
	if limit <= 0 {
		limit = -1
	}
	rows, err := s.db.Query(`
		SELECT id, ref, provider_name, url, status, attempts, last_error, payload, headers, trace_id, signature_valid, history, created_at, updated_at
		FROM webhook_attempts
		ORDER BY created_at DESC, rowid DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list webhook attempts: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanAttempts(rows)
}

func scanAttempt(row *sql.Row) (*Attempt, error) {
	var a Attempt
	var headersJSON, historyJSON string
	var createdMS, updatedMS int64
	var sigValid sql.NullInt32

	err := row.Scan(
		&a.ID, &a.Ref, &a.ProviderName, &a.URL, &a.Status, &a.Attempts, &a.LastError, &a.Payload,
		&headersJSON, &a.TraceID, &sigValid, &historyJSON, &createdMS, &updatedMS,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan webhook attempt: %w", err)
	}

	if err := json.Unmarshal([]byte(headersJSON), &a.Headers); err != nil {
		return nil, fmt.Errorf("unmarshal headers: %w", err)
	}
	if err := json.Unmarshal([]byte(historyJSON), &a.History); err != nil {
		return nil, fmt.Errorf("unmarshal history: %w", err)
	}
	if sigValid.Valid {
		v := sigValid.Int32 == 1
		a.SignatureValid = &v
	}
	a.CreatedAt = time.UnixMilli(createdMS)
	a.UpdatedAt = time.UnixMilli(updatedMS)
	return &a, nil
}

func scanAttempts(rows *sql.Rows) ([]*Attempt, error) {
	var result []*Attempt
	for rows.Next() {
		var a Attempt
		var headersJSON, historyJSON string
		var createdMS, updatedMS int64
		var sigValid sql.NullInt32

		err := rows.Scan(
			&a.ID, &a.Ref, &a.ProviderName, &a.URL, &a.Status, &a.Attempts, &a.LastError, &a.Payload,
			&headersJSON, &a.TraceID, &sigValid, &historyJSON, &createdMS, &updatedMS,
		)
		if err != nil {
			return nil, fmt.Errorf("scan webhook attempt: %w", err)
		}

		if err := json.Unmarshal([]byte(headersJSON), &a.Headers); err != nil {
			return nil, fmt.Errorf("unmarshal headers: %w", err)
		}
		if err := json.Unmarshal([]byte(historyJSON), &a.History); err != nil {
			return nil, fmt.Errorf("unmarshal history: %w", err)
		}
		if sigValid.Valid {
			v := sigValid.Int32 == 1
			a.SignatureValid = &v
		}
		a.CreatedAt = time.UnixMilli(createdMS)
		a.UpdatedAt = time.UnixMilli(updatedMS)
		result = append(result, &a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate webhook attempts: %w", err)
	}
	return result, nil
}
