# 017 — No package-level documentation (`go doc` / pkg.go.dev shows nothing)

**Problem:** `package rptgen` has no package doc comment, so `go doc github.com/lfrystak/rptgen` and the pkg.go.dev landing page have no overview, despite a good README that the godoc tooling does not surface.

## Background / context

Every `.go` file starts directly with `package rptgen` and no preceding doc comment (`rptgen.go:1`, `html.go:1`, etc.). There is no `doc.go`. Exported types/methods have decent inline comments, but there is:

- No package synopsis (the one-liner shown in `go doc` and pkg.go.dev search results).
- No package overview / quick-start in godoc (README content is not visible there).
- No runnable `Example` functions, which pkg.go.dev renders and runs as compiled documentation.

For a library whose primary discovery surface is pkg.go.dev, this materially hurts usability and perceived quality.

## Severity

**Nice-to-have** — does not affect behaviour, but meaningfully improves adoption and is low effort.

## Proposed change / acceptance criteria

1. Add a `doc.go` with a package comment: one-sentence synopsis followed by an overview and a minimal usage snippet (mirroring, not duplicating, the README).
2. Add at least one runnable `Example` (e.g. `ExampleHtmlRenderer_Render`) in a `_test.go` file so it appears and is compiled on pkg.go.dev.
3. `go doc github.com/lfrystak/rptgen` shows a useful synopsis + overview; `go test` runs the example.
4. Ensure exported symbols changed by other specs keep accurate doc comments.

## Dependencies

- Write after API-shaping specs settle ([003](003-chart-data-ordering.md), [006](006-render-to-io-writer.md), [008](008-api-consistency.md)) so the example/synopsis reflect the final API; coordinate wording with [013-documentation-drift.md](013-documentation-drift.md).
