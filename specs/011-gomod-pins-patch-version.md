# 011 — `go.mod` pins a patch toolchain version (`go 1.26.3`) for a library

**Problem:** `go.mod` declares `go 1.26.3`, forcing every consumer onto Go ≥ 1.26.3 even though nothing in the code requires it, which is needlessly restrictive for an importable library.

## Background / context

`go.mod`:

```
module github.com/lfrystak/rptgen
go 1.26.3
```

Since Go 1.21 the `go` directive is a real minimum-version requirement: a module importing `rptgen` with a toolchain older than 1.26.3 will be told to upgrade (or auto-download a newer toolchain). The code uses only long-stable standard-library APIs (`strings.Builder`, `encoding/json`, `html`, `sort`, `time`, `fmt`) — there is no 1.26-specific feature in use. CI uses `GO_VERSION: stable` (`ci.yml:20`), so the pin is not even matching a deliberate floor; it is just whatever was installed when `go mod init` ran.

Libraries should declare the **lowest** Go version they actually support to maximize consumer compatibility; pinning a patch release is an anti-pattern.

## Severity

**Important** — unnecessarily narrows the library's usable audience; trivial to fix.

## Proposed change / acceptance criteria

1. Lower the `go` directive to the minimum version actually required (e.g. `go 1.21`, or the project's chosen support floor), without a patch suffix.
2. Verify `go build ./...` and `go test ./...` pass on that declared minimum (add it to the CI matrix or at least test locally).
3. Optionally add a CI job that builds against the declared minimum Go version in addition to `stable` to prevent regressions.

## Dependencies

- None blocking. CI matrix addition coordinates with [009-ci-workflow-defects.md](009-ci-workflow-defects.md).
