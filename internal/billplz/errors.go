package billplz

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/errcode"
)

var (
	errCollectionNotFound = errcode.New(errcode.ETransactionNotFound, "collection not found")
	errInvalidAmount      = errcode.New(errcode.EInvalidRequest, "amount must be greater than zero")
)

func errRequiredField(name string) error {
	return errcode.New(errcode.EInvalidRequest, fmt.Sprintf("%s is required", name))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(v)
}
