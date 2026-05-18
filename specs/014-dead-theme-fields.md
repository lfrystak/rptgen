# 014 — Dead / misleading `Theme` fields: `AccentColor` and `EnableGradients`

**Problem:** `Theme` exposes configuration that has little or no effect, so callers set values expecting a result that never happens.

## Background / context

- `Theme.AccentColor` (`rptgen.go:77`) is set to `#10b981` by `DefaultTheme()` (`rptgen.go:94`) but is **never referenced** anywhere in `html.go` (`grep AccentColor *.go` shows only the struct/default). It styles nothing.
- `Theme.EnableGradients` (`rptgen.go:83`) is used at exactly one place — `html.go:549` — where it only adds a gradient to the **report header background**. It does not affect charts, despite README describing it as "Gradient fills on chart bars" ([013-documentation-drift.md](013-documentation-drift.md)).

Exposed-but-inert configuration is a maintainability and trust problem: it widens the public API surface (every exported field is a compatibility commitment) while doing nothing, and it makes the theme harder to reason about.

## Severity

**Important** — public API includes knobs that silently do nothing; either wire them up or remove them before more users depend on them.

## Proposed change / acceptance criteria

Decide per field:

1. **`AccentColor`:** either (a) actually use it (e.g. for `Subtitle`, tooltip icon, chart accent, links) and add a test asserting it appears in output, or (b) remove the field and its default and document the removal.
2. **`EnableGradients`:** either (a) extend it to charts/cards so behaviour matches a corrected description, or (b) rename/scope its documentation to "header gradient only." Add a test asserting the on/off output difference.
3. Whatever is decided, `README.md` theme table must match exactly (coordinate with [013-documentation-drift.md](013-documentation-drift.md)).
4. Add tests that fail if a documented theme field stops affecting output (guards against future inert fields).

## Dependencies

- README reconciliation in [013-documentation-drift.md](013-documentation-drift.md).
- Test additions overlap [016-test-quality.md](016-test-quality.md).
