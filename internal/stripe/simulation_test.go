package stripe

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmuara/openmuara/internal/engine"
)

func TestFailureSimulationHandler(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	session := &CheckoutSession{ID: "cs_fail_1", Status: "open", PaymentStatus: "unpaid"}
	sessions.Save(session.ID, session)
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "checkout_session", Reference: session.ID, Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewFailureSimulationHandler(sessions, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/fail?session_id=cs_fail_1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	stored, ok := sessions.Load("cs_fail_1")
	if !ok {
		t.Fatal("expected session")
	}
	if stored.PaymentStatus != "unpaid" {
		t.Errorf("payment status: want unpaid, got %q", stored.PaymentStatus)
	}

	tx, ok, err := ledger.GetByReference("cs_fail_1")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok || tx.Status != engine.TransactionStatusUnpaid {
		t.Errorf("ledger status: want unpaid, got %q", tx.Status)
	}
}

func TestCancelSimulationHandler(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	session := &CheckoutSession{ID: "cs_cancel_1", Status: "open", PaymentStatus: "unpaid"}
	sessions.Save(session.ID, session)
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "checkout_session", Reference: session.ID, Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewCancelSimulationHandler(sessions, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/cancel?session_id=cs_cancel_1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}

	stored, ok := sessions.Load("cs_cancel_1")
	if !ok || stored.Status != "canceled" {
		t.Errorf("session status: want canceled, got %q", stored.Status)
	}
}

func TestSimulationHandlerMissingSession(t *testing.T) {
	sessions := NewMemorySessionStore()
	ledger := engine.NewMemoryStore()
	handler := NewFailureSimulationHandler(sessions, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/fail?session_id=missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: want 404, got %d", rec.Code)
	}
}

func TestSimulationHandlerMethodNotAllowed(t *testing.T) {
	handler := NewFailureSimulationHandler(NewMemorySessionStore(), engine.NewMemoryStore())
	req := httptest.NewRequest(http.MethodGet, "/_admin/stripe/fail?session_id=cs_1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestSimulationHandlerMissingSessionID(t *testing.T) {
	handler := NewFailureSimulationHandler(NewMemorySessionStore(), engine.NewMemoryStore())
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/fail", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestSimulationHandlerLedgerLookupError(t *testing.T) {
	sessions := NewMemorySessionStore()
	sessions.Save("cs_err", &CheckoutSession{ID: "cs_err"})

	handler := NewFailureSimulationHandler(sessions, &errStore{err: errors.New("boom")})
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/fail?session_id=cs_err", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestSimulationHandlerTransitionConflict(t *testing.T) {
	sessions := NewMemorySessionStore()
	sessions.Save("cs_conflict", &CheckoutSession{ID: "cs_conflict"})

	ledger := engine.NewMemoryStore()
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: ProviderName, Type: "checkout_session", Reference: "cs_conflict", Amount: 10.0, Currency: "USD", Status: engine.TransactionStatusPaid}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewFailureSimulationHandler(sessions, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/stripe/fail?session_id=cs_conflict", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status: want 409, got %d", rec.Code)
	}
}

type errStore struct {
	engine.TransactionStore
	err error
}

func (e *errStore) GetByReference(string) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, e.err
}
