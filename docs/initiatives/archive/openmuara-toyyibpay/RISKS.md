# OpenMuara ToyyibPay — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | ToyyibPay API shape changes after implementation. | Low | Medium | Scope to documented subset; keep contract/golden tests. |
| R02 | MD5 callback hash algorithm differs from real ToyyibPay. | Medium | High | Match documented formula exactly: `MD5(userSecretKey + status + order_id + refno + "ok")`. Test against known examples. |
| R03 | Existing provider tests break when adding new routes. | Low | High | Run full `go test ./...` after every sub-step. |
| R04 | Form-encoded request parsing differs from real ToyyibPay expectations. | Low | Medium | Use standard `application/x-www-form-urlencoded` parsing; test with real field names. |
