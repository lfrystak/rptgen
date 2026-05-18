# 018 — Generated HTML has accessibility / quality gaps

**Problem:** The rendered document hardcodes a non-configurable language, uses a static non-descriptive logo `alt`, and emits tables without header semantics — reducing accessibility and output quality.

## Background / context

- `html.go:31` always emits `<html lang="en">` regardless of report content/locale; no way to override.
- `html.go:44` renders the logo as `<img src="…" alt="logo">` — a static, non-descriptive alt for every report; should derive from the report title or be caller-configurable.
- `renderTable` (`html.go:230`) emits `<th>` cells but without `scope="col"`, and the table has no `<caption>` (the title is rendered as a separate `<h3>` outside the table, breaking the table/caption association for assistive tech).
- Tooltips are conveyed only via a CSS `::after` on a `<span data-tooltip>` (`html.go:282`, CSS `html.go:704`) with no `aria`/`title`, so the information is invisible to screen readers and keyboard users (`cursor: help`, hover-only).

These do not break rendering but lower the quality bar for a tool whose entire output is a shareable HTML document.

## Severity

**Nice-to-have** — output works; accessibility and polish improvements.

## Proposed change / acceptance criteria

1. Make document language configurable (e.g. `Report.Lang`, default `"en"`).
2. Use a meaningful logo `alt` (e.g. `"<Report.Title> logo"`) or expose a `LogoAlt` field.
3. Add `scope="col"` to `<th>` and render the table title as `<caption>` inside `<table>` (keep visual styling).
4. Expose tooltip text to assistive tech (e.g. add `title`/`aria-label`, make focusable) so it is not hover-only.
5. Tests asserting the above attributes are present (coordinate with [016-test-quality.md](016-test-quality.md)).

## Dependencies

- Touches `renderTable`/header rendering — sequence with [015-html-monolith-refactor.md](015-html-monolith-refactor.md) to avoid rework.
- Test additions overlap [016-test-quality.md](016-test-quality.md); doc updates overlap [013-documentation-drift.md](013-documentation-drift.md).
