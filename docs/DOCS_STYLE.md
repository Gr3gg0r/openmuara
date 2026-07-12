---
id: docs-style
title: Documentation Style Guide
---

# Documentation Style Guide

This guide keeps OpenMuara docs consistent, accurate, and easy to maintain.

## Tone

- **Direct and concise.** Assume the reader is busy.
- **Persona-aware.** Separate developer, contributor, tester, and maintainer paths.
- **No marketing fluff.** State what OpenMuara does and how to use it.

## Code examples

- Use `127.0.0.1` instead of `localhost`.
- Use clearly fake secrets and keys:
  - `sk_test_muara`
  - `muara-fawry-secret`
  - `whsec_muara`
- Make examples copy-paste runnable where possible.
- Prefer `curl` for HTTP examples.
- Format JSON responses so readers can scan them.

## File conventions

- Front matter for Docusaurus:

  ```markdown
  ---
  id: page-id
  title: Page Title
  ---
  ```

- One H1 (`#`) per file.
- Use sentence case for headings.
- Keep lines under 120 characters where practical.

## Links

- Use relative links inside the repo: `[Contributing](/CONTRIBUTING.md)`.
- Do not link to internal consumer repos or private paths.
- Update `website/sidebars.ts` when adding or renaming docs.

## Provider docs

Every provider page should include:

1. Configuration snippet
2. First request
3. Signature algorithm
4. Simulation / escape routes
5. Webhook payload example
6. Common errors table
7. See also links

## Markdown lint

Run the linter before committing doc changes:

```bash
npx markdownlint-cli2 "docs/**/*.md" "runbooks/**/*.md"
```

## One source of truth

Never duplicate instructions. If the same content belongs in two places, keep
the canonical version in one file and link or redirect from the other.
