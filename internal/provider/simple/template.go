package simple

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/openmuara/openmuara/internal/engine"
)

// templateData holds the values available to response and webhook templates.
type templateData struct {
	Reference  string
	Status     string
	Amount     float64
	Currency   string
	Provider   string
	RequestID  string
	EventID    string
	StatusID   string
	Request    map[string]any
	RawRequest map[string]any
}

func renderTemplateData(tx engine.Transaction) templateData {
	return templateData{
		Reference: tx.Reference,
		Status:    string(tx.Status),
		Amount:    tx.Amount,
		Currency:  string(tx.Currency),
		Provider:  tx.Provider,
		RequestID: tx.TraceID,
		EventID:   tx.TraceID,
		StatusID:  statusID(string(tx.Status)),
	}
}

func (p *Provider) renderResponse(values map[string]any, tx engine.Transaction) map[string]any {
	if len(p.runtime.ResponseTemplate) == 0 {
		return map[string]any{
			"status":    "ok",
			"reference": tx.Reference,
		}
	}

	data := renderTemplateData(tx)
	data.RawRequest = values
	data.Request = values
	return renderMap(p.runtime.ResponseTemplate, data)
}

func renderMap(src map[string]any, data templateData) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = renderValue(v, data)
	}
	return out
}

func renderValue(v any, data templateData) any {
	switch x := v.(type) {
	case string:
		return renderString(x, data)
	case map[string]any:
		return renderMap(x, data)
	case []any:
		arr := make([]any, len(x))
		for i, item := range x {
			arr[i] = renderValue(item, data)
		}
		return arr
	default:
		return v
	}
}

func renderString(s string, data templateData) string {
	if !strings.Contains(s, "{{") {
		return s
	}
	tmpl, err := template.New("t").Parse(s)
	if err != nil {
		return s
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return s
	}
	return buf.String()
}

func statusID(status string) string {
	switch strings.ToLower(status) {
	case "paid", "completed", "success":
		return "1"
	case "unpaid", "failed", "failure":
		return "0"
	default:
		return "0"
	}
}
