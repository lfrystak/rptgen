# 021 — No area chart support

**Problem:** Chart.js renders an area chart as a line dataset with `fill` enabled, but `LineChart` hardcodes `Fill: false` (`html.go:408,419-424`), so there is no way to produce a filled/area chart — a standard, commonly-requested chart type.

## Background / context

In Chart.js an "area chart" is not a distinct `type`; it is `type: "line"` with `dataset.fill` set (e.g. `fill: true` or `fill: 'origin'`) and typically a semi-transparent `backgroundColor`. `rptgen`'s `LineChart` already builds line datasets but unconditionally sets `falseVal := false; ds.Fill = &falseVal` (`html.go:408,422`), foreclosing the area variant. The data shape (`[]LineSeries` with ordered points) is already exactly what an area chart needs — only the `fill` flag and a translucent fill color differ.

## Severity

**Nice-to-have** — discrete, low-risk feature; small surface area because it reuses the entire `LineChart` path.

## Proposed change / acceptance criteria

1. Add an opt-in fill on `LineChart` — a single boolean (e.g. `LineChart.Area bool`, default `false` preserving today's output) plus an `NewAreaChart(title string, series []LineSeries) *LineChart` convenience constructor that sets it. Do **not** add a new element type or `ElementType()` value — it stays a `LineChart` so the dispatch and data shape are unchanged.
2. When `Area` is true: emit `fill: true` (single series) / stacked-appropriate fill, and derive a translucent `backgroundColor` from the series color (the existing per-series theme color at reduced alpha) while keeping the opaque `borderColor`.
3. When `Area` is false: byte-identical output to today.
4. Respect common options from [020-shared-chart-options.md](020-shared-chart-options.md) if present; otherwise current defaults.
5. Tests: `Area: true` yields `fill` truthy and a translucent background; `Area: false` unchanged; multi-series area renders one filled dataset per series.
6. README/example: one short area-chart example.

## Dependencies

- Independently completable against the current architecture (extend `renderLineChartScript`). No dependency on 001–019.
- If [005-extensibility-element-interface.md](005-extensibility-element-interface.md)/[015-html-monolith-refactor.md](015-html-monolith-refactor.md) have landed, implement via the open render/config seam instead of editing the central switch — behaviour identical either way.
- Picks up axis titles/legend from [020-shared-chart-options.md](020-shared-chart-options.md) when available; not blocked by it.
