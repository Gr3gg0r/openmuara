> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Repository Hygiene Appendix

> **Created:** 2026-07-10  
> **Status:** Draft

This appendix contains exact templates and checklists for the repo-hygiene initiative. During implementation, copy these into the target files unless a better version emerges.

---

## A. Secret-scan checklist

```bash
# Full-history scan with gitleaks
gitleaks detect --source . --verbose

# Inspect tracked files for obvious secrets or local configs
git ls-files | grep -E '\.(env|key|pem|p12|pfx)$'
git ls-files | grep -E 'config\.yml|secrets|credentials'

# Look for large or binary files
git ls-files | while read -r f; do
  size=$(stat -f%z "$f" 2>/dev/null || echo 0)
  if [ "$size" -gt 1048576 ]; then
    echo "$size $f"
  fi
done
```

Expected result: zero leaks, no `.env`/`.key` files, no binaries over 1 MiB.

---

## B. `.editorconfig` (root)

```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

[*.go]
indent_style = tab

[*.{yml,yaml}]
indent_size = 2

[*.{md,mdx}]
trim_trailing_whitespace = false
max_line_length = 120

[*.{ts,tsx,js,jsx}]
indent_size = 2

[*.{css,scss}]
indent_size = 2

[*.sh]
indent_size = 2
end_of_line = lf

[Makefile]
indent_style = tab
```

---

## C. `.gitattributes` (root)

```gitattributes
# Auto-detect text files and normalize line endings
* text=auto

# Shell scripts must use LF even on Windows
*.sh text eol=lf
*.bash text eol=lf

# Go source
go.mod text
go.sum text
*.go text

# Generated / vendored files — hide from GitHub diff and language stats
web/dashboard/node_modules/** linguist-generated
website/node_modules/** linguist-generated
internal/ui/dashboard-dist/assets/** linguist-generated
internal/ui/dashboard-dist/*.gz linguist-generated
web/dashboard/coverage/** linguist-generated
coverage.out linguist-generated
coverage.html linguist-generated
```

---

## D. Proposed `.gitignore` updates

Add or ensure the following entries exist (remove duplicates):

```gitignore
# Agent workspace — local planning artifacts, not committed
# .toyol/ is legacy; kept until local workspace is migrated.
.agents/
.toyol/
.gstack/
.playwright-mcp/

# Runtime artifacts — generated during local runs, do not commit
queue-logs/
screenshots/
qa/
state/
*.log

# Build artifacts
bin/
coverage.out
coverage.html

# Generated screenshots / QA artifacts
dashboard-mobile-*.png

# Environment and secrets
.env
.env.*
!.env.example

# OS and editor files
.DS_Store
*.swp
*.swo
.vscode/settings.json
.idea/

# Dependencies and build outputs
node_modules/
dist/
build/
*.tmp
internal/ui/dashboard-dist/assets/
internal/ui/dashboard-dist/*.gz

# Playwright outputs
test-results/
playwright-report/
```

---

## E. Proposed `.dockerignore` updates

```dockerignore
.git/
.muara/
.toyol/
.agents/
.gstack/
.playwright-mcp/
bin/
coverage.out
coverage.html
*.log
.DS_Store
.idea/
.vscode/
node_modules/
dist/
build/
*.tmp
test-results/
playwright-report/
web/dashboard/coverage/
internal/ui/dashboard-dist/assets/
internal/ui/dashboard-dist/*.gz
```

---

## F. Root `SECURITY.md` redirect

Replace the current root `SECURITY.md` with:

```markdown
# Security Policy

The canonical security policy for OpenMuara lives in
[`.github/SECURITY.md`](.github/SECURITY.md).

Please report vulnerabilities privately as described there.
```

---

## G. `.github/SUPPORT.md`

