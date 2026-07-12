package audit

import "time"

// Store persists and retrieves audit events.
type Store interface {
	Save(event Event) error
	List(limit, offset int) ([]Event, error)
	ListSince(limit, offset int, since time.Time) ([]Event, error)
	Clear() error
}
