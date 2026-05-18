# 010 — GoReleaser is configured to build binaries for a library

**Problem:** `rptgen` is an importable library with no command, yet `.goreleaser.yaml` builds and ships OS/arch binaries — which, because the only `main` package is the demo, would release the example report generator as the project's "binary."

## Background / context

- The project is a library (`README.md:12`, package `rptgen`, no `cmd/`).
- `.goreleaser.yaml:7-17` defines `builds:` with `main: ./...`, `goos: [linux, windows, darwin]`, `goarch: [arm64, amd64]`.
- The only `package main` in the module is `example/main.go` (a demo that hardcodes sample data and writes `report.html`). `main: ./...` will resolve to it, so `goreleaser release` produces cross-platform archives of the **example program**, presented on the Releases page as if it were the product.
- The release job runs on tags (`ci.yml:75-100`), so this ships on every `v*` tag.

For a library, releases are git tags + `go get`; binary artifacts are misleading and imply a CLI that does not exist (the README's "Latest Release" badge points here).

## Severity

**Important** — public release artifacts misrepresent the project; wrong release model for a library.

## Proposed change / acceptance criteria

1. Decide the intended distribution model:
   - **Library only (likely):** Remove the `builds:` binary section; use GoReleaser only for source archives / changelog / GitHub release notes, or drop GoReleaser and rely on tags + `go.mod` versioning. Keep the changelog/release-notes config.
   - **Library + real CLI (if desired):** Create a proper `cmd/rptgen` (not `example/`), restrict `builds.main` to it, and document the CLI.
2. Ensure `example/` is never the released binary (exclude it explicitly or move it under `examples/` / `internal`).
3. Verify a dry run (`goreleaser release --snapshot --clean`) produces the intended artifacts only.
4. README "Latest Release" badge/expectations reconciled with the chosen model ([013-documentation-drift.md](013-documentation-drift.md)).

## Dependencies

- Documentation reconciliation in [013-documentation-drift.md](013-documentation-drift.md).
- Related CI changes in [009-ci-workflow-defects.md](009-ci-workflow-defects.md) (same workflow file, release job).
