> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Production Dependency License Matrix

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Verified

---

This matrix lists every production dependency of OpenMuara, its version, SPDX license identifier, and compatibility assessment against the project's MIT License. It is generated from `go.mod`, `web/dashboard/package.json`, and `website/package.json`.

## License compatibility legend

| Status | Meaning |
|---|---|
| ✅ Compatible | Permissive license compatible with MIT distribution |
| ⚠️ Conditional | Copyleft or weak copyleft; may require source disclosure or have linking conditions |
| ❌ Incompatible | Strong copyleft (GPL/AGPL) or proprietary; must be replaced or explicitly accepted |
| ❓ Unknown | License not identified; must be resolved before release |

## Go production dependencies

Direct dependencies are listed first; indirect dependencies follow.

| Package | Version | License | Status | Notes |
|---|---|---|---|---|
| `github.com/google/uuid` | `v1.6.0` | BSD-3-Clause | ✅ Compatible | — |
| `github.com/mattn/go-isatty` | `v0.0.22` | MIT | ✅ Compatible | Promoted to direct by `go mod tidy` |
| `github.com/prometheus/client_golang` | `v1.23.2` | Apache-2.0 | ✅ Compatible | — |
| `github.com/spf13/cobra` | `v1.8.1` | Apache-2.0 | ✅ Compatible | — |
| `github.com/spf13/viper` | `v1.19.0` | MIT | ✅ Compatible | — |
| `golang.org/x/crypto` | `v0.53.0` | BSD-3-Clause | ✅ Compatible | Promoted to direct by `go mod tidy` |
| `gopkg.in/yaml.v3` | `v3.0.1` | MIT | ✅ Compatible | — |
| `modernc.org/sqlite` | `v1.53.0` | BSD-3-Clause | ✅ Compatible | — |
| `github.com/beorn7/perks` | `v1.0.1` | MIT | ✅ Compatible | Indirect |
| `github.com/cespare/xxhash/v2` | `v2.3.0` | MIT | ✅ Compatible | Indirect |
| `github.com/dustin/go-humanize` | `v1.0.1` | MIT | ✅ Compatible | Indirect |
| `github.com/fsnotify/fsnotify` | `v1.7.0` | BSD-3-Clause | ✅ Compatible | Indirect |
| `github.com/hashicorp/hcl` | `v1.0.0` | MPL-2.0 | ⚠️ Conditional | Weak copyleft; used unmodified via Viper; see D06 |
| `github.com/inconshreveable/mousetrap` | `v1.1.0` | Apache-2.0 | ✅ Compatible | Indirect |
| `github.com/magiconair/properties` | `v1.8.7` | BSD-2-Clause | ✅ Compatible | Indirect |
| `github.com/mitchellh/mapstructure` | `v1.5.0` | MIT | ✅ Compatible | Indirect |
| `github.com/munnerz/goautoneg` | `v0.0.0-20191010083416-a7dc8b61c822` | BSD-3-Clause | ✅ Compatible | Indirect |
| `github.com/ncruces/go-strftime` | `v1.0.0` | MIT | ✅ Compatible | Indirect |
| `github.com/pelletier/go-toml/v2` | `v2.2.2` | MIT | ✅ Compatible | Indirect |
| `github.com/prometheus/client_model` | `v0.6.2` | Apache-2.0 | ✅ Compatible | Indirect |
| `github.com/prometheus/common` | `v0.66.1` | Apache-2.0 | ✅ Compatible | Indirect |
| `github.com/prometheus/procfs` | `v0.16.1` | Apache-2.0 | ✅ Compatible | Indirect |
| `github.com/remyoudompheng/bigfft` | `v0.0.0-20230129092748-24d4a6f8daec` | BSD-3-Clause | ✅ Compatible | Indirect |
| `github.com/sagikazarmark/locafero` | `v0.4.0` | MIT | ✅ Compatible | Indirect |
| `github.com/sagikazarmark/slog-shim` | `v0.1.0` | BSD-3-Clause | ✅ Compatible | Indirect |
| `github.com/sourcegraph/conc` | `v0.3.0` | MIT | ✅ Compatible | Indirect |
| `github.com/spf13/afero` | `v1.11.0` | Apache-2.0 | ✅ Compatible | Indirect |
| `github.com/spf13/cast` | `v1.6.0` | MIT | ✅ Compatible | Indirect |
| `github.com/spf13/pflag` | `v1.0.5` | BSD-3-Clause | ✅ Compatible | Indirect |
| `github.com/subosito/gotenv` | `v1.6.0` | MIT | ✅ Compatible | Indirect |
| `go.uber.org/atomic` | `v1.9.0` | MIT | ✅ Compatible | Indirect |
| `go.uber.org/multierr` | `v1.9.0` | MIT | ✅ Compatible | Indirect |
| `go.yaml.in/yaml/v2` | `v2.4.2` | Apache-2.0 | ✅ Compatible | Indirect |
| `golang.org/x/exp` | `v0.0.0-20250305212735-054e65f0b394` | BSD-3-Clause | ✅ Compatible | Indirect |
| `golang.org/x/sys` | `v0.46.0` | BSD-3-Clause | ✅ Compatible | Indirect |
| `golang.org/x/text` | `v0.38.0` | BSD-3-Clause | ✅ Compatible | Indirect |
| `google.golang.org/protobuf` | `v1.36.8` | BSD-3-Clause | ✅ Compatible | Indirect |
| `gopkg.in/ini.v1` | `v1.67.0` | Apache-2.0 | ✅ Compatible | Indirect |
| `modernc.org/libc` | `v1.73.4` | MIT | ✅ Compatible | Indirect |
| `modernc.org/mathutil` | `v1.7.1` | BSD-3-Clause | ✅ Compatible | Indirect; go-licenses cannot auto-detect; manually verified |
| `modernc.org/memory` | `v1.11.0` | BSD-3-Clause | ✅ Compatible | Indirect |

