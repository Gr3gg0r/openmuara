# OpenMuara v1 Solid Gold — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation | Status |
|---|---|---|---|---|---|
| R01 | Adding CI jobs increases CI duration | Medium | Low | Run jobs in parallel; keep smoke test fast | ⬜ |
| R02 | Coverage backfill writes low-value tests | Medium | Medium | Require meaningful assertions; review manually | ⬜ |
| R03 | P03 trace-ID changes affect webhook signature verification | Low | High | Add header only, do not change payload body | ⬜ |
| R04 | Dashboard UI changes break existing smoke test | Low | Medium | Update smoke test selectors if needed | ⬜ |
| R05 | Stronger linters flag a lot of existing code | Medium | Medium | Enable gradually, grandfather existing issues | ⬜ |
