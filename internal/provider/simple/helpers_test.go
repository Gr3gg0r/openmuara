package simple

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"testing"
)

func computeFawrySignature(t *testing.T, body map[string]any, secret string) string {
	t.Helper()
	mc := body["merchantCode"].(string)
	ref := body["merchantRefNum"].(string)
	cust := ""
	if c, ok := body["customerProfileId"].(string); ok {
		cust = c
	}
	ret := body["returnUrl"].(string)

	rawItems := body["chargeItems"].([]any)
	type item struct {
		ID       string
		Price    float64
		Quantity int
	}
	items := make([]item, len(rawItems))
	for i, raw := range rawItems {
		m := raw.(map[string]any)
		qty := m["quantity"]
		qtyInt := 0
		switch q := qty.(type) {
		case int:
			qtyInt = q
		case float64:
			qtyInt = int(q)
		}
		items[i] = item{
			ID:       m["itemId"].(string),
			Price:    m["price"].(float64),
			Quantity: qtyInt,
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })

	var itemPart string
	for _, it := range items {
		itemPart += fmt.Sprintf("%s%d%s", it.ID, it.Quantity, fmt.Sprintf("%.2f", it.Price))
	}

	text := fmt.Sprintf("%s%s%s%s%s%s", mc, ref, cust, ret, itemPart, secret)
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
