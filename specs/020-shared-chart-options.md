# 020 — No shared model for the common Chart.js options every chart type supports

**Problem:** The internal Chart.js options model exposes only a tiny, ad-hoc subset (responsive, aspect ratio, one legend toggle, `indexAxis`, axis `stacked`), and none of it is reachable from the public API — so callers cannot set the *standard* options Chart.js shares across every chart type (chart title, legend position, axis titles, value-axis min/max, tooltip toggle). This is the foundation the chart-type expansion (021–027) builds on.

## Background / context

Chart.js configuration is uniformly `{ type, data, options }`. A large, well-defined set of `options` is **common to all chart types** — notably `responsive`/`maintainAspectRatio`/`aspectRatio`, and the `plugins.legend`, `plugins.title`, `plugins.subtitle`, `plugins.tooltip` blocks. Cartesian charts (bar, line, area, scatter, bubble, mixed) additionally share the `scales.x` / `scales.y` model (axis `title`, `min`, `max`, `grid.display`, `beginAtZero`); radial charts (radar, polar area) share `scales.r`; circular charts (pie, doughnut) have no scales.

Today (`html.go:335-358`) the Go model is:

```go
type chartOptions struct {
    Responsive  bool
    AspectRatio *float64
    IndexAxis   string
    Plugins     *chartPlugins // only Legend.Display
    Scales      *chartScales  // only X/Y.Stacked
}
```

`ChartBase` (`charts.go:3-8`) exposes only `Title` (rendered as an HTML `<h3>`, *not* the Chart.js title) and a card `Tooltip`. There is no way to set a legend position, an axis label, or a y-axis minimum on any chart — current or future. Every new chart type in 021–027 would otherwise reinvent or hardcode these.

## Severity

**Important** — foundational enabler for the chart-type expansion; without a shared options model each new chart type duplicates option plumbing and the "options common where possible" design goal is unmet. No behaviour regression on its own.

## Proposed change / acceptance criteria

1. Introduce one exported, chart-agnostic options struct (e.g. `ChartOptions`) covering only the **common** Chart.js options, all optional with zero-values meaning "Chart.js default":
   - `LegendPosition string` (`""`/`top`/`bottom`/`left`/`right`/`none`; `none` hides the legend; empty = current default behaviour).
   - `XAxisTitle`, `YAxisTitle string` (cartesian charts only; ignored by pie/doughnut/polar/radar).
   - `YMin`, `YMax *float64` (value-axis bounds; applied to the value axis regardless of orientation).
   - `ShowTooltips *bool` (toggles `plugins.tooltip.enabled`; nil = default on).
   - `AspectRatio *float64` (already plumbed internally; expose it).
   - Optionally `ShowChartTitle bool` to additionally emit the Chart.js native `plugins.title` from the existing `ChartBase.Title` (default keeps today's HTML-`<h3>` behaviour to avoid a visual regression).
2. Embed it in `ChartBase` (e.g. `ChartBase.Options ChartOptions`) so **every** existing chart (`BarChart`, `LineChart`, `PieChart`, `StackedBarChart`) and every future one (021–027) inherits it with no per-type code.
3. Expand the internal `chartOptions`/`chartPlugins`/`chartScales` model additively so these map to the correct Chart.js JSON; fields stay `omitempty`/pointer so unset options serialize exactly as today (no golden-file churn for charts that don't set them).
4. Axis-title/min/max are silently ignored for chart types without that scale (documented), not errors.
5. Tests: a chart with `LegendPosition`, axis titles, and `YMin`/`YMax` produces the expected `options` JSON; a chart with a zero-value `ChartOptions` produces byte-identical output to before this change.
6. Document the common-options model once, in package docs and README, and have 021–027 reference it rather than re-document per type.

## Dependencies

- **No hard dependency on 001–019.** Compatible with — and cleaner after — [015-html-monolith-refactor.md](015-html-monolith-refactor.md) (config model "open for extension"): if 015 has landed, add these as the shared/common slice of the open model; if not, extend the existing closed structs additively. Either order works; note the sequencing, do not block on it.
- Foundation reused by [021-area-chart.md](021-area-chart.md)–[027-mixed-chart.md](027-mixed-chart.md). Those specs are independently completable without this one (they fall back to current defaults); they only *gain* common options once this lands.
- Doc addition coordinates with [013-documentation-drift.md](013-documentation-drift.md) (keep README claims accurate).
