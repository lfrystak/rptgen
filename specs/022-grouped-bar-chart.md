# 022 — `BarChart` is single-series only; no grouped (multi-series) bar chart

**Problem:** `BarChart.Data` is a single `map[string]float64` (`charts.go:13`), so it can only render one series. Chart.js bar charts natively support multiple datasets rendered side-by-side (grouped bars) — a common need (e.g. revenue vs. costs per quarter, unstacked) that today has no representation. `StackedBarChart` covers the *stacked* case but not the *grouped* one.

## Background / context

`renderBarChartScript` (`html.go:360-391`) emits exactly one `chartDataset`. The grouped-bar case — multiple series across shared categories, bars adjacent rather than stacked — is structurally identical to `StackedBarChart`'s `[]StackedBarSeries` input but **without** `scales.{x,y}.stacked: true`. The data model and ordering approach already exist in the codebase (`StackedBarSeries`, `charts.go:76-97`); only the stacked flag differs.

## Severity

**Nice-to-have** — fills an obvious gap between single-series `BarChart` and `StackedBarChart`; reuses an existing data shape.

## Proposed change / acceptance criteria

1. Support multi-series grouped bars without breaking single-series `BarChart`. Preferred: a `GroupedBarChart` type reusing the existing ordered series shape (mirror `StackedBarChart`: `Series []StackedBarSeries` or a shared/renamed equivalent) with `NewGroupedBarChart(title string, series []...) *GroupedBarChart` and an `IsHorizontal bool`. Reuse, do not duplicate, the stacked-bar dataset-building logic; the only config delta is omitting `stacked`.
2. Do not change `BarChart`'s single-series API or output.
3. Use the ordered slice shape (consistent with `LineChart`/`StackedBarChart` and the direction of [003-chart-data-ordering.md](003-chart-data-ordering.md)); do **not** introduce a new `map[string]float64` field (that would add ordering debt 003 must then also fix — see Dependencies).
4. Per-series colors from the theme palette; legend shown for >1 series.
5. Tests: N series over M categories produce N datasets, M labels, no `stacked`, deterministic order; horizontal variant sets `indexAxis`.
6. README/example: grouped vs. stacked shown side by side.

## Dependencies

- **Conflict-avoidance with [003-chart-data-ordering.md](003-chart-data-ordering.md):** 003 migrates existing map-based chart data to an ordered representation. This spec must adopt the ordered slice shape from the outset so it neither depends on 003 nor adds new map-ordering debt. If 003 introduces a canonical ordered series/`DataPoint` type, this chart should converge on it; until then reuse `StackedBarSeries`.
- Independently completable against current architecture; cleaner via the open seam if [005](005-extensibility-element-interface.md)/[015](015-html-monolith-refactor.md) landed. No dependency on 001–019.
- Common options via [020-shared-chart-options.md](020-shared-chart-options.md) when available; not blocked by it.
