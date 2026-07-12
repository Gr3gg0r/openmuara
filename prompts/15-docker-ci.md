# Prompt 15 — Docker & CI

## Goal
Containerize OpenMuara and set up CI.

## Acceptance Criteria
- [x] Multi-stage `Dockerfile`
- [x] `docker-compose.yml` with volume mount for `.muara/`
- [x] GitHub Actions workflow:
  - lint (`golangci-lint`)
  - test (`go test ./...`)
  - build
  - docker build (optional)
- [x] `.dockerignore`
- [ ] CI badge in README

## Files Changed
- `Dockerfile`
- `docker-compose.yml`
- `.dockerignore`
- `.github/workflows/ci.yml`
- `README.md`

## Response Shape
Return:
1. Dockerfile stages
2. Compose service definition
3. CI job list

## Test Notes
- `docker build -t openmuara .`
- `docker compose up --build`
- Push to branch and verify CI run
