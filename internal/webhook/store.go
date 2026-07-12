package webhook

import (
	"sync"
	"time"
)

// AttemptStatus represents the delivery state of a webhook attempt.
type AttemptStatus string

const (
	// AttemptStatusPending means the webhook has not been delivered yet.
	AttemptStatusPending AttemptStatus = "pending"
	// AttemptStatusDelivered means the webhook was accepted (2xx response).
	AttemptStatusDelivered AttemptStatus = "delivered"
	// AttemptStatusFailed means all retries were exhausted.
	AttemptStatusFailed AttemptStatus = "failed"
)

// AttemptHistory records one delivery try (initial or retry) for a webhook.
type AttemptHistory struct {
	Time   time.Time `json:"time"`
	Status int       `json:"status"`
	Error  string    `json:"error,omitempty"`
}

// Attempt records a webhook delivery attempt.
type Attempt struct {
	ID             string            `json:"id"`
	Ref            string            `json:"ref"`
	ProviderName   string            `json:"provider_name,omitempty"`
	URL            string            `json:"url"`
	Status         AttemptStatus     `json:"status"`
	Attempts       int               `json:"attempts"`
	LastError      string            `json:"last_error,omitempty"`
	Payload        []byte            `json:"payload,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	TraceID        string            `json:"trace_id,omitempty"`
	SignatureValid *bool             `json:"signature_valid,omitempty"`
	History        []AttemptHistory  `json:"attempt_events,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// AttemptStore persists and retrieves webhook attempts.
type AttemptStore interface {
	Save(a *Attempt) error
	Get(ref string) (*Attempt, error)
	// List returns the most recent attempts up to limit, skipping offset results.
	// A limit of zero or less returns all remaining items.
	List(limit, offset int) ([]*Attempt, error)

	// Clear removes all attempts from the store.
	Clear() error
}

// MemoryStore is an in-memory AttemptStore.
type MemoryStore struct {
	mu   sync.RWMutex
	refs map[string]*Attempt
	list []*Attempt
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		refs: make(map[string]*Attempt),
		list: make([]*Attempt, 0),
	}
}

// Clear removes all attempts from the store.
func (s *MemoryStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refs = make(map[string]*Attempt)
	s.list = make([]*Attempt, 0)
	return nil
}

// Save stores or updates an attempt.
func (s *MemoryStore) Save(a *Attempt) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	if existing, ok := s.refs[a.Ref]; ok {
		*existing = *a
	} else {
		clone := *a
		s.refs[a.Ref] = &clone
		s.list = append(s.list, &clone)
	}

	return nil
}

// Get retrieves an attempt by reference.
func (s *MemoryStore) Get(ref string) (*Attempt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	a, ok := s.refs[ref]
	if !ok {
		return nil, nil
	}

	cloned := *a
	return &cloned, nil
}

// List returns the most recent attempts up to limit, skipping offset results.
// A limit of zero or less returns all remaining items.
func (s *MemoryStore) List(limit, offset int) ([]*Attempt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset < 0 {
		offset = 0
	}
	if offset >= len(s.list) {
		return []*Attempt{}, nil
	}

	start := len(s.list) - 1 - offset
	result := make([]*Attempt, 0)
	if limit > 0 {
		result = make([]*Attempt, 0, limit)
	}
	for i := start; i >= 0; i-- {
		if limit > 0 && len(result) >= limit {
			break
		}
		cloned := *s.list[i]
		result = append(result, &cloned)
	}

	return result, nil
}
