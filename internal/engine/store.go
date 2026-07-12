// Package engine provides the in-memory transaction ledger for OpenMuara.
package engine

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

var transactionsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "openmuara_transactions_total",
	Help: "Total transactions stored by provider and status.",
}, []string{"provider", "status"})

func init() {
	prometheus.MustRegister(transactionsTotal)
}

func recordTransaction(provider, status string) {
	if provider == "" {
		provider = "unknown"
	}
	if status == "" {
		status = "unknown"
	}
	transactionsTotal.WithLabelValues(provider, string(status)).Inc()
}

// TransactionStore is the read/write contract for the transaction ledger.
type TransactionStore interface {
	// CreateOrGet stores tx if it is new and returns the stored transaction
	// plus a boolean that is true when the transaction was just created.
	// If tx.IdempotencyKey is non-empty and already maps to an existing
	// transaction, the existing transaction is returned and the bool is false.
	CreateOrGet(tx Transaction) (Transaction, bool, error)

	// GetByID returns a transaction by its ID.
	GetByID(id string) (Transaction, bool, error)

	// GetByReference returns the transaction matching the reference.
	GetByReference(ref string) (Transaction, bool, error)

	// List returns the most recent transactions up to limit, skipping the first
	// offset results. Use offset 0 and a negative limit to return all results.
	List(limit, offset int) ([]Transaction, error)

	// Clear removes all transactions from the store.
	Clear() error
}

// MemoryStore is a thread-safe in-memory TransactionStore.
type MemoryStore struct {
	mu            sync.RWMutex
	byID          map[string]Transaction
	byReference   map[string]Transaction
	byIdempotency map[string]string
	order         []string
}

// NewMemoryStore creates a new in-memory transaction store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		byID:          make(map[string]Transaction),
		byReference:   make(map[string]Transaction),
		byIdempotency: make(map[string]string),
		order:         make([]string, 0),
	}
}

// CreateOrGet stores tx and returns the stored transaction. See TransactionStore.
func (s *MemoryStore) CreateOrGet(tx Transaction) (Transaction, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if tx.IdempotencyKey != "" {
		if id, ok := s.byIdempotency[tx.IdempotencyKey]; ok {
			return s.byID[id], false, nil
		}
	}

	if tx.ID == "" {
		tx.ID = uuid.Must(uuid.NewRandom()).String()
	}

	now := time.Now()
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = now
	}
	if tx.UpdatedAt.IsZero() {
		tx.UpdatedAt = now
	}

	// Preserve the original trace ID when a reference is reused.
	if existing, ok := s.byReference[tx.Reference]; ok && existing.TraceID != "" && tx.TraceID == "" {
		tx.TraceID = existing.TraceID
	}

	s.byID[tx.ID] = tx
	s.byReference[tx.Reference] = tx
	if tx.IdempotencyKey != "" {
		s.byIdempotency[tx.IdempotencyKey] = tx.ID
	}
	s.order = append(s.order, tx.ID)

	recordTransaction(tx.Provider, string(tx.Status))
	return tx, true, nil
}

// GetByID returns a transaction by ID.
func (s *MemoryStore) GetByID(id string) (Transaction, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.byID[id]
	return tx, ok, nil
}

// Clear removes all transactions from the store.
func (s *MemoryStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.byID = make(map[string]Transaction)
	s.byReference = make(map[string]Transaction)
	s.byIdempotency = make(map[string]string)
	s.order = make([]string, 0)
	return nil
}

// GetByReference returns a transaction by reference.
func (s *MemoryStore) GetByReference(ref string) (Transaction, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.byReference[ref]
	return tx, ok, nil
}

// List returns the most recent transactions up to limit, skipping offset results.
// A limit of zero or less returns all remaining items.
func (s *MemoryStore) List(limit, offset int) ([]Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset < 0 {
		offset = 0
	}
	if offset >= len(s.order) {
		return []Transaction{}, nil
	}

	start := len(s.order) - 1 - offset
	result := make([]Transaction, 0)
	if limit > 0 {
		result = make([]Transaction, 0, limit)
	}
	for i := start; i >= 0; i-- {
		if limit > 0 && len(result) >= limit {
			break
		}
		result = append(result, s.byID[s.order[i]])
	}
	return result, nil
}
