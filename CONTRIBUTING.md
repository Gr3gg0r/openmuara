# Contributing to OpenMuara

Thanks for helping make OpenMuara better. This document covers the basics for
code, tests, and commits.

## Development setup

1. Install Go 1.25+ and [Task](https://taskfile.dev/installation/).
2. Install Node.js 20+ if you will work on the embedded admin dashboard or Docusaurus website.
3. Clone the repo and switch to the `dev` branch.
4. Install optional but recommended tools:
   - `golangci-lint`
   - `govulncheck` (`go install golang.org/x/vuln/cmd/govulncheck@latest`)
   - `pre-commit` (`pip install pre-commit` or `brew install pre-commit`)

```bash
git checkout dev
pre-commit install   # optional but recommended
task check
```

To build the embedded web UI before running or testing:

```bash
task ui:build
```

## Branching

- `main` and `dev` are protected. Do not push directly to them.
- The default working branch is `dev`. Most work should branch from and merge back into `dev`.
- Use feature branches: `feat/<description>` or `fix/<description>`.
- Open a pull request from your feature branch to `dev`.
- `main` is updated from `dev` at release time (see `runbooks/release.md`).

## Quality gates

Before pushing, run:

```bash
task quality
```

This runs formatting, vet, lint, race tests, coverage gate, smoke test,
vulnerability scan, forbidden-pattern check, shell-script check, size advisory,
and tracker audit.

You can also run individual gates:

```bash
task check      # fmt + vet + lint + test + coverage
task smoke      # end-to-end smoke test
task vuln       # govulncheck
task forbidden  # no fmt.Println / os.Exit in library code
task scripts    # shellcheck
task sizes      # advisory size report
```

## Validating workflow changes locally

### With act

Install [act](https://nektosact.com/installation/index.html) and use the
provided `.actrc`:

```bash
act -j docker-build
act -j install-dry-run
act -j lint
```

To test the release workflow, create `.github/test-events/release.json`:

```json
{
  "ref": "refs/tags/v0.0.0-test.1",
  "ref_name": "v0.0.0-test.1"
}
```

Then run:

```bash
act -j release --eventpath .github/test-events/release.json
```

> Note: `act` cannot fully emulate OIDC-based cosign signing. Release signing
> tests must be done on a fork.

### On a fork

For workflow changes that affect releases or container signing:

1. Fork the repository.
2. Push your branch to your fork.
3. Tag a test release (e.g., `v0.0.0-test.1`) and push it.
4. Verify the workflow artifacts, signatures, and container image on your fork.
5. Delete the test tag and release when done.

## Writing code

- Follow the existing Go style: explicit types, small functions, table-driven
  tests.
- Keep files, functions, and lines within the recommended size limits. Run
  `task sizes` to see current advisory warnings.
- Do not add `fmt.Println` in production code or `os.Exit` outside `cmd/`.
- Do not commit real `.muara/config.yml` files or API secrets.

## Writing tests

- New features and bug fixes should include tests.
- Use the shared helpers in `internal/testutil` where they fit.
- Aim to cover error paths, not just happy paths.
- Run tests with the race detector (`task test`).

## Commits

- One logical change per commit.
- Follow [Conventional Commits](https://www.conventionalcommits.org/):

  ```text
  <type>(<scope>): <short summary>
  ```

  Common types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `security`.

  Examples:
  - `feat(stripe): add PaymentIntent confirm simulation`
  - `fix(fawry): validate merchant_security_key length`
  - `test(webhook): add chaos retry-exhaustion case`
  - `docs(readme): update Docker run example`
  - `chore(hygiene): add .editorconfig`

- Reference the bug register ID in the commit body or footer when fixing a tracked bug.

## Opening issues and pull requests

1. Check existing issues and pull requests before opening a new one.
2. For bugs, use the **Bug report** issue template and include:
   - Steps to reproduce.
   - The output of `muara doctor --json`.
   - The output of `task quality` or the failing gate.
3. For features, use the **Feature request** template and describe the use case.
4. Pull requests should target the `dev` branch and include the provided PR template.

Read [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) before participating.

## Provider contributions

When adding or changing a provider emulation:

- Keep the public provider routes contract-faithful to the real gateway.
- Add signature/validation tests for the provider.
- Update `docs/providers.md` and `docs/providers/<provider>.md` if applicable.
- Add a smoke-test step only if it exercises a unique flow.

## Documentation

Documentation lives in `docs/` and `runbooks/`. Public-facing changes (install,
security, provider setup) should update:

- `README.md`
- `docs/quickstart.md`
- `runbooks/local-development.md`

## AI-assisted development

OpenMuara's code, documentation, and runbooks are developed with the assistance
of AI coding agents and reviewed by human maintainers. We treat AI-generated
output as a draft: it is tested, linted, and validated against the same quality
gates as human-written code before it is merged.
