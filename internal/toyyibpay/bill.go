package toyyibpay

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/google/uuid"
)

// BillStore is a thread-safe in-memory bill repository.
type BillStore struct {
	mu      sync.RWMutex
	byCode  map[string]Bill
	byOrder map[string]string
}

// NewBillStore creates an empty bill store.
func NewBillStore() *BillStore {
	return &BillStore{
		byCode:  make(map[string]Bill),
		byOrder: make(map[string]string),
	}
}

// Create stores a bill and indexes it by order ID.
func (s *BillStore) Create(bill Bill) Bill {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bill.BillCode == "" {
		bill.BillCode = uuid.Must(uuid.NewRandom()).String()
	}
	if bill.OrderID == "" {
		bill.OrderID = bill.BillCode
	}
	s.byCode[bill.BillCode] = bill
	s.byOrder[bill.OrderID] = bill.BillCode
	return bill
}

// GetByCode returns a bill by code.
func (s *BillStore) GetByCode(code string) (Bill, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.byCode[code]
	return b, ok
}

// GetByOrderID returns a bill code by order ID.
func (s *BillStore) GetByOrderID(orderID string) (Bill, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, ok := s.byOrder[orderID]
	if !ok {
		return Bill{}, false
	}
	b, ok := s.byCode[code]
	return b, ok
}

// Update replaces a bill.
func (s *BillStore) Update(bill Bill) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byCode[bill.BillCode] = bill
}

func (p *Provider) billCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}
		if !p.authenticate(r.FormValue("userSecretKey"), w, r) {
			return
		}

		bill, err := p.parseCreateBill(r)
		if err != nil {
			httputil.RespondError(w, r, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}
		bill = p.bills.Create(bill)
		bill.BillPaymentLink = p.baseURL + "/_admin/toyyibpay/pay/" + bill.BillCode
		p.bills.Update(bill)

		p.recordBillTransaction(r.Context(), bill)
		writeJSON(w, BillCreateResponse{Status: "1", Msg: "success", Bill: bill})
	}
}

func (p *Provider) parseCreateBill(r *http.Request) (Bill, error) {
	amount, err := strconv.Atoi(r.FormValue("billAmount"))
	if err != nil || amount <= 0 {
		return Bill{}, errcode.New(errcode.EInvalidRequest, "billAmount must be a positive integer")
	}
	if r.FormValue("billName") == "" {
		return Bill{}, errcode.New(errcode.EInvalidRequest, "billName is required")
	}
	if r.FormValue("billReturnUrl") == "" {
		return Bill{}, errcode.New(errcode.EInvalidRequest, "billReturnUrl is required")
	}
	if r.FormValue("billCallbackUrl") == "" {
		return Bill{}, errcode.New(errcode.EInvalidRequest, "billCallbackUrl is required")
	}

	categoryCode := r.FormValue("categoryCode")
	if categoryCode == "" {
		categoryCode = p.defaultCategory
	}

	channel := r.FormValue("billPaymentChannel")
	if channel == "" {
		channel = "2"
	}

	return Bill{
		BillName:           r.FormValue("billName"),
		BillDescription:    r.FormValue("billDescription"),
		BillTo:             r.FormValue("billTo"),
		BillEmail:          r.FormValue("billEmail"),
		BillPhone:          r.FormValue("billPhone"),
		BillAmount:         amount,
		BillStatus:         "1",
		CategoryCode:       categoryCode,
		BillReturnURL:      r.FormValue("billReturnUrl"),
		BillCallbackURL:    r.FormValue("billCallbackUrl"),
		BillPaymentChannel: channel,
		BillExpiryDate:     r.FormValue("billExpiryDate"),
		BillExpiryDays:     r.FormValue("billExpiryDays"),
		BillPriceSetting:   defaultIfEmpty(r.FormValue("billPriceSetting"), "0"),
		BillPayorInfo:      defaultIfEmpty(r.FormValue("billPayorInfo"), "0"),
		OrderID:            r.FormValue("billExternalReferenceNo"),
	}, nil
}

func (p *Provider) recordBillTransaction(ctx context.Context, bill Bill) {
	if p.store == nil {
		return
	}
	tx := engine.NewTransaction(engine.Transaction{
		Provider:    ProviderName,
		Type:        "bill",
		Amount:      float64(bill.BillAmount) / 100.0,
		Currency:    "MYR",
		Status:      engine.TransactionStatusNew,
		CustomerRef: bill.BillEmail,
		Reference:   bill.OrderID,
		TraceID:     httputil.TraceIDFromContext(ctx),
	})
	_, _, _ = p.store.CreateOrGet(tx)
}

func (p *Provider) billTransactionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}
		if !p.authenticate(r.FormValue("userSecretKey"), w, r) {
			return
		}

		billCode := r.FormValue("billCode")
		bill, ok := p.bills.GetByCode(billCode)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}
		writeJSON(w, p.buildBillTransactions(bill))
	}
}

func (p *Provider) buildBillTransactions(bill Bill) []BillTransaction {
	if p.store == nil {
		return []BillTransaction{}
	}
	tx, ok, _ := p.store.GetByReference(bill.OrderID)
	if !ok {
		return []BillTransaction{}
	}
	status := billPaymentStatusCode(tx.Status)
	return []BillTransaction{{
		BillPaymentStatus:  status,
		BillPaymentChannel: bill.BillPaymentChannel,
		BillPaymentAmount:  strconv.Itoa(bill.BillAmount),
		BillPaymentRefNo:   tx.ID,
		BillPaymentTime:    tx.UpdatedAt.Format("2006-01-02 15:04:05"),
	}}
}

func (p *Provider) billInactiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httputil.RespondError(w, r, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if err := r.ParseForm(); err != nil {
			httputil.RespondError(w, r, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid form")
			return
		}
		if !p.authenticate(r.FormValue("userSecretKey"), w, r) {
			return
		}

		billCode := r.FormValue("billCode")
		bill, ok := p.bills.GetByCode(billCode)
		if !ok {
			httputil.RespondError(w, r, httputil.ErrNotFound, http.StatusNotFound, "bill not found")
			return
		}
		bill.BillStatus = "2"
		p.bills.Update(bill)
		writeJSON(w, map[string]string{"status": "1", "msg": "success"})
	}
}

func defaultIfEmpty(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// billPaymentStatusCode maps an engine transaction status to the ToyyibPay
// bill payment status code: 1 = successful, 2 = pending, 3 = unsuccessful.
func billPaymentStatusCode(status engine.TransactionStatus) string {
	switch status {
	case engine.TransactionStatusPaid:
		return "1"
	case engine.TransactionStatusUnpaid, engine.TransactionStatusRefunded:
		return "3"
	default:
		return "2"
	}
}
