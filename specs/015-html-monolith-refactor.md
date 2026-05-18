# 015 — `html.go` is a 765-line monolith mixing five unrelated concerns

**Problem:** A single file hand-builds HTML, generates CSS, defines Chart.js JSON config structs, dispatches element types, and generates IDs — all by raw string concatenation — making the renderer hard to maintain, test, and extend.

## Background / context

`html.go` (765 lines) contains, intermixed:

1. Document assembly (`Render`, `html.go:21`).
2. Element type-switch dispatch + per-element HTML builders (`renderElement`…`renderChartContainer`, `html.go:151-283`).
3. Chart.js config model (`chartConfig` and friends, `html.go:312-358`) + per-chart builders (`html.go:360-524`).
4. A ~225-line `generateCSS` (`html.go:527-752`) that `fmt.Fprintf`s a stylesheet as Go string literals with `%%`-escaped percent signs.
5. ID generation + slug logic (`slugify`, `idGen`, `html.go:107-148`).

All HTML/CSS is produced via `strings.Builder` + `fmt.Fprintf` rather than `html/template`/`text/template`, which is the root cause of the escaping risk in [002-script-context-injection.md](002-script-context-injection.md) and makes structural changes error-prone (indentation strings hardcoded, `%%` littered through CSS). The coupling also blocks a second renderer ([005-extensibility-element-interface.md](005-extensibility-element-interface.md)) and makes targeted unit testing impossible (everything funnels through the whole-document `Render`).

## Severity

**Important** — primary maintainability/extensibility bottleneck; enabler for [002](002-script-context-injection.md) and [005](005-extensibility-element-interface.md).

## Proposed change / acceptance criteria

1. Split into focused files within the package (or an `internal/` subpackage), e.g.: `render_html.go` (document assembly), `render_elements.go` (per-element), `render_charts.go` (chart builders), `chartjs.go` (config model), `css.go` (theme→CSS), `idgen.go` (slug/ID).
2. **The Chart.js config model (`chartjs.go`) must be open for extension.** The current `chartConfig`/`chartDataset`/`chartOptions` structs (`html.go:314-358`) are a closed, hand-rolled subset of Chart.js; new chart types (scatter `{x,y}` points, radar `r` scale, bubble, area, mixed) need fields it lacks. The refactored model must let a new chart type supply its dataset/options shape **without editing the shared structs** — e.g. per-chart-type config structs, and/or a typed passthrough for arbitrary Chart.js options (a `map[string]any`/`json.RawMessage` "extra" escape hatch that is still context-safe per [002](002-script-context-injection.md)). Acceptance: adding the `ScatterChart` from [005-extensibility-element-interface.md](005-extensibility-element-interface.md) requires no modification to the shared `chartConfig`/`chartDataset`/`chartOptions` structs.
3. Move static CSS into an embedded `.css` template (`//go:embed`) with theme values injected via a typed struct, eliminating `%%`-escaping and Go-string CSS.
4. Use `html/template` / `text/template` with context-correct escaping for HTML and the script block (directly resolves [002-script-context-injection.md](002-script-context-injection.md)).
5. Make per-element/per-chart renderers independently unit-testable (export within package or via small interfaces) so [016-test-quality.md](016-test-quality.md) can add focused tests.
6. No behavioural change beyond the escaping fix; existing tests (updated for [001](001-charts-not-self-contained.md)) stay green; golden-file tests added per [016](016-test-quality.md).

## Dependencies

- Enables/co-lands with [002-script-context-injection.md](002-script-context-injection.md) and [005-extensibility-element-interface.md](005-extensibility-element-interface.md).
- Touches the same files as [001](001-charts-not-self-contained.md) and [012](012-gofmt-not-clean.md) — sequence to avoid churn.
- Golden tests in [016-test-quality.md](016-test-quality.md) should land first or together to catch regressions during the split.