```markdown
# Getting Support

Thanks for using OpenMuara. This page lists the best ways to get help, report
issues, or stay informed.

## Documentation

- [Quick start](docs/quickstart.md)
- [Local development runbook](runbooks/local-development.md)
- [Quality gates](runbooks/quality-gates.md)
- [Full docs site](https://openmuara.github.io/openmuara/)

## Ask a question

For questions, ideas, or general discussion, use
[GitHub Discussions](https://github.com/openmuara/openmuara/discussions).

## Report a bug

Open a [bug report](https://github.com/openmuara/openmuara/issues/new?template=bug_report.yml)
and include:

- Steps to reproduce.
- The output of `muara doctor --json`.
- The failing quality gate, if any.

## Request a feature or provider

- [Feature request](https://github.com/openmuara/openmuara/issues/new?template=feature_request.yml)
- [Provider emulation request](https://github.com/openmuara/openmuara/issues/new?template=provider_request.yml)

## Report a security issue

Please **do not** open a public issue. See [SECURITY.md](SECURITY.md) for
private reporting instructions.

## Code of Conduct

All interactions are governed by our [Code of Conduct](../CODE_OF_CONDUCT.md).
```

---

## H. `MAINTAINERS.md`

```markdown
# Maintainers

This file lists the current maintainers of OpenMuara.

| Name | GitHub | Areas |
|---|---|---|
| Shahfiq | @shahfiq | Project lead, architecture, provider contracts |

## Emeritus

None yet.

## How to become a maintainer

Consistent, high-quality contributions over time are the path to maintainer
status. Maintainers are added by consensus of existing maintainers. See
[CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.
```

> **Note:** Replace `@shahfiq` with the actual maintainer GitHub handle before committing.

---

## I. `.github/FUNDING.yml`

```yaml
# OpenMuara does not currently accept sponsorships.
# This file is a placeholder so the sponsorship section is intentionally empty.
# When a funding model is established, uncomment and configure the entries below.
#
# github: [openmuara]
# open_collective: openmuara
# custom: ["https://example.com/donate"]
```

---

## J. `.github/ISSUE_TEMPLATE/provider_request.yml`

```yaml
name: Provider emulation request
description: Suggest a new payment provider to emulate in OpenMuara
title: '[provider] '
labels: ['provider', 'enhancement']
body:
  - type: markdown
    attributes:
      value: |
        Thanks for suggesting a provider. Please include links to the real
        provider's API docs and the endpoints you need to emulate.
  - type: input
    id: provider_name
    attributes:
      label: Provider name
      placeholder: e.g., PayPal, Midtrans, Razorpay
    validations:
      required: true
  - type: input
    id: docs_url
    attributes:
      label: API documentation URL
    validations:
      required: true
  - type: textarea
    id: endpoints
    attributes:
      label: Endpoints to emulate
      placeholder: |
        - POST /v1/charges
        - POST /v1/webhooks
    validations:
      required: true
  - type: textarea
    id: signature
    attributes:
      label: Signature or authentication scheme
      description: How does the provider sign or authenticate requests?
    validations:
      required: false
  - type: textarea
    id: use_case
    attributes:
      label: Use case
      description: What are you testing that requires this provider?
    validations:
      required: true
```

---

## K. `.github/ISSUE_TEMPLATE/docs_issue.yml`

```yaml
name: Documentation issue
description: Report a docs error, gap, or unclear explanation
title: '[docs] '
labels: ['docs']
body:
  - type: input
    id: page
    attributes:
      label: Page or file
      placeholder: e.g., docs/quickstart.md
    validations:
      required: true
  - type: textarea
    id: problem
    attributes:
      label: What is wrong or missing?
    validations:
      required: true
  - type: textarea
    id: suggestion
    attributes:
      label: Suggested improvement
    validations:
      required: false
```

---

## L. `.github/ISSUE_TEMPLATE/config.yml`

```yaml
blank_issues_enabled: false
contact_links:
  - name: Ask a question
    url: https://github.com/openmuara/openmuara/discussions
    about: Use GitHub Discussions for general questions and ideas.
  - name: Report a security issue
    url: https://github.com/openmuara/openmuara/security/advisories/new
    about: Please report security vulnerabilities privately.
```

---

## M. `.github/release.yml`

```yaml
changelog:
  exclude:
    labels:
      - ignore-for-release
  categories:
    - title: Breaking Changes
      labels:
        - breaking change
    - title: New Providers
      labels:
        - provider
    - title: Security
      labels:
        - security
    - title: Added
      labels:
        - enhancement
    - title: Fixed
      labels:
        - bug
    - title: Documentation
      labels:
        - docs
    - title: Maintenance
      labels:
        - ci/cd
        - tech debt
```

