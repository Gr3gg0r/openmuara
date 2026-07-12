package audit

import (
	"sync"
	"time"
)

// MemoryStore is an in-memory audit store used for tests and memory-mode runs.
type MemoryStore struct {
	mu     sync.RWMutex
	events []Event
}

// NewMemoryStore creates an empty in-memory audit store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{events: make([]Event, 0)}
}

// Save appends an event to the in-memory store.
func (s *MemoryStore) Save(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
	return nil
}

// List returns the most recent events, skipping offset and limiting to limit.
func (s *MemoryStore) List(limit, offset int) ([]Event, error) {
	return s.ListSince(limit, offset, time.Time{})
}

// ListSince returns events at or after since, ordered newest first.
func (s *MemoryStore) ListSince(limit, offset int, since time.Time) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = len(s.events)
	}

	filtered := make([]Event, 0)
	for i := len(s.events) - 1; i >= 0; i-- {
		if !since.IsZero() && s.events[i].Timestamp.Before(since) {
			continue
		}
		filtered = append(filtered, s.events[i])
	}

	if offset >= len(filtered) {
		return []Event{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}

// Clear removes all events from the store.
func (s *MemoryStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = make([]Event, 0)
	return nil
}
