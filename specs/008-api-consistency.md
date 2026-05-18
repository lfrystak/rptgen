# 008 — Inconsistent public API surface: constructors vs raw struct literals

**Problem:** Construction is inconsistent across element types — charts and tables have `NewXxx` constructors, tiles and `FreeText` do not — so callers must mix two styles, and there is no single documented convention.

## Background / context

- Constructors exist: `NewBarChart`, `NewLineChart`, `NewLineChartSingle`, `NewPieChart`, `NewStackedBarChart` (`charts.go`), `NewTable`, `NewTableWithColumns`, `NewTableFromColumns` (`elements.go`).
- No constructors: `NumberTile`, `DateTile`, `FreeText`, `Section`. The README/example build these with raw struct literals (`example/main.go:58`, `README.md:39`) while building charts with constructors in the same file — visibly inconsistent.
- `Section` and `Canvas` both model "a column grid of elements" with `AddElement` chaining, but `Section` has no constructor and `Canvas` has `NewCanvas`; `Section.ColumnWidths` vs `Canvas.ColumnWidths` duplicate the same concept with no shared type.
- Constructors also disagree on error strategy: `NewTableFromColumns` panics (see [004-error-handling.md](004-error-handling.md)) while others cannot fail. No documented rule for when a constructor exists or whether it validates.

The mixed style leaks into every example and makes the API harder to learn and document; it also blocks invariants (e.g. requiring `BaseElement`/seal, validating `Format`) from being enforced at construction.

## Severity

**Important** — coherence of the primary public API; compounds with [004](004-error-handling.md), [007](007-baseelement-dead-boilerplate.md), [019](019-numbertile-format-robustness.md).

## Proposed change / acceptance criteria

1. Decide and document one construction convention, e.g.: every element has a `NewXxx` constructor for required fields; optional fields set via exported fields or functional options. Apply uniformly (add `NewNumberTile`, `NewDateTile`, `NewFreeText`, `NewSection`).
2. Unify the section/canvas grid concept (shared `ColumnWidths` type or shared embedded grid), or clearly document why they differ.
3. Standardize constructor error behaviour with [004-error-handling.md](004-error-handling.md) (no panics for caller input).
4. Keep raw struct literals working where reasonable (don't force constructors via unexported required fields unless [005](005-extensibility-element-interface.md)/[007](007-baseelement-dead-boilerplate.md) require sealing); the goal is a consistent *recommended* path.
5. README and `example/main.go` updated to use the single recommended style throughout.

## Dependencies

- [004-error-handling.md](004-error-handling.md) (constructor error convention), [007-baseelement-dead-boilerplate.md](007-baseelement-dead-boilerplate.md) (what constructors must set), [003-chart-data-ordering.md](003-chart-data-ordering.md) (constructor input types), [013-documentation-drift.md](013-documentation-drift.md) (docs/example rewrite).
