# rptgen — Engineering Review Tracking Index

Single source of truth for work items from the standalone engineering review of `rptgen`.
Ordered by priority (severity, then enabling/blocking relationships). Update the **Status**
column as work progresses: `🔜 Not started` · `🔄 In progress` · `✅ Done` · `❌ Won't do`.

| ID | Title | Severity | Status |
|----|-------|----------|--------|
| [001](001-charts-not-self-contained.md) | Reports are not self-contained: Chart.js loaded from CDN | Critical | ✅ Done |
| [002](002-script-context-injection.md) | User data can break out of the `<script>` block (XSS / corruption) | Critical | ✅ Done |
| [003](003-chart-data-ordering.md) | Chart category ordering destroyed by `map` + alphabetical sort | Critical | ✅ Done |
| [004](004-error-handling.md) | Library panics on bad input, swallows errors, returns always-nil error | Important | ✅ Done |
| [005](005-extensibility-element-interface.md) | `Element` interface is illusory: rendering is a hardcoded type switch | Important | 🔜 Not started |
| [006](006-render-to-io-writer.md) | `Render` returns buffered `string` instead of writing to `io.Writer` | Important | 🔜 Not started |
| [007](007-baseelement-dead-boilerplate.md) | `BaseElement` / `newBaseElement()` is dead boilerplate | Important | 🔜 Not started |
| [008](008-api-consistency.md) | Inconsistent public API: constructors vs raw struct literals | Important | 🔜 Not started |
| [009](009-ci-workflow-defects.md) | CI: placeholder step, non-enforcing lint, missing quality gates | Important | ✅ Done |
| [010](010-goreleaser-library-misconfig.md) | GoReleaser builds binaries for a library | Important | ✅ Done |
| [011](011-gomod-pins-patch-version.md) | `go.mod` pins a patch toolchain version (`go 1.26.3`) | Important | ✅ Done |
| [012](012-gofmt-not-clean.md) | Source is not `gofmt`-clean | Important | ✅ Done |
| [013](013-documentation-drift.md) | README documents behaviour the code does not implement | Important | 🔜 Not started |
| [014](014-dead-theme-fields.md) | Dead/misleading `Theme` fields: `AccentColor`, `EnableGradients` | Important | 🔜 Not started |
| [015](015-html-monolith-refactor.md) | `html.go` is a 765-line monolith mixing five concerns | Important | ✅ Done |
| [016](016-test-quality.md) | Tests are smoke-level; key paths uncovered | Important | ✅ Done |
| [017](017-package-godoc.md) | No package-level documentation (godoc / pkg.go.dev) | Nice-to-have | 🔜 Not started |
| [018](018-html-accessibility-quality.md) | Generated HTML has accessibility / quality gaps | Nice-to-have | 🔜 Not started |
| [019](019-numbertile-format-robustness.md) | `NumberTile` formatting is fragile (unvalidated verb, naive grouping) | Nice-to-have | 🔜 Not started |

## Chart.js expansion (020–027)

Feature expansion to support all Chart.js chart types holistically. Each item is
independently completable in a small step against the current architecture; no hard
dependency on 001–019 (see each spec's **Dependencies** for coordination notes).

| ID | Title | Severity | Status |
|----|-------|----------|--------|
| [020](020-shared-chart-options.md) | No shared model for common Chart.js options across chart types | Important | 🔜 Not started |
| [021](021-area-chart.md) | No area chart support | Nice-to-have | 🔜 Not started |
| [022](022-grouped-bar-chart.md) | `BarChart` is single-series only; no grouped (multi-series) bar chart | Nice-to-have | 🔜 Not started |
| [023](023-radar-chart.md) | No radar chart support | Nice-to-have | 🔜 Not started |
| [024](024-polar-area-chart.md) | No polar area chart support | Nice-to-have | 🔜 Not started |
| [025](025-scatter-chart.md) | No scatter chart support; no `{x,y}` point data shape | Nice-to-have | 🔜 Not started |
| [026](026-bubble-chart.md) | No bubble chart support | Nice-to-have | 🔜 Not started |
| [027](027-mixed-chart.md) | No mixed (combo) chart support | Nice-to-have | 🔜 Not started |

## Suggested sequencing

The fixes interlock; a reasonable order that minimizes rework:

1. **Hygiene first:** [012](012-gofmt-not-clean.md) → [009](009-ci-workflow-defects.md), [011](011-gomod-pins-patch-version.md), [010](010-goreleaser-library-misconfig.md) (independent, unblock a trustworthy CI).
2. **Safety net:** [016](016-test-quality.md) golden tests before structural change.
3. **Structure:** [015](015-html-monolith-refactor.md) (templating) — enables [002](002-script-context-injection.md) and [005](005-extensibility-element-interface.md).
4. **Behaviour/API:** [001](001-charts-not-self-contained.md), [002](002-script-context-injection.md), [003](003-chart-data-ordering.md), [004](004-error-handling.md), [005](005-extensibility-element-interface.md), [006](006-render-to-io-writer.md), [007](007-baseelement-dead-boilerplate.md), [008](008-api-consistency.md), [014](014-dead-theme-fields.md), [019](019-numbertile-format-robustness.md).
5. **Polish/docs last (reflect final behaviour):** [013](013-documentation-drift.md), [017](017-package-godoc.md), [018](018-html-accessibility-quality.md).

### Chart.js expansion ordering

Independent of 001–019 and schedulable in parallel, but cleaner *after* the structural
specs land:

1. **Cleaner after** [015](015-html-monolith-refactor.md) (config model open for extension) and
   [005](005-extensibility-element-interface.md) (open element/render seam) — new chart types
   then add no central-switch edits. Not a hard block; each chart spec is also completable
   against the current closed architecture.
2. **Foundation:** [020](020-shared-chart-options.md) (others gain common options once it lands; not a blocker).
3. **Standalone, any order:** [021](021-area-chart.md), [022](022-grouped-bar-chart.md), [023](023-radar-chart.md), [024](024-polar-area-chart.md).
4. **XY family (ordered):** [025](025-scatter-chart.md) → [026](026-bubble-chart.md) (hard dep on 025).
5. **Last (reuses bar/line builders):** [027](027-mixed-chart.md).

**Coordination flags:** [025](025-scatter-chart.md)'s `ScatterChart` is the named acceptance
vehicle for [005](005-extensibility-element-interface.md) and [015](015-html-monolith-refactor.md)
— pick one owner/ordering so it is introduced once. New chart types must use ordered
slice data shapes (not `map[string]float64`) to stay consistent with, and avoid adding
debt for, [003](003-chart-data-ordering.md).
