# 003 — Chart category ordering is destroyed by `map` + alphabetical sort

**Problem:** Single-series charts store data as `map[string]float64` and render categories sorted alphabetically, so logically-ordered data (months, quarters, buckets) renders in the wrong order.

## Background / context

- `BarChart.Data`, `PieChart.Data`, `LineSeries.Points` are all `map[string]float64` (`charts.go:13,29,63`).
- Maps have no order, so `html.go` uses `sortedKeys` (`html.go:293`) which `sort.Strings` the keys for "deterministic chart rendering."

Deterministic, but **wrong**: the library's own example (`example/main.go:24`) builds `monthlyRevenue` for a `NewLineChartSingle("Monthly Revenue Trend", …)`. With alphabetical sorting the x-axis renders `April, February, January, June, March, May` instead of `January … June`. `salesByRegion` and any numeric-bucket labels (`"0-10"`, `"10-20"`, `"100+"`) are similarly mangled. For a charting library this is a core correctness defect, not a cosmetic one.

`StackedBarChart` and `LineChart` already got this right by using **ordered slices** (`[]StackedBarSeries`, `[]LineSeries`) with comments explicitly noting "slice to guarantee deterministic ... order" — so the design intent is known; it just wasn't applied to the single-series/data-point types.

## Severity

**Critical** — produces incorrect charts from the documented, primary use case.

## Proposed change / acceptance criteria

1. Replace order-losing maps with an ordered representation for `BarChart`, `PieChart`, and `LineSeries` data points. Options (pick one, apply consistently):
   - An exported ordered pair type, e.g. `type DataPoint struct { Label string; Value float64 }` and `Data []DataPoint`.
   - Parallel `Labels []string` + `Values []float64` with constructor validation.
2. Preserve caller insertion order end-to-end through `renderElement` → chart config; remove `sortedKeys` from these paths (it may remain only where alphabetical order is genuinely desired and documented).
3. Keep ergonomic constructors: `NewBarChart`, `NewPieChart`, `NewLineChartSingle` should accept the ordered type; consider a map-accepting helper that documents it sorts keys, for callers who explicitly want that.
4. Update `README.md`, `example/main.go`, and tests; add a test asserting that input order `[Mar, Jan, Feb]` renders as `[Mar, Jan, Feb]`, not `[Feb, Jan, Mar]`.
5. Re-render `example/report.html` to confirm month/region order is correct.

## Dependencies

- API-shape change coordinates with [008-api-consistency.md](008-api-consistency.md).
- Doc/example updates overlap [013-documentation-drift.md](013-documentation-drift.md).
- New ordering tests overlap [016-test-quality.md](016-test-quality.md).
