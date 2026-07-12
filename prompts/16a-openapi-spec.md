# Prompt 16a — OpenAPI Specification

## Goal
Produce and validate an OpenAPI 3.x spec for OpenMuara public endpoints.

## Acceptance Criteria
- [ ] `openapi/openmuara-v1.yaml` generated/maintained
- [ ] Covers all provider-matched paths and admin endpoints
- [ ] Components: schemas, parameters, responses, security schemes
- [ ] Validation in CI (`swagger-codegen validate` or `redocly lint`)
- [ ] Spec served at `GET /_admin/openapi.json`

## Files to Create/Change
- `openapi/openmuara-v1.yaml`
- `.github/workflows/ci.yml` — spec validation
- `internal/server/router.go` — `/openapi.json` route

## Response Shape
Return:
1. Spec version and base URL
2. Endpoint count by tag
3. Validation command

## Test Notes
- Validate spec with chosen tool
- Verify `/openapi.json` returns JSON
