package simple

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/audit"
	"github.com/Gr3gg0r/openmuara/internal/engine"
	"github.com/Gr3gg0r/openmuara/internal/errcode"
	"github.com/Gr3gg0r/openmuara/internal/httputil"
	"github.com/Gr3gg0r/openmuara/internal/plugin"
	"github.com/Gr3gg0r/openmuara/internal/provider"
)

func (p *Provider) handlerFor(r plugin.Route) http.Handler {
	switch {
	case isChargeAction(r.Action):
		return p.chargeHandler(r)
	case isWebhookAction(r.Action):
		return p.webhookHandler(r)
	case isEscapePageAction(r.Action):
		return p.escapePageHandler(r)
	case isEscapeAction(r.Action):
		return p.escapeActionHandler(r)
	case isStatusAction(r.Action):
		return p.statusHandler(r)
	default:
		return p.defaultHandler(r)
	}
}

func (p *Provider) chargeHandler(r plugin.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != r.Method {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		body, err := readBody(req)
		if err != nil {
			httputil.RespondError(w, req, httputil.ErrInvalidJSON, http.StatusBadRequest, "failed to read body")
			return
		}

		var values map[string]any
		if err := json.Unmarshal(body, &values); err != nil {
			httputil.RespondError(w, req, httputil.ErrInvalidJSON, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if err := p.validateRequest(values, r.SchemaRef); err != nil {
			httputil.RespondError(w, req, httputil.ErrMissingField, http.StatusBadRequest, errcode.Message(err))
			return
		}

		if p.cfg.Signature != nil && p.cfg.Signature.Algorithm != "none" {
			if !p.verifySignature(values) {
				httputil.RespondError(w, req, httputil.ErrInvalidSignature, http.StatusBadRequest, "invalid signature")
				return
			}
		}

		ref, _ := stringValue(values, p.runtime.ReferenceField)
		amount := p.amount(values)
		customerRef, _ := stringValue(values, p.runtime.CustomerField)

		tx := engine.NewTransaction(engine.Transaction{
			Provider:       p.name,
			Type:           "charge",
			Amount:         amount,
			Currency:       p.runtime.Currency,
			Status:         engine.TransactionStatusNew,
			CustomerRef:    customerRef,
			IdempotencyKey: req.Header.Get("Idempotency-Key"),
			Reference:      ref,
			TraceID:        httputil.TraceIDFromContext(req.Context()),
		})
		if _, _, err := p.store.CreateOrGet(tx); err != nil {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusInternalServerError, "failed to record transaction")
			return
		}

		audit.FromContext(req.Context()).Log(req.Context(), "charge.created", "transaction", ref, audit.JSON(values), "ok")

		resp := p.renderResponse(values, tx)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func (p *Provider) webhookHandler(r plugin.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != r.Method {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "acknowledged"})
	}
}

func (p *Provider) statusHandler(r plugin.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != r.Method {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ref := req.URL.Query().Get("ref")
		if ref == "" {
			httputil.RespondError(w, req, httputil.ErrMissingField, http.StatusBadRequest, "ref is required")
			return
		}

		tx, ok, err := p.store.GetByReference(ref)
		if err != nil {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusInternalServerError, "failed to lookup transaction")
			return
		}
		if !ok {
			httputil.RespondError(w, req, httputil.ErrNotFound, http.StatusNotFound, "transaction not found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"reference": tx.Reference,
			"status":    tx.Status,
			"amount":    tx.Amount,
			"currency":  tx.Currency,
		})
	}
}

func (p *Provider) defaultHandler(r plugin.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != r.Method {
			httputil.RespondError(w, req, httputil.ErrInternal, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func (p *Provider) buildPayload(_ context.Context, tx provider.Transaction) ([]byte, error) {
	tmpl := p.firstWebhookTemplate()
	if tmpl == nil {
		return json.Marshal(map[string]any{
			"provider":  p.name,
			"reference": tx.Reference,
			"status":    tx.Status,
		})
	}
	data := templateData{
		Reference: tx.Reference,
		Status:    tx.Status,
		Provider:  p.name,
		StatusID:  statusID(tx.Status),
	}
	if p.store != nil {
		if etx, ok, _ := p.store.GetByReference(tx.Reference); ok {
			data.Amount = etx.Amount
			data.Currency = string(etx.Currency)
			data.RequestID = etx.TraceID
			data.EventID = etx.TraceID
		}
	}
	rendered := renderMap(tmpl.Template, data)
	return json.Marshal(rendered)
}

func (p *Provider) firstWebhookTemplate() *plugin.Webhook {
	if len(p.cfg.Webhooks) == 0 {
		return nil
	}
	return &p.cfg.Webhooks[0]
}

func isChargeAction(action string) bool {
	return action == "charge" || action == "simple_charge" || endsWith(action, "_charge")
}

func isWebhookAction(action string) bool {
	return action == "webhook" || action == "simple_webhook" || endsWith(action, "_webhook")
}

func isEscapePageAction(action string) bool {
	return action == "escape_page" || action == "simple_escape_page" || endsWith(action, "_escape_page")
}

func isEscapeAction(action string) bool {
	return action == "escape_action" || action == "simple_escape_action" || endsWith(action, "_escape_action")
}

func isStatusAction(action string) bool {
	return action == "status" || action == "simple_status" || endsWith(action, "_status")
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func readBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}
