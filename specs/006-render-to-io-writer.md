# 006 — `Render` returns a fully-buffered `string` instead of writing to `io.Writer`

**Problem:** The renderer materializes the entire HTML document (CSS + embedded Chart.js + all data) as a single in-memory `string`, with no streaming option — un-idiomatic for Go document/output generation and memory-heavy for large reports.

## Background / context

`Renderer.Render(report *Report, theme *Theme) (string, error)` (`rptgen.go:104`, impl `html.go:21`) builds everything in a `strings.Builder` and returns `b.String()`. Every caller then does its own I/O — `example/main.go:201` and `README.md:47` both immediately `os.WriteFile([]byte(html), …)`, round-tripping the whole document string→[]byte.

Idiomatic Go for emitting a document is to accept an `io.Writer` (cf. `template.Execute`, `json.NewEncoder`, `http.ResponseWriter`). Current design:

- Forces full buffering; a report with many large tables holds the entire HTML plus the embedded Chart.js bundle (see [001](001-charts-not-self-contained.md)) in memory twice (string + `[]byte` at the call site).
- Cannot stream to an `http.ResponseWriter` or a file without an extra copy.
- No place to thread cancellation (`context.Context`) for very large generations.

## Severity

**Important** — API ergonomics/idiom and memory behaviour for the library's primary operation.

## Proposed change / acceptance criteria

1. Make the primary renderer method write to an `io.Writer`, e.g. `Render(w io.Writer, report *Report, theme *Theme) error` (optionally `RenderContext(ctx, w, …)`), updating the `Renderer` interface accordingly.
2. Provide a thin string convenience helper (e.g. `RenderString(report, theme) (string, error)`) so simple callers and existing examples stay one-liners.
3. Internally write directly to the `io.Writer` (wrap in `bufio.Writer`); keep `strings.Builder` only inside the convenience wrapper.
4. Update `README.md` example, `example/main.go`, and tests to the writer form (e.g. render straight into the opened file).
5. Decide error semantics jointly with [004-error-handling.md](004-error-handling.md) (write errors are now real).

## Dependencies

- Signature change must be coordinated with [004-error-handling.md](004-error-handling.md) and [005-extensibility-element-interface.md](005-extensibility-element-interface.md) (same interface).
- Doc/example updates overlap [013-documentation-drift.md](013-documentation-drift.md).
