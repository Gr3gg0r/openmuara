// Package engine provides the transaction ledger for OpenMuara.
package engine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite" // register sqlite driver with database/sql
)

// SQLiteStore is a file-backed TransactionStore using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (and migrates) the SQLite ledger at path.
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

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			provider TEXT NOT NULL,
			type TEXT NOT NULL,
			amount REAL NOT NULL,
			currency TEXT NOT NULL,
			status TEXT NOT NULL,
			customer_ref TEXT NOT NULL,
			idempotency_key TEXT,
			reference TEXT NOT NULL UNIQUE,
			trace_id TEXT,
			items TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_transactions_reference ON transactions(reference);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_idempotency ON transactions(idempotency_key) WHERE idempotency_key <> '';
	`)
	if err != nil {
		return err
	}
	return s.addTraceIDColumn()
}

func (s *SQLiteStore) addTraceIDColumn() error {
	var count int
	if err := s.db.QueryRow(`
		SELECT COUNT(*) FROM pragma_table_info('transactions') WHERE name = 'trace_id'
	`).Scan(&count); err != nil {
		return fmt.Errorf("check trace_id column: %w", err)
	}
	if count == 0 {
		if _, err := s.db.Exec(`ALTER TABLE transactions ADD COLUMN trace_id TEXT`); err != nil {
			return fmt.Errorf("add trace_id column: %w", err)
		}
	}
	return nil
}

// CreateOrGet stores tx and returns the stored transaction. See TransactionStore.
func (s *SQLiteStore) CreateOrGet(tx Transaction) (Transaction, bool, error) {
	if tx.IdempotencyKey != "" {
		existing, ok, err := s.getByIdempotencyKey(tx.IdempotencyKey)
		if err != nil {
			return Transaction{}, false, err
		}
		if ok {
			return existing, false, nil
		}
	}

	now := time.Now()
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = now
	}
	if tx.UpdatedAt.IsZero() {
		tx.UpdatedAt = now
	}

	// Preserve the existing ID when a reference is reused so that callers
	// can keep referencing the same transaction across status updates.
	referenceExisted := false
	existing, ok, err := s.GetByReference(tx.Reference)
	if err != nil {
		return Transaction{}, false, err
	}
	if ok {
		referenceExisted = true
		tx.ID = existing.ID
		tx.CreatedAt = existing.CreatedAt
	} else if tx.ID == "" {
		tx.ID = uuid.Must(uuid.NewRandom()).String()
	}

	itemsJSON, err := json.Marshal(tx.Items)
	if err != nil {
		return Transaction{}, false, fmt.Errorf("marshal items: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO transactions (id, provider, type, amount, currency, status, customer_ref, idempotency_key, reference, trace_id, items, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(reference) DO UPDATE SET
			provider = excluded.provider,
			type = excluded.type,
			amount = excluded.amount,
			currency = excluded.currency,
			status = excluded.status,
			customer_ref = excluded.customer_ref,
			idempotency_key = excluded.idempotency_key,
			trace_id = COALESCE(excluded.trace_id, trace_id),
			items = excluded.items,
			updated_at = excluded.updated_at
	`, tx.ID, tx.Provider, tx.Type, tx.Amount, tx.Currency, tx.Status, tx.CustomerRef, tx.IdempotencyKey, tx.Reference, tx.TraceID, string(itemsJSON), tx.CreatedAt.UnixMilli(), tx.UpdatedAt.UnixMilli())
	if err != nil {
		return Transaction{}, false, fmt.Errorf("insert transaction: %w", err)
	}

	stored, ok, err := s.GetByID(tx.ID)
	if err != nil {
		return Transaction{}, false, err
	}
	if !ok {
		return Transaction{}, false, fmt.Errorf("transaction not found after insert")
	}

	recordTransaction(stored.Provider, string(stored.Status))
	return stored, !referenceExisted, nil
}

