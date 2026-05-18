# 027 — No mixed (combo) chart support

**Problem:** Chart.js supports mixed charts — multiple datasets of *different* types (e.g. bars + an overlaid line) in one chart by setting `type` per dataset. `rptgen` has no way to combine chart types, a common reporting need (e.g. monthly revenue bars with a trend line).

## Background / context

A mixed chart is a single `{ type, data, options }` where the top-level `type` is typically `"bar"` and individual datasets set their own `type` (e.g. `type: "line"`). `chartDataset` (`html.go:325-333`) has no per-dataset `type` field, so this is impossible today. All the per-type dataset construction logic (bar fill/colors, line tension/fill/points) already exists in `renderBarChartScript`/`renderLineChartScript` and can be reused per dataset. This is the last and most composite item; scheduling it after the single-type charts lets it reuse their builders.

## Severity

**Nice-to-have** — composite feature; deliberately last in the expansion so it reuses the bar/line dataset builders rather than reimplementing them.

## Proposed change / acceptance criteria

1. Add a `MixedChart` whose series each declare a kind, e.g. `MixedSeries struct{ Name string; Kind string /* "bar" | "line" */; Points []DataPoint }` over a shared ordered category-label set; `NewMixedChart(title string, series []MixedSeries) *MixedChart`. Keep `Kind` a small closed set (bar, line) initially.
2. Add an optional per-dataset `Type string` (`omitempty`) to the internal dataset model — additive; existing single-type charts omit it and serialize byte-identically.
3. Reuse existing per-kind dataset construction (colors, line tension/fill, bar background) rather than duplicating; the top-level `type` stays `"bar"` with line datasets overriding via their own `type`.
4. Shared, ordered category labels across all series (same first-appearance/sorted collection as `renderLineChartScript`, `html.go:396-406`); align the point shape with [003-chart-data-ordering.md](003-chart-data-ordering.md)'s ordered representation / the shape used by 022.
5. Legend shown (mixed charts are inherently multi-series); per-series theme colors.
6. Tests: a bar+line `MixedChart` → datasets with correct individual `type`s, shared labels, deterministic order; all other charts' output unchanged.
7. README/example: one bars+trend-line example.

## Dependencies

- Independently completable against current architecture (additive per-dataset `type` field + new case). No dependency on 001–019.
- Reuses the bar/line dataset builders — schedule after [021-area-chart.md](021-area-chart.md)/[022-grouped-bar-chart.md](022-grouped-bar-chart.md) so those builders are factored for reuse (soft ordering; not a hard block — current builders can be reused as-is).
- Data-shape consistency tracked against [003-chart-data-ordering.md](003-chart-data-ordering.md) (same note as [022](022-grouped-bar-chart.md)); cleaner via the open seam if [005](005-extensibility-element-interface.md)/[015](015-html-monolith-refactor.md) landed.
- Common options via [020-shared-chart-options.md](020-shared-chart-options.md) when available.
