# 019 — `NumberTile` formatting is fragile: unvalidated verb, naive thousands separator

**Problem:** `NumberTile.Format` is fed straight to `fmt.Sprintf` with no validation, so a wrong verb produces visible garbage in the report, and the thousands-separator logic post-processes the formatted string in a way that breaks on non-`%f`-style outputs.

## Background / context

`NumberTile.FormatValue` (`elements.go:25`):

```go
s = fmt.Sprintf(n.Format, n.Value)   // n.Value is float64
```

- If a caller sets `Format: "%d"` (a natural mistake, integers are common in tiles), Go emits `%!d(float64=42)` — that string is then HTML-escaped and shown verbatim in the report (`html.go:189`). No validation, no error ([004-error-handling.md](004-error-handling.md)).
- `addThousandsSep` (`elements.go:42`) operates on the already-formatted string and only understands a literal `.` decimal point and 3-digit grouping. For `Format: "%e"`/`"%g"` (scientific) or `%x` it inserts commas into the mantissa/exponent producing nonsense. It also assumes US grouping with no locale option.
- The `ThousandsSep` + `Format` interaction is implicit and only exercised by `%f`-style formats in tests (`elements_test.go:10`).

This is a correctness/robustness gap in the most commonly used element.

## Severity

**Nice-to-have** — narrower impact than the structural issues, but it silently produces wrong-looking reports from reasonable caller mistakes.

## Proposed change / acceptance criteria

1. Validate or constrain `Format`: either restrict to a known-safe set of float verbs and return/error on others (coordinate with [004-error-handling.md](004-error-handling.md) / [008-api-consistency.md](008-api-consistency.md)), or detect `fmt`'s `%!` error sentinel and surface it instead of rendering it.
2. Apply thousands grouping on the numeric value (e.g. format the integer/fractional parts deterministically) rather than string-scanning arbitrary `fmt` output; explicitly document that `ThousandsSep` is only valid with decimal float formats, and ignore/error otherwise.
3. Consider an explicit, typed formatting API (decimals, prefix/suffix, grouping, percent) instead of a raw `fmt` verb string, so misuse is unrepresentable.
4. Add tests for: `%d` misuse, scientific format + `ThousandsSep`, integer values, negative + grouping (already partially covered at `elements_test.go:22`).

## Dependencies

- Error/validation strategy shared with [004-error-handling.md](004-error-handling.md) and [008-api-consistency.md](008-api-consistency.md).
- Tests overlap [016-test-quality.md](016-test-quality.md).
