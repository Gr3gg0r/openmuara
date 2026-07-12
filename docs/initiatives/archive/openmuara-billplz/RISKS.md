# OpenMuara Billplz — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| R01 | Billplz v3 API shape changes after implementation. | Low | Medium | Scope to documented subset; keep contract/golden tests. |
| R02 | `x_signature` verification differs from real Billplz due to payload ordering or encoding. | Medium | High | Match documented algorithm exactly: sort keys case-insensitively ascending, concatenate `key+value`, join with `|`. Test against known examples. |
| R03 | Bill `url` conflicts with existing admin route patterns. | Low | Low | Use distinct `/_admin/billplz/pay/{id}` path. |
| R04 | Existing provider tests break when adding new routes. | Low | High | Run full `go test ./...` after every sub-step. |
| R05 | Callback vs redirect terminology confuses implementers. | Low | Medium | Name routes and docs to match real Billplz: `callback_url` = server POST, `redirect_url` = browser GET. |
