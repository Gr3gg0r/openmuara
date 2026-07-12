# OpenMuara iPay88 — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | iPay88 API or signature algorithm changes after implementation. | Low | Medium | Scope to documented subset; keep contract/golden tests. |
| R02 | SHA256 signature algorithm differs from real iPay88 due to field ordering, amount normalization, or missing `SignatureType`. | Medium | High | Match documented algorithm exactly; test against known examples. |
| R03 | Existing provider tests break when adding new routes. | Low | High | Run full `go test ./...` after every sub-step. |
| R04 | ResponseURL/BackendURL SSRF if not validated. | Low | High | Validate URLs are HTTP(S) and reject private/internal hosts. |
