// Package toyyibpay emulates the ToyyibPay payment gateway.
package toyyibpay

// Category represents a ToyyibPay category.
type Category struct {
	CategoryCode        string `json:"categoryCode"`
	CategoryName        string `json:"categoryName"`
	CategoryDescription string `json:"categoryDescription"`
	CategoryStatus      string `json:"categoryStatus"`
}

// CategoryResponse wraps the create category response.
type CategoryResponse struct {
	Status string   `json:"status"`
	Msg    string   `json:"msg"`
	Data   Category `json:"data"`
}