---

## N. GitHub label taxonomy

Create these labels in the repository (hex colors included):

| Label | Color | Description |
|---|---|---|
| `bug` | `#d73a4a` | Something is broken |
| `enhancement` | `#a2eeef` | New feature or improvement |
| `provider` | `#7057ff` | Provider emulation work |
| `docs` | `#0075ca` | Documentation only |
| `security` | `#d93f0b` | Security-related |
| `good first issue` | `#7057ff` | Friendly for newcomers |
| `help wanted` | `#008672` | Maintainer wants community help |
| `breaking change` | `#b60205` | Changes provider contract or CLI |
| `tech debt` | `#cccccc` | Cleanup or refactor |
| `ci/cd` | `#f9d0c4` | Build/release pipeline |
| `needs triage` | `#ededed` | Not yet reviewed by a maintainer |
| `wontfix` | `#ffffff` | Intentionally not fixing |
| `ignore-for-release` | `#ffffff` | Omit from generated release notes |

---

## O. Repository settings checklist

Apply these settings in the GitHub repository UI before or immediately after transfer:

### General

- [ ] **Wiki:** Disable if unused.
- [ ] **Projects:** Disable if unused.
- [ ] **Discussions:** Enable (linked from `SUPPORT.md`).
- [ ] **Sponsorships:** Disable until `FUNDING.yml` is populated.
- [ ] **Preserve this repository:** Optional.

### Branch protection — `main`

- [ ] Require a pull request before merging.
- [ ] Require approvals: at least 1.
- [ ] Dismiss stale PR approvals when new commits are pushed.
- [ ] Require status checks to pass before merging:
  - `docs`
  - `ui-build`
  - `ui-test`
  - `lint`
  - `unit`
  - `smoke`
  - `vuln`
  - `gosec`
  - `secrets`
  - `quality`
  - `dependency-license`
  - `docker-build`
  - `install-dry-run`
  - `changelog-check`
- [ ] Require signed commits.
- [ ] Include administrators.
- [ ] Restrict pushes that create files larger than 100 MB.

### Branch protection — `dev`

- [ ] Require a pull request before merging.
- [ ] Require approvals: at least 1.
- [ ] Require status checks to pass before merging (same list as `main`, or a subset).
- [ ] Do not require signed commits (keeps contribution barrier lower), but encourage them.

### Tags

- [ ] Restrict tag creation to maintainers.

---

## P. AI-assisted development disclosure

Add this short paragraph to `README.md` (near the License or Contributing section) and to `CONTRIBUTING.md`:

```markdown
## AI-assisted development

OpenMuara's code, documentation, and runbooks are developed with the assistance
of AI coding agents and reviewed by human maintainers. We treat AI-generated
output as a draft: it is tested, linted, and validated against the same quality
gates as human-written code before it is merged.
```

---

## Q. Conventional Commits quick reference

```text
<type>(<scope>): <short summary>

<body>

<footer>
```

Common types:

| Type | Use |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, no logic change |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or correcting tests |
| `chore` | Build/tooling/config changes |
| `ci` | CI/CD changes |
| `security` | Security fix |

Examples:

```text
feat(stripe): add PaymentIntent confirm simulation
fix(fawry): validate merchant_security_key length
docs(readme): update Docker run example
test(webhook): add chaos retry-exhaustion case
chore(hygiene): add .editorconfig
```

---

## R. Worktree cleanup commands

```bash
# List worktrees
git worktree list

# Remove a merged worktree (safe; does not delete the branch)
git worktree remove /path/to/worktree

# Delete a stale branch after its worktree is removed
git branch -d feat/old-branch

# Delete a suspended branch that was never merged (use with care)
git branch -D feat/suspended-branch
```

---

## S. Useful one-liners

```bash
# Find all files containing the legacy brand
grep -Ril --exclude-dir=.git --exclude-dir=node_modules toyol .

# Show tracked files that should probably be ignored
git ls-files | grep -E 'coverage\.out|coverage\.html|bin/|node_modules/|\.muara/|\.toyol/'

# Count commits per author
git shortlog -sn

# List branches merged into dev but not main
git branch --merged dev | grep -v '^\*' | grep -v main
```
