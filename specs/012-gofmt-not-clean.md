# 012 — Source is not `gofmt`-clean

**Problem:** `gofmt -l .` reports `html.go` and `rptgen.go` as unformatted; the canonical Go formatting baseline is not met and nothing enforces it.

## Background / context

Running `gofmt -l .` in the repo root outputs:

```
html.go
rptgen.go
```

Likely culprits are misaligned struct field/tag columns, e.g. `rptgen.go:82-83`:

```go
ShadowIntensity string
EnableAnimations bool
EnableGradients  bool
```

and the chart-config structs in `html.go` (e.g. `html.go:314-318`, `Options chartOptions` tag alignment). Unformatted code in a Go library is an immediate credibility/maintainability signal and produces noisy diffs once contributors' editors auto-format.

## Severity

**Important** — baseline hygiene for a Go project; prerequisite for an enforcing CI gate.

## Proposed change / acceptance criteria

1. Run `gofmt -w .` (or `go fmt ./...`) and commit the result; `gofmt -l .` must output nothing.
2. Land this **before/with** the CI gofmt gate in [009-ci-workflow-defects.md](009-ci-workflow-defects.md) so CI does not immediately fail.
3. No functional changes — formatting only; tests unaffected.

## Dependencies

- Prerequisite for the gofmt gate in [009-ci-workflow-defects.md](009-ci-workflow-defects.md).
- Will conflict-rebase with any other in-flight edits to `html.go`/`rptgen.go` ([001](001-charts-not-self-contained.md), [015](015-html-monolith-refactor.md)); sequence accordingly.
