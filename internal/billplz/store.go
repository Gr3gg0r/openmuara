package billplz

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// collectionStore holds collections in memory.
type collectionStore struct {
	mu   sync.RWMutex
	byID map[string]Collection
}

func newCollectionStore() *collectionStore {
	return &collectionStore{byID: make(map[string]Collection)}
}

func (s *collectionStore) create(req CreateCollectionRequest) Collection {
	now := time.Now()
	c := Collection{
		ID:        uuid.Must(uuid.NewRandom()).String(),
		Title:     req.Title,
		Logo:      req.Logo,
		Status:    "active",
		Region:    "MY",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[c.ID] = c
	return c
}

func (s *collectionStore) get(id string) (Collection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.byID[id]
	return c, ok
}

// billStore holds bills in memory.
type billStore struct {
	mu   sync.RWMutex
	byID map[string]Bill
}

func newBillStore() *billStore {
	return &billStore{byID: make(map[string]Bill)}
}

func (s *billStore) create(req CreateBillRequest, baseURL string) Bill {
	now := time.Now()
	id := uuid.Must(uuid.NewRandom()).String()
	b := Bill{
		ID:              id,
		CollectionID:    req.CollectionID,
		Paid:            false,
		State:           BillStateDue,
		Amount:          req.Amount,
		Description:     req.Description,
		Name:            req.Name,
		Email:           req.Email,
		Mobile:          req.Mobile,
		Reference1:      req.Reference1,
		Reference1Label: req.Reference1Label,
		Reference2:      req.Reference2,
		Reference2Label: req.Reference2Label,
		CallbackURL:     req.CallbackURL,
		RedirectURL:     req.RedirectURL,
		URL:             baseURL + "/_admin/billplz/pay/" + id,
		DueAt:           &now,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[b.ID] = b
	return b
}

func (s *billStore) get(id string) (Bill, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.byID[id]
	return b, ok
}

func (s *billStore) update(b Bill) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[b.ID] = b
}

func (s *billStore) delete(id string) (Bill, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.byID[id]
	if !ok {
		return Bill{}, false
	}
	b.State = BillStateDeleted
	s.byID[id] = b
	return b, true
}
