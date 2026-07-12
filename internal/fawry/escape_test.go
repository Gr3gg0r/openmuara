package fawry

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/provider"
	"github.com/Gr3gg0r/openmuara/internal/webhook"
)

func TestEscapeHandlerRequiresGet(t *testing.T) {
	handler := NewEscapeHandler()
	req := httptest.NewRequest(http.MethodPost, "/escape?ref=r&returnUrl=http://localhost", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestEscapeHandlerRequiresRefAndReturnURL(t *testing.T) {
	handler := NewEscapeHandler()
	req := httptest.NewRequest(http.MethodGet, "/escape?ref=r", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", rec.Code)
	}
}

func TestEscapeHandlerRendersPage(t *testing.T) {
	handler := NewEscapeHandler()
	req := httptest.NewRequest(http.MethodGet, "/escape?ref=r&returnUrl=http://localhost&amount=10.00", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "r") {
		t.Error("expected body to contain ref")
	}
}

func TestBuildCallbackURLAddsStatus(t *testing.T) {
	got, err := buildCallbackURL("http://localhost/callback", "PAID")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "orderStatus=PAID") {
		t.Errorf("expected orderStatus=PAID in %q", got)
	}
	if !strings.Contains(got, "statusCode=200") {
		t.Errorf("expected statusCode=200 in %q", got)
	}
}

func TestBuildCallbackURLRejectsInvalidURL(t *testing.T) {
	_, err := buildCallbackURL("://invalid", "PAID")
	if err == nil {
		t.Fatal("expected error for invalid url")
	}
}

func TestEscapeActionHandlerUpdatesLedgerToPaid(t *testing.T) {
	ledger := engine.NewMemoryStore()
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: "fawry", Type: "charge", Reference: "ref-1", Amount: 10.0, Currency: "EGP", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewEscapeActionHandler(nil, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=ref-1&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", rec.Code, rec.Body.String())
	}

	tx, ok, err := ledger.GetByReference("ref-1")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok || tx.Status != engine.TransactionStatusPaid {
		t.Errorf("ledger status: want paid, got %q", tx.Status)
	}
}

func TestEscapeActionHandlerUpdatesLedgerToUnpaid(t *testing.T) {
	ledger := engine.NewMemoryStore()
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: "fawry", Type: "charge", Reference: "ref-2", Amount: 10.0, Currency: "EGP", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewEscapeActionHandler(nil, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=ref-2&returnUrl=http://localhost/callback&status=CANCELED"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: want 303, got %d, body: %s", rec.Code, rec.Body.String())
	}

	tx, ok, err := ledger.GetByReference("ref-2")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !ok || tx.Status != engine.TransactionStatusUnpaid {
		t.Errorf("ledger status: want unpaid, got %q", tx.Status)
	}
}

func TestEscapeActionHandlerReturns404ForMissingReference(t *testing.T) {
	ledger := engine.NewMemoryStore()
	handler := NewEscapeActionHandler(nil, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=missing&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

type getErrStore struct {
	engine.TransactionStore
	err error
}

func (e *getErrStore) GetByReference(string) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, e.err
}

type createErrStore struct {
	engine.TransactionStore
	err error
}

func (e *createErrStore) CreateOrGet(engine.Transaction) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, e.err
}

func TestEscapeActionHandlerMethodNotAllowed(t *testing.T) {
	handler := NewEscapeActionHandler(nil, engine.NewMemoryStore())
	req := httptest.NewRequest(http.MethodGet, "/_admin/fawry-escape", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: want 405, got %d", rec.Code)
	}
}

func TestEscapeActionHandlerMissingFields(t *testing.T) {
	handler := NewEscapeActionHandler(nil, engine.NewMemoryStore())
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=r&returnUrl=http://localhost"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestEscapeActionHandlerNilLedger(t *testing.T) {
	handler := NewEscapeActionHandler(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=r&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestEscapeActionHandlerGetByReferenceError(t *testing.T) {
	handler := NewEscapeActionHandler(nil, &getErrStore{err: errors.New("boom")})
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=r&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d", rec.Code)
	}
}

func TestEscapeActionHandlerCreateOrGetError(t *testing.T) {
	ledger := engine.NewMemoryStore()
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: "fawry", Type: "charge", Reference: "ref-err", Amount: 10.0, Currency: "EGP", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewEscapeActionHandler(nil, &createErrStore{TransactionStore: ledger, err: errors.New("boom")})
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=ref-err&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestEscapeActionHandlerTransitionConflict(t *testing.T) {
	ledger := engine.NewMemoryStore()
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: "fawry", Type: "charge", Reference: "ref-unpaid", Amount: 10.0, Currency: "EGP", Status: engine.TransactionStatusUnpaid}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewEscapeActionHandler(nil, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=ref-unpaid&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status: want 409, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestEscapeActionHandlerDispatchError(t *testing.T) {
	ledger := engine.NewMemoryStore()
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: "fawry", Type: "charge", Reference: "ref-disp", Amount: 10.0, Currency: "EGP", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	d := &webhook.Dispatcher{
		URL: "http://localhost/webhook",
		Builder: func(context.Context, provider.Transaction) ([]byte, error) {
			return nil, errors.New("build failed")
		},
		Store: webhook.NewMemoryStore(),
	}
	handler := NewEscapeActionHandler(d, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=ref-disp&returnUrl=http://localhost/callback&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestEscapeActionHandlerInvalidReturnURL(t *testing.T) {
	ledger := engine.NewMemoryStore()
	if _, _, err := ledger.CreateOrGet(engine.Transaction{Provider: "fawry", Type: "charge", Reference: "ref-url", Amount: 10.0, Currency: "EGP", Status: engine.TransactionStatusNew}); err != nil {
		t.Fatalf("create: %v", err)
	}

	handler := NewEscapeActionHandler(nil, ledger)
	req := httptest.NewRequest(http.MethodPost, "/_admin/fawry-escape", strings.NewReader("ref=ref-url&returnUrl=://bad&status=PAID"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want 500, got %d, body: %s", rec.Code, rec.Body.String())
	}
}
