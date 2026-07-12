package billplz_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

// errStore is a TransactionStore that always returns an error.
type errStore struct{}

func (errStore) CreateOrGet(engine.Transaction) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, errors.New("store error")
}
func (errStore) GetByID(string) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, errors.New("store error")
}
func (errStore) GetByReference(string) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, errors.New("store error")
}
func (errStore) List(int, int) ([]engine.Transaction, error) {
	return nil, errors.New("store error")
}
func (errStore) Clear() error {
	return errors.New("store error")
}

func TestChargeHandler(t *testing.T) {
	p := newInitializedProvider(t)
	handler := p.ChargeHandler()

	req := httptest.NewRequest(http.MethodGet, "/charge", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
}

func TestEscapeHandler(t *testing.T) {
	p := newInitializedProvider(t)
	if p.EscapeHandler() != nil {
		t.Fatal("expected EscapeHandler to return nil")
	}
}

func TestUpdateTransactionWithErrStore(t *testing.T) {
	p := newInitializedProvider(t)
	c := createCollection(t, p)
	b := createBill(t, p, c.ID)

	p.SetStore(errStore{})
	handler := routeFor(t, p, http.MethodPost, "/_admin/billplz/pay/{id}")
	req := mustPayRequest(b.ID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 200 or 303, got %d", rec.Code)
	}
}
