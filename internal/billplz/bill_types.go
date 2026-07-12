package billplz

import "time"

// BillState represents the lifecycle state of a Billplz bill.
type BillState string

const (
	// BillStateDue means the bill has not been paid.
	BillStateDue BillState = "due"
	// BillStatePaid means the bill has been paid.
	BillStatePaid BillState = "paid"
	// BillStateDeleted means the bill was deleted.
	BillStateDeleted BillState = "deleted"
)

// Bill represents a Billplz v3 bill.
type Bill struct {
	ID              string     `json:"id"`
	CollectionID    string     `json:"collection_id"`
	Paid            bool       `json:"paid"`
	State           BillState  `json:"state"`
	Amount          int64      `json:"amount"`
	Description     string     `json:"description"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Mobile          string     `json:"mobile,omitempty"`
	Reference1      string     `json:"reference_1,omitempty"`
	Reference1Label string     `json:"reference_1_label,omitempty"`
	Reference2      string     `json:"reference_2,omitempty"`
	Reference2Label string     `json:"reference_2_label,omitempty"`
	CallbackURL     string     `json:"callback_url"`
	RedirectURL     string     `json:"redirect_url,omitempty"`
	URL             string     `json:"url"`
	PaidAmount      *int64     `json:"paid_amount"`
	DueAt           *time.Time `json:"due_at"`
	PaidAt          *time.Time `json:"paid_at"`
}

// CreateBillRequest is the body for POST /api/v3/bills.
type CreateBillRequest struct {
	CollectionID    string `json:"collection_id"`
	Email           string `json:"email"`
	Mobile          string `json:"mobile"`
	Name            string `json:"name"`
	Amount          int64  `json:"amount"`
	CallbackURL     string `json:"callback_url"`
	Description     string `json:"description"`
	RedirectURL     string `json:"redirect_url"`
	Reference1      string `json:"reference_1"`
	Reference1Label string `json:"reference_1_label"`
	Reference2      string `json:"reference_2"`
	Reference2Label string `json:"reference_2_label"`
}

// BillResponse wraps a single bill for JSON responses.
type BillResponse struct {
	Bill Bill `json:"bill"`
}
