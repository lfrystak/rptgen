# 009 — CI workflow has a placeholder step, non-enforcing lint, and missing quality gates

**Problem:** `.github/workflows/ci.yml` contains a dead placeholder step, runs the linter in non-failing mode, and omits `go vet`, the race detector, and a `gofmt` gate — so CI is largely decorative.

## Background / context

`.github/workflows/ci.yml`:

- `ci.yml:48` — a step **named `Test`** whose `run:` is `ls -al`. It does nothing and is confusingly named identically to the real test step at `ci.yml:54`.
- `ci.yml:42` — `golangci-lint` has `continue-on-error: true`, so lint findings never fail the build. Combined with results only flowing into SonarQube, lint is advisory at best.
- No `go vet ./...` step.
- `go test` (`ci.yml:55`) runs without `-race`.
- No `gofmt`/`goimports` enforcement — and the repo is currently **not gofmt-clean** (`gofmt -l` flags `html.go` and `rptgen.go`; see [012-gofmt-not-clean.md](012-gofmt-not-clean.md)). CI would not catch this.
- `on.push.branches` includes `develop` (`ci.yml:7`); confirm that branch exists / is intended.

A green check currently guarantees only "compiles and smoke tests pass," not "vetted, formatted, race-clean, lint-clean."

## Severity

**Important** — CI is the maintainability backbone; right now it gives false assurance.

## Proposed change / acceptance criteria

1. Delete the placeholder `ls -al` step (`ci.yml:48-49`).
2. Add a formatting gate that fails on unformatted code: `test -z "$(gofmt -l .)"` (or `gofmt -l . && exit on output`). Coordinate with [012-gofmt-not-clean.md](012-gofmt-not-clean.md) so the repo is clean first.
3. Add `go vet ./...` as a failing step.
4. Run tests with `-race` (`go test -race ./... -coverprofile=coverage.out`).
5. Make lint enforcing: remove `continue-on-error: true` (optionally keep SARIF upload) so lint failures break CI; fix or explicitly `//nolint` existing findings.
6. Confirm/trim the `develop` branch trigger.
7. CI passes on a clean tree with all gates enforcing.

## Dependencies

- Must land with or after [012-gofmt-not-clean.md](012-gofmt-not-clean.md) (else the new gofmt gate fails immediately).
- Lint-enforcing may surface issues related to [004](004-error-handling.md)/[007](007-baseelement-dead-boilerplate.md).
- Release job correctness tracked separately in [010-goreleaser-library-misconfig.md](010-goreleaser-library-misconfig.md).
