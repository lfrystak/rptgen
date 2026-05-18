# 013 — README documents behaviour the code does not implement

**Problem:** Several README statements are factually wrong against the current code, misleading users about offline support, defaults, and features.

## Background / context

Concrete mismatches between `README.md` and the implementation:

1. **"Self-contained / no external dependencies / Chart.js embedded inline" (`README.md:12,196,323`)** — false: `html.go:79` loads Chart.js from `cdn.jsdelivr.net`. (Root cause tracked in [001-charts-not-self-contained.md](001-charts-not-self-contained.md); the doc must be reconciled once behaviour is decided.)
2. **`BackgroundColor` default (`README.md:311`)** — README says `#ffffff`; `DefaultTheme()` (`rptgen.go:92`) sets `#f1f5f9`.
3. **`EnableGradients` "Gradient fills on chart bars" (`README.md:319`)** — false: the flag only alters the **header** background CSS (`html.go:549`); charts are unaffected. (See [014-dead-theme-fields.md](014-dead-theme-fields.md).)
4. **`AccentColor` "Highlight elements" (`README.md:313`)** — the field is set in `DefaultTheme()` but never referenced in `generateCSS`/charts; it highlights nothing. (See [014-dead-theme-fields.md](014-dead-theme-fields.md).)
5. **"Custom renderers can be implemented by satisfying the `Renderer` interface" (`README.md:329`)** — technically true but practically unsupported; no reusable dispatch is exported (see [005-extensibility-element-interface.md](005-extensibility-element-interface.md)).
6. Minimal example (`README.md:43-47`) ignores `Render`'s returned error after the `log.Fatal` branch and drops the `os.WriteFile` error — sets a poor pattern (relates to [004-error-handling.md](004-error-handling.md)).

## Severity

**Important** — documentation is the public contract for a library; these errors actively mislead.

## Proposed change / acceptance criteria

1. Reconcile every point above so README matches actual behaviour **after** the corresponding behaviour specs are resolved (don't document aspirational behaviour).
2. Fix the `BackgroundColor` default value in the theme table to match `DefaultTheme()` (or change the default — pick one, keep consistent).
3. Correct or remove the `EnableGradients`/`AccentColor` descriptions per [014-dead-theme-fields.md](014-dead-theme-fields.md).
4. Correct the extensibility section per [005-extensibility-element-interface.md](005-extensibility-element-interface.md).
5. Update README/example to the final API shape from [003](003-chart-data-ordering.md), [006](006-render-to-io-writer.md), [008](008-api-consistency.md) and to handle errors properly.
6. Add a doc test or check (where feasible) so default-value drift is caught automatically — e.g. a test asserting documented defaults equal `DefaultTheme()` values.

## Dependencies

- Downstream of [001](001-charts-not-self-contained.md), [003](003-chart-data-ordering.md), [005](005-extensibility-element-interface.md), [006](006-render-to-io-writer.md), [008](008-api-consistency.md), [014](014-dead-theme-fields.md). Update docs last, once behaviour is final.
