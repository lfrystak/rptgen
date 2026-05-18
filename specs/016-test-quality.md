# 016 — Tests are smoke-level: substring assertions, no golden files, key paths uncovered

**Problem:** The test suite asserts coarse substring presence on the rendered document and does not validate chart configuration correctness, ID-collision handling, security escaping, or output structure — so regressions in the renderer can pass CI.

## Background / context

- `html_test.go` checks `strings.Contains(out, "new Chart(")`, `"cdn.jsdelivr.net"`, `"1fr 2fr"`, etc. — it never parses the HTML, never validates the emitted Chart.js JSON, and never checks element structure/order.
- Coverage gaps (from `go test -coverprofile`): `idGen.next` 63.6% (the `-2` collision-suffix branch — i.e. two charts with the same title in a section — is **untested**), `shadowCSS` 40% (only "medium"/empty hit; `"none"/"subtle"/"strong"` untested), `tooltipIcon`/`chartColors` 66.7%.
- No test for: horizontal bar (`IsHorizontal`), donut (`IsDonut`), multi-series line label union (`html.go:397`), `StackedBarChart` series ordering, `EnableGradients`/`EnableAnimations` on/off output difference, `NewTableFromColumns` value/order correctness (`elements_test.go:102` only asserts row count, not contents/order), table cell rendering of non-string `any` values.
- No security regression test for the injection issue in [002-script-context-injection.md](002-script-context-injection.md).
- No golden-file test, so structural HTML/CSS regressions during the [015](015-html-monolith-refactor.md) refactor would be invisible.
- `example/` has 0% coverage and is excluded from Sonar coverage — acceptable, but there is no compile/run guard that the example still works (it is the de-facto integration scenario).

## Severity

**Important** — test weakness undermines every other change here; refactors ([015](015-html-monolith-refactor.md)) and behaviour fixes need a real safety net.

## Proposed change / acceptance criteria

1. Add golden-file tests: render `buildFullReport()` (and a custom-theme variant) and compare against committed `testdata/*.html`, with an `-update` flag. Covers structure during refactors.
2. Add focused tests: parse/inspect the emitted chart JSON config and assert type, labels (incl. order per [003](003-chart-data-ordering.md)), datasets, stacked axes, indexAxis for horizontal, doughnut for donut.
3. Cover `idGen.next` collision suffix (two same-titled charts → `-2`), all `shadowCSS` branches, `tooltipIcon` empty/non-empty, `chartColors` theme-override path.
4. Add the security regression test from [002-script-context-injection.md](002-script-context-injection.md).
5. Strengthen `NewTableFromColumns` test to assert per-row values and column order, and add a non-string `any` cell test.
6. Add a build/run guard for `example/` (e.g. `go build ./example/...` in CI, or a test invoking the report-building code with output to a temp dir).
7. Target meaningful coverage (e.g. ≥90% lib statements with assertions that actually verify output, not just presence).

## Dependencies

- Golden tests should land with/before [015-html-monolith-refactor.md](015-html-monolith-refactor.md).
- Specific assertions depend on final shapes from [001](001-charts-not-self-contained.md), [002](002-script-context-injection.md), [003](003-chart-data-ordering.md), [014](014-dead-theme-fields.md).
- CI wiring (race, example build) coordinates with [009-ci-workflow-defects.md](009-ci-workflow-defects.md).
