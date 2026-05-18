# 024 — No polar area chart support

**Problem:** Chart.js supports `type: "polarArea"` (pie-like segments on a radial value axis), but `rptgen` has no polar area chart.

## Background / context

A polar area chart has the *same input shape as `PieChart`*: a set of labelled values, one dataset, per-segment colors. It differs only in `type: "polarArea"` and the presence of a radial `scales.r`. `renderPieChartScript` (`html.go:445-470`) is almost exactly the needed builder — swap the type string and add the optional radial scale. This is the lowest-effort of the new chart types after area.

## Severity

**Nice-to-have** — small, standalone; closely mirrors the existing pie path.

## Proposed change / acceptance criteria

1. Add a `PolarAreaChart` type with the same ordered labelled-value data shape used by `PieChart` (align with whatever ordered shape [003-chart-data-ordering.md](003-chart-data-ordering.md) settles on; until then mirror `PieChart` exactly so the two stay consistent) and `NewPolarAreaChart(title string, data ...) *PolarAreaChart`.
2. Emit `type: "polarArea"`; one dataset, per-segment theme colors (reuse the pie color-cycling logic), legend shown by default (as pie).
3. Reuse the optional radial scale added for [023-radar-chart.md](023-radar-chart.md) if that has landed; otherwise add the same additive `scales.r` (kept `omitempty`). Do not duplicate the radial-scale struct — share it with radar.
4. Tests: labelled values → one dataset, correct labels/colors, `type:"polarArea"`; output of other charts unchanged.
5. README/example: one polar area example.

## Dependencies

- Independently completable against current architecture. No dependency on 001–019.
- Shares the radial-scale model with [023-radar-chart.md](023-radar-chart.md); if both are scheduled, sequence 023 first so the `scales.r` struct exists, or factor it in whichever lands first. Not a hard ordering requirement (each can add it independently and dedupe on merge).
- Data-shape consistency tracked against [003-chart-data-ordering.md](003-chart-data-ordering.md) (same note as [022](022-grouped-bar-chart.md)).
- Common options via [020-shared-chart-options.md](020-shared-chart-options.md) when available.
