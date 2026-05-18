# 026 — No bubble chart support

**Problem:** Chart.js `type: "bubble"` plots `{x, y, r}` points (r = bubble radius in px). `rptgen` has no bubble chart; once [025-scatter-chart.md](025-scatter-chart.md) lands the only gap is the radius dimension.

## Background / context

Bubble is scatter plus a per-point radius. It reuses the entire XY infrastructure introduced by [025-scatter-chart.md](025-scatter-chart.md): linear `scales.x`/`scales.y`, point-array datasets, per-series colors. The single addition is an `R` field on the point and `type: "bubble"`. Keeping bubble dependent on the shared XY shape (rather than a parallel one) is what makes "switching chart types easy" hold.

## Severity

**Nice-to-have** — smallest of the XY family once 025 exists; thin delta over scatter.

## Proposed change / acceptance criteria

1. Add a bubble point reusing the XY shape — extend with radius rather than defining a wholly new type. Options (pick the one consistent with how 025 modelled `XYPoint`): a `BubblePoint struct{ X, Y, R float64 }`, or reuse `XYPoint` plus a parallel radius. Prefer one explicit `BubblePoint` to keep `r` mandatory and discoverable.
2. Add `BubbleChart` with `Series []BubbleSeries` (`Name string; Points []BubblePoint`) and `NewBubbleChart(title string, series []BubbleSeries) *BubbleChart`.
3. Emit `type: "bubble"`; reuse the scatter dataset/scale builder, adding `r` to each serialized point. No change to scatter or any other chart's output.
4. Tests: points serialize as `{"x":..,"y":..,"r":..}`; reuses linear scales; multi-series → one dataset per series.
5. README/example: one bubble example.

## Dependencies

- **Depends on [025-scatter-chart.md](025-scatter-chart.md)** for the XY point/dataset/scale infrastructure — the only intra-expansion hard ordering. Schedule strictly after 025. Not independently completable before 025 (would otherwise duplicate the XY shape, violating the shared-shape design goal).
- No dependency on 001–019; cleaner via the open seam if [005](005-extensibility-element-interface.md)/[015](015-html-monolith-refactor.md) landed.
- Common options via [020-shared-chart-options.md](020-shared-chart-options.md) when available.
