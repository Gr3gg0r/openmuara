package billplz

import "time"

// Collection represents a Billplz v3 collection.
type Collection struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Logo      *Logo     `json:"logo,omitempty"`
	Status    string    `json:"status"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Logo represents an optional collection logo object.
type Logo struct {
	ThumbnailURL string `json:"thumbnail_url"`
}

// CreateCollectionRequest is the body for POST /api/v3/collections.
type CreateCollectionRequest struct {
	Title string `json:"title"`
	Logo  *Logo  `json:"logo,omitempty"`
}

// CollectionResponse wraps a single collection for JSON responses.
type CollectionResponse struct {
	Collection Collection `json:"collection"`
}

// PaymentMethod represents one available payment method.
type PaymentMethod struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// PaymentMethodsResponse is the response for GET /api/v3/collections/{id}/payment_methods.
type PaymentMethodsResponse struct {
	PaymentMethods []PaymentMethod `json:"payment_methods"`
}
