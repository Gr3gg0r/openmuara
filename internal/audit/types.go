// Package audit provides structured audit logging for OpenMuara.
package audit

import "time"

// Event is a single structured audit log entry.
type Event struct {
	ID           string    `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	Actor        string    `json:"actor"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Payload      string    `json:"payload,omitempty"`
	Result       string    `json:"result,omitempty"`
}