## npm production dependencies — web/dashboard

| Package | Version | License | Status | Notes |
|---|---|---|---|---|
| `preact` | `10.29.7` | MIT | ✅ Compatible | — |

## npm production dependencies — website

| Package | Version | License | Status | Notes |
|---|---|---|---|---|
| `@docusaurus/core` | `3.10.1` | MIT | ✅ Compatible | — |
| `@docusaurus/faster` | `3.10.1` | MIT | ✅ Compatible | — |
| `@docusaurus/preset-classic` | `3.10.1` | MIT | ✅ Compatible | — |
| `@easyops-cn/docusaurus-search-local` | `0.55.2` | MIT | ✅ Compatible | — |
| `@mdx-js/react` | `3.1.1` | MIT | ✅ Compatible | — |
| `clsx` | `2.1.1` | MIT | ✅ Compatible | — |
| `prism-react-renderer` | `2.4.1` | MIT | ✅ Compatible | — |
| `react` | `19.2.7` | MIT | ✅ Compatible | — |
| `react-dom` | `19.2.7` | MIT | ✅ Compatible | — |

## npm transitive dependency summary

Transitive production dependencies in both npm packages were scanned with `npx license-checker --production`. All encountered licenses are permissive and compatible with MIT distribution:

- MIT
- Apache-2.0
- BSD-2-Clause
- BSD-3-Clause
- ISC
- MIT-0
- CC-BY-4.0 (build-time data only; not bundled into shipped artifacts)
- Python-2.0 (build-time tooling only; not bundled into shipped artifacts)

The `website` package carries accepted build-time vulnerabilities in Docusaurus transitive tooling (see `KNOWN_ISSUES.md` F05); these are security, not license, findings.

## Exceptions and decisions

| Package | Exception rationale | Decision ID |
|---|---|---|
| `github.com/hashicorp/hcl` | MPL-2.0 weak copyleft; used unmodified as a transitive dependency of Viper; compatible with MIT distribution when unmodified | D06 |
| `modernc.org/mathutil` | go-licenses cannot auto-detect license text; manually verified as BSD-3-Clause | D07 |

## Generation commands

```bash
# Go
./scripts/check-licenses.sh
go-licenses csv ./... > /tmp/licenses-go.csv

# npm (per package directory)
npx license-checker --direct --production --csv
npm sbom --package-lock-only --sbom-format=spdx
```

## Verification checklist

- [x] All direct Go production dependencies are listed.
- [x] All indirect Go production dependencies are listed.
- [x] All npm production dependencies are listed.
- [x] Every dependency has an SPDX license identifier.
- [x] No dependency has an ❌ Incompatible or ❓ Unknown status without a recorded decision.
- [x] This file is regenerated before every release.

*This file should be regenerated and reviewed before every release.*
