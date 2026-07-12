# OpenMuara MKP Fawry Integration — Risk Register

| ID | Risk | Likelihood | Impact | Mitigation | Status |
|---|---|---|---|---|---|
| R01 | Adding `canceled`/`expired` states breaks existing state-machine assumptions | Medium | High | Add transitions carefully; fuzz test all state pairs | ⬜ |
| R02 | `response_delay_ms` delays smoke tests or CI | Medium | Medium | Default to 0; only enable when configured | ⬜ |
| R03 | `billing_type` payload shape diverges from real Fawry | Low | Medium | Document as OpenMuara extension, not real-provider contract | ⬜ |
| R04 | Escape page UI becomes cluttered with new options | Low | Low | Use progressive disclosure / collapsible advanced options | ⬜ |
