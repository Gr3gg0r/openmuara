package toyyibpay

// Bill represents a ToyyibPay bill.
type Bill struct {
	BillCode           string `json:"billCode"`
	BillName           string `json:"billName"`
	BillDescription    string `json:"billDescription"`
	BillTo             string `json:"billTo"`
	BillEmail          string `json:"billEmail"`
	BillPhone          string `json:"billPhone"`
	BillAmount         int    `json:"billAmount"`
	BillStatus         string `json:"billStatus"`
	CategoryCode       string `json:"categoryCode"`
	BillReturnURL      string `json:"billReturnUrl"`
	BillCallbackURL    string `json:"billCallbackUrl"`
	BillPaymentChannel string `json:"billPaymentChannel"`
	BillExpiryDate     string `json:"billExpiryDate"`
	BillExpiryDays     string `json:"billExpiryDays"`
	BillPriceSetting   string `json:"billPriceSetting"`
	BillPayorInfo      string `json:"billPayorInfo"`
	BillPaymentLink    string `json:"billPaymentLink"`
	OrderID            string `json:"order_id,omitempty"`
}

// BillCreateResponse is the response shape for createBill.
type BillCreateResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Bill   Bill   `json:"bill"`
}

// BillTransaction represents one payment transaction for a bill.
type BillTransaction struct {
	BillPaymentStatus  string `json:"billpaymentStatus"`
	BillPaymentChannel string `json:"billpaymentChannel"`
	BillPaymentAmount  string `json:"billpaymentAmount"`
	BillPaymentRefNo   string `json:"billpaymentRefNo"`
	BillPaymentTime    string `json:"billpaymentTime"`
}
