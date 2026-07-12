package simple

import (
	"crypto/hmac"
	"crypto/sha256"

	// #nosec G501 -- simple runtime emulates providers that use MD5
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/Gr3gg0r/openmuara/internal/plugin"
)

func (p *Provider) verifySignature(values map[string]any) bool {
	if p.cfg.Signature == nil {
		return true
	}

	given, _ := signatureValue(values, p.cfg.Signature.Fields)
	if given == "" {
		return false
	}

	expected := p.sign(values)
	return strings.EqualFold(given, expected)
}

func (p *Provider) sign(values map[string]any) string {
	if p.cfg.Signature == nil {
		return ""
	}

	switch p.cfg.Signature.Algorithm {
	case "fawry_sha256":
		return p.signFawrySHA256(values)
	case "hmac_sha256":
		return p.signHMACSHA256(values)
	case "md5_concat":
		return p.signMD5Concat(values)
	case "senangpay_md5":
		return p.signSenangpayMD5(values)
	default:
		return ""
	}
}

func (p *Provider) signFawrySHA256(values map[string]any) string {
	// Fawry canonical string:
	// merchantCode + merchantRefNum + customerProfileId + returnUrl +
	// sorted(itemId + quantity + price(2 decimals)) + secureKey
	merchantCode := coalesceString(values, "merchant_code", "merchantCode")
	merchantRefNum := coalesceString(values, "merchant_ref_num", "merchantRefNum")
	customerProfileID := coalesceString(values, "customer_profile_id", "customerProfileId")
	returnURL := coalesceString(values, "return_url", "returnUrl")

	items := itemList(values)
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })

	var itemPart string
	for _, it := range items {
		itemPart += fmt.Sprintf("%s%d%s", it.ID, it.Quantity, fmt.Sprintf("%.2f", it.Price))
	}

	text := fmt.Sprintf("%s%s%s%s%s%s",
		merchantCode,
		merchantRefNum,
		customerProfileID,
		returnURL,
		itemPart,
		p.secret,
	)
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (p *Provider) signHMACSHA256(values map[string]any) string {
	// Build a flat key=value map from the configured signature fields.
	flat := make(map[string]string, len(p.cfg.Signature.Fields))
	for _, f := range p.cfg.Signature.Fields {
		if s, ok := stringValue(values, f); ok {
			flat[f] = s
			continue
		}
		if s, ok := stringValue(values, camelCase(f)); ok {
			flat[f] = s
		}
	}

	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+flat[k])
	}
	msg := strings.Join(pairs, "|")

	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func (p *Provider) signMD5Concat(values map[string]any) string {
	var parts []string
	for _, f := range p.cfg.Signature.Fields {
		s := coalesceString(values, f, camelCase(f))
		if s != "" {
			parts = append(parts, s)
		}
	}
	msg := p.secret + strings.Join(parts, "")
	// #nosec G401 -- MD5 is used only to emulate provider signatures
	sum := md5.Sum([]byte(msg))
	return fmt.Sprintf("%x", sum[:])
}

func (p *Provider) signSenangpayMD5(values map[string]any) string {
	detail := coalesceString(values, "detail")
	amount := coalesceFloat(values, "amount")
	orderID := coalesceString(values, "order_id", "orderId")
	msg := fmt.Sprintf("%s%s%.2f%s", p.secret, detail, amount, orderID)
	// #nosec G401 -- MD5 is used only to emulate provider signatures
	sum := md5.Sum([]byte(msg))
	return fmt.Sprintf("%x", sum[:])
}

func signatureValue(values map[string]any, fields []string) (string, bool) {
	// The signature field is conventionally the last configured field, or any
	// field named signature/hash.
	for _, f := range fields {
		if f == "signature" || f == "hash" {
			if s, ok := stringValue(values, f); ok {
				return s, true
			}
			if s, ok := stringValue(values, camelCase(f)); ok {
				return s, true
			}
		}
	}
	if len(fields) > 0 {
		f := fields[len(fields)-1]
		if s, ok := stringValue(values, f); ok {
			return s, true
		}
		return stringValue(values, camelCase(f))
	}
	return "", false
}

type item struct {
	ID       string
	Price    float64
	Quantity int
}

func itemList(values map[string]any) []item {
	raw, ok := values["charge_items"]
	if !ok {
		raw, ok = values["chargeItems"]
	}
	if !ok {
		return nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil
	}

	var items []item
	for _, r := range list {
		m, ok := r.(map[string]any)
		if !ok {
			continue
		}
		id := coalesceString(m, "itemId", "item_id")
		price := coalesceFloat(m, "price")
		qty := int(coalesceFloatOrInt(m, "quantity"))
		items = append(items, item{ID: id, Price: price, Quantity: qty})
	}
	return items
}

func coalesceString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if s, ok := stringValue(m, k); ok {
			return s
		}
	}
	return ""
}

func coalesceFloat(m map[string]any, keys ...string) float64 {
	for _, k := range keys {
		if f, ok := floatValue(m, k); ok {
			return f
		}
	}
	return 0
}

func coalesceFloatOrInt(m map[string]any, keys ...string) float64 {
	for _, k := range keys {
		if f, ok := floatValue(m, k); ok {
			return f
		}
		if i, ok := intValue(m, k); ok {
			return float64(i)
		}
	}
	return 0
}

func camelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func resolveSecret(name string, sig *plugin.Signature, cfg map[string]any) (string, error) {
	if sig == nil {
		return "", nil
	}
	if sig.SecretEnv != "" {
		return "", fmt.Errorf("secret_env not supported by simple runtime")
	}
	if sig.SecretKey == "" {
		return "", nil
	}
	key := sig.SecretKey
	// Gateway manifests use a full dotted path such as
	// "providers.<name>.config.<key>" for validation, but the simple runtime
	// only receives the provider's own config map. Strip the prefix when it
	// matches the current provider.
	prefix := "providers." + name + ".config."
	key = strings.TrimPrefix(key, prefix)
	return dottedLookup(cfg, key), nil
}

func dottedLookup(cfg map[string]any, path string) string {
	parts := strings.Split(path, ".")
	current := cfg
	for i, part := range parts {
		if i == len(parts)-1 {
			if v, ok := current[part]; ok {
				if s, ok := v.(string); ok {
					return s
				}
			}
			return ""
		}
		next, ok := current[part].(map[string]any)
		if !ok {
			return ""
		}
		current = next
	}
	return ""
}