// GetByID returns a transaction by ID.
func (s *SQLiteStore) GetByID(id string) (Transaction, bool, error) {
	row := s.db.QueryRow(`
		SELECT id, provider, type, amount, currency, status, customer_ref, idempotency_key, reference, trace_id, items, created_at, updated_at
		FROM transactions
		WHERE id = ?
	`, id)
	return scanTransaction(row)
}

// Clear removes all transactions from the store.
func (s *SQLiteStore) Clear() error {
	_, err := s.db.Exec(`DELETE FROM transactions`)
	if err != nil {
		return fmt.Errorf("clear transactions: %w", err)
	}
	return nil
}

// GetByReference returns a transaction by reference.
func (s *SQLiteStore) GetByReference(ref string) (Transaction, bool, error) {
	row := s.db.QueryRow(`
		SELECT id, provider, type, amount, currency, status, customer_ref, idempotency_key, reference, trace_id, items, created_at, updated_at
		FROM transactions
		WHERE reference = ?
	`, ref)
	return scanTransaction(row)
}

func (s *SQLiteStore) getByIdempotencyKey(key string) (Transaction, bool, error) {
	row := s.db.QueryRow(`
		SELECT id, provider, type, amount, currency, status, customer_ref, idempotency_key, reference, trace_id, items, created_at, updated_at
		FROM transactions
		WHERE idempotency_key = ?
	`, key)
	return scanTransaction(row)
}

// List returns the most recent transactions up to limit, skipping offset results.
// A limit of zero or less returns all remaining items.
func (s *SQLiteStore) List(limit, offset int) ([]Transaction, error) {
	if limit <= 0 {
		limit = -1
	}
	rows, err := s.db.Query(`
		SELECT id, provider, type, amount, currency, status, customer_ref, idempotency_key, reference, trace_id, items, created_at, updated_at
		FROM transactions
		ORDER BY created_at DESC, rowid DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	return scanTransactions(rows)
}

func scanTransaction(row *sql.Row) (Transaction, bool, error) {
	var tx Transaction
	var itemsJSON string
	var createdMS, updatedMS int64
	var traceID sql.NullString

	err := row.Scan(
		&tx.ID, &tx.Provider, &tx.Type, &tx.Amount, &tx.Currency, &tx.Status,
		&tx.CustomerRef, &tx.IdempotencyKey, &tx.Reference, &traceID, &itemsJSON, &createdMS, &updatedMS,
	)
	if err == sql.ErrNoRows {
		return Transaction{}, false, nil
	}
	if err != nil {
		return Transaction{}, false, fmt.Errorf("scan transaction: %w", err)
	}

	if traceID.Valid {
		tx.TraceID = traceID.String
	}
	if err := json.Unmarshal([]byte(itemsJSON), &tx.Items); err != nil {
		return Transaction{}, false, fmt.Errorf("unmarshal items: %w", err)
	}
	tx.CreatedAt = time.UnixMilli(createdMS)
	tx.UpdatedAt = time.UnixMilli(updatedMS)
	return tx, true, nil
}

func scanTransactions(rows *sql.Rows) ([]Transaction, error) {
	var result []Transaction
	for rows.Next() {
		var tx Transaction
		var itemsJSON string
		var createdMS, updatedMS int64
		var traceID sql.NullString

		err := rows.Scan(
			&tx.ID, &tx.Provider, &tx.Type, &tx.Amount, &tx.Currency, &tx.Status,
			&tx.CustomerRef, &tx.IdempotencyKey, &tx.Reference, &traceID, &itemsJSON, &createdMS, &updatedMS,
		)
		if err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		if traceID.Valid {
			tx.TraceID = traceID.String
		}
		if err := json.Unmarshal([]byte(itemsJSON), &tx.Items); err != nil {
			return nil, fmt.Errorf("unmarshal items: %w", err)
		}
		tx.CreatedAt = time.UnixMilli(createdMS)
		tx.UpdatedAt = time.UnixMilli(updatedMS)
		result = append(result, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transactions: %w", err)
	}
	return result, nil
}
