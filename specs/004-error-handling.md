# 004 — Error handling: library panics on bad input, swallows errors, returns always-nil error

**Problem:** The package mixes three anti-patterns — panicking on recoverable user input, silently swallowing a marshal failure, and a `Render` signature that promises an `error` it never returns.

## Background / context

1. **Panic as input validation.** `NewTableFromColumns` (`elements.go:138`) `panic("rptgen: NewTableFromColumns: column slices have different lengths")` on mismatched lengths. This is normal caller-data error, not a programmer invariant; a library should return an `error`. `TestNewTableFromColumnsPanicsOnMismatch` (`elements_test.go:114`) currently codifies the panic as intended behaviour.
2. **Swallowed error.** `chartInitScript` (`html.go:305`) on `json.Marshal` failure does `configJSON = []byte("{}")` and proceeds, emitting a chart with no data and no signal to the caller. The failure is invisible.
3. **Phantom error return.** `HtmlRenderer.Render` (`html.go:21`) is typed `(string, error)` but every `return` is `nil` error. The `Renderer` interface (`rptgen.go:104`) propagates this. Callers (and `example/main.go:198`) write `if err != nil` branches that are dead code, and a future renderer that *can* fail has no precedent for how errors surface.
4. **Silent unknown-element fallthrough.** `renderElement` default branch (`html.go:179`) emits `<!-- unknown element -->` and returns no error — partial/incorrect output with no diagnostic. (See also [005-extensibility-element-interface.md](005-extensibility-element-interface.md).)

## Severity

**Important** — affects API safety, debuggability, and correctness; #1 can crash a caller's process.

## Proposed change / acceptance criteria

1. Replace the `NewTableFromColumns` panic with an error return: either `NewTableFromColumns(...) (*Table, error)` or a validating variant. Decide and document the convention for all constructors (see [008-api-consistency.md](008-api-consistency.md)). Update the test to assert the returned error.
2. `chartInitScript` must propagate a `json.Marshal` failure up through `renderElement` → `Render` rather than substituting `{}`. (In practice the config is always marshalable, so this becomes a guard that surfaces programmer errors instead of hiding them.)
3. Either make `Render` genuinely able to return errors (it will, once #1/#2 propagate) and document when, or — if it truly cannot fail — drop `error` from the interface. Do not keep a permanently-nil error.
4. Unknown element types should produce a returned error (or be impossible by construction — see [005](005-extensibility-element-interface.md)), not a silent HTML comment.
5. `go vet` and tests stay green; add tests covering the new error paths.

## Dependencies

- Interacts with [005-extensibility-element-interface.md](005-extensibility-element-interface.md) (unknown-element handling) and [008-api-consistency.md](008-api-consistency.md) (constructor conventions).
