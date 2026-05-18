# 023 — No radar chart support

**Problem:** Chart.js supports `type: "radar"` for comparing multiple series across shared quantitative axes, but `rptgen` has no radar chart and the internal config model has no radial (`scales.r`) support.

## Background / context

A radar chart is a category/multi-series chart: shared axis labels (the spokes) and one dataset per series — the *same logical shape* as `LineChart` (`[]LineSeries`, ordered points keyed by spoke label). It differs only in `type: "radar"`, a radial `scales.r` (not `x`/`y`), and that datasets are typically filled. The current `chartScales` struct only models `X`/`Y` (`html.go:351-358`), so radial scale config is a model gap, not just a missing type.

## Severity

**Nice-to-have** — distinct standard chart type; reuses the line-style multi-series data shape.

## Proposed change / acceptance criteria

1. Add a `RadarChart` type with `Series []LineSeries` (reuse the existing ordered multi-series shape — do not invent a parallel one) and `NewRadarChart(title string, series []LineSeries) *RadarChart`.
2. Emit `type: "radar"`; build one dataset per series with per-series theme colors and a translucent fill (radar is conventionally filled).
3. Extend the internal scales model additively with an optional radial axis (`scales.r`: at minimum `beginAtZero`, and `min`/`max`/title from [020](020-shared-chart-options.md) when present). Keep additions `omitempty` so non-radar charts are byte-identical.
4. All series must share the spoke label set; collect labels in deterministic first-appearance/sorted order exactly as `renderLineChartScript` already does (`html.go:396-406`).
5. Legend shown for >1 series.
6. Tests: 3 series over 5 spokes → 3 datasets, 5 labels, `type:"radar"`, radial scale present; non-radar charts unchanged.
7. README/example: one radar example.

## Dependencies

- Independently completable against current architecture (new case + additive `scales.r` field). No dependency on 001–019.
- Radial-scale field is the kind of model extension [015-html-monolith-refactor.md](015-html-monolith-refactor.md) requires the config model to absorb "without editing shared structs"; if 015 landed, add it as a per-type/extension config. Compatible either way; not blocked.
- Common options (axis/legend) via [020-shared-chart-options.md](020-shared-chart-options.md) when available.
