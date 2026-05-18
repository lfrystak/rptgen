# 025 — No scatter chart support; no `{x,y}` point data shape

**Problem:** Chart.js `type: "scatter"` plots `{x, y}` numeric points, but every `rptgen` chart and the internal `chartDataset` model assume category-labelled `[]float64` data (`html.go:325-333`). There is no XY point shape, so scatter (and later bubble) cannot be expressed at all.

## Background / context

Existing charts map a `labels []string` array to parallel `data []float64`. Scatter datasets instead carry `data: [{x, y}, ...]` with a numeric `scales.x` (linear, not category). `chartDataset.Data` is typed `[]float64` (`html.go:328`), which cannot represent points. This spec introduces the shared **XY point family** that [026-bubble-chart.md](026-bubble-chart.md) and [027-mixed-chart.md](027-mixed-chart.md) also reuse, so the new shape is designed once here.

## Severity

**Nice-to-have** — standalone feature, but it establishes the shared XY data shape, so its design matters for 026/027.

## Proposed change / acceptance criteria

1. Add one exported point type, e.g. `type XYPoint struct{ X, Y float64 }`, and a `ScatterSeries struct{ Name string; Points []XYPoint }` (ordered slice — consistent with the `[]LineSeries`/`[]StackedBarSeries` convention).
2. Add `ScatterChart` with `Series []ScatterSeries` and `NewScatterChart(title string, series []ScatterSeries) *ScatterChart`.
3. Extend the internal dataset model to carry point data **without breaking the `[]float64` path**: either a separate point-dataset struct or a `Data any` that holds `[]float64` or `[]XYPoint`. Existing charts must serialize byte-identically.
4. Emit `type: "scatter"` with a linear `scales.x` (and `scales.y`); per-series theme colors; legend for >1 series.
5. Tests: points serialize as `{"x":..,"y":..}`; existing category charts unchanged; multi-series scatter → one dataset per series.
6. README/example: one scatter example.

## Dependencies

- **Coordination, not conflict, with [005-extensibility-element-interface.md](005-extensibility-element-interface.md) and [015-html-monolith-refactor.md](015-html-monolith-refactor.md):** both name `ScatterChart` as their *acceptance vehicle* ("a new chart type can be added without modifying `renderElement`/the shared config structs"). Two valid orderings, pick one and note it on scheduling: (a) land 005/015 first, then this spec implements `ScatterChart` purely through the open seam and *is* their acceptance test; or (b) land this first against the current architecture (closed switch + additive config field), and 005/015 later refactor it to prove the seam. Either is conflict-free; the specs must not both claim to introduce `ScatterChart` independently — flag the dependency at scheduling time.
- Establishes the XY shape consumed by [026-bubble-chart.md](026-bubble-chart.md) and [027-mixed-chart.md](027-mixed-chart.md); schedule 025 before 026.
- Common options via [020-shared-chart-options.md](020-shared-chart-options.md) when available.
