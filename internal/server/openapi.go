package server

import (
	"embed"
	"net/http"
)

//go:embed openapi.yaml
var openapiSpec embed.FS

// OpenAPIHandler serves the OpenAPI specification document.
func OpenAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		data, err := openapiSpec.ReadFile("openapi.yaml")
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load openapi spec"})
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
