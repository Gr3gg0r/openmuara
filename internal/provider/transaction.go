// Package provider defines a minimal payment transaction model used by providers.
// When a dedicated transaction package is introduced, this can become an alias.
package provider

// Transaction is a minimal payment transaction model used by providers.
// When internal/transaction is available, this can be replaced by an alias.
type Transaction struct {
	ID        string
	Reference string
	Amount    float64
	Currency  string
	Status    string
}
