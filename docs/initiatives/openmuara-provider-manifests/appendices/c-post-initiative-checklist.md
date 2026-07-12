> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Appendix C — Post-Initiative Checklist

Complete these steps after P05 (Final gates & PR).

---

## Final Sweep

- [ ] `TRACKING.md` reflects all completed prompts and gate results.
- [ ] `HANDOFF.md` is updated with final state.
- [ ] `DECISIONS.md` has no open decisions.
- [ ] `RISKS.md` has no active high-score risks.
- [ ] `KNOWN_ISSUES.md` is updated; resolved issues are marked.

## Documentation

- [ ] `docs/provider-contract.md` is accurate.
- [ ] `docs/contributing-providers.md` is accurate.
- [ ] Migration guide exists at `docs/migration/provider-manifests.md` if needed.
- [ ] `CHANGELOG.md` has a release-notes snippet.

## Quality

- [ ] `go build ./...` passes.
- [ ] `go test ./...` passes.
- [ ] `go vet ./...` passes.
- [ ] `golangci-lint run ./...` passes.
- [ ] `go test -race ./...` passes.
- [ ] Coverage maintained on changed modules.

## End-to-End

- [ ] `muara provider validate plugins/*/gateway.yml` passes.
- [ ] `examples/checkout-store` Fawry flow works.
- [ ] `examples/checkout-store` Stripe flow works.

## PR & Handoff

- [ ] Branch is rebased on latest `dev`.
- [ ] `REVIEW_CHECKLIST.md` is complete.
- [ ] PR is opened from `dev` to `main` (or per project flow).
- [ ] Human reviewer is notified with summary of changes and risks.

## After Merge

- [ ] Update root `TRACKING.md` / v1 master backlog.
- [ ] Close the initiative in `TRACKING.md`.
- [ ] Archive or update any related prompt/task files if needed.
