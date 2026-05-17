# rptgen — Migration Overview

## Background

This is a Go rewrite of the C# `SharpReports` NuGet library. SharpReports generates static HTML (and JSON) reports containing metric tiles, tables, charts, and free text. The Go version is a drop-in spiritual equivalent, not a strict 1:1 translation — it must follow Go idioms and best practices.

## Decisions

| Topic | Decision |
|---|---|
| Output type | Library only (`package rptgen`) |
| Go API style | Idiomatic structs — plain struct literals + functional options where needed. No fluent builder. |
| Chart.js delivery | Embed in binary via `go:embed` (Chart.js v4.4.0). Reports work fully offline. |
| Tests | Yes — unit tests for renderers, element formatting, and section/canvas layout. |
| External dependencies | Minimize. Zero required runtime dependencies; Chart.js embedded at build time. |

## Module

```
module github.com/lfrystak/rptgen
```

## Package Layout

```
rptgen/
├── rptgen.go           # Top-level types: Report, ReportSection, Theme
├── elements.go         # Element interface + base types (NumberTile, DateTile, FreeText, Table, Canvas)
├── charts.go           # Chart element types (BarChart, LineChart, PieChart, StackedBarChart)
├── html.go             # HtmlRenderer (implements Renderer interface)
├── assets/
│   └── chartjs/
│       └── chart.umd.min.js   # Embedded Chart.js v4.4.0
└── *_test.go           # Tests alongside source files
```

All types and functions live in package `rptgen`. No sub-packages — the library is small enough that a single package is the right Go call.

## Renderer Interface

```go
type Renderer interface {
    Render(report *Report, theme *Theme) (string, error)
}
```

`HtmlRenderer` implements this. `theme` may be `nil` (defaults apply).

## Go Type Mapping from C#

| C# | Go |
|---|---|
| `IReportElement` | `Element` interface |
| `ReportElementBase` | embedded `BaseElement` struct |
| `Report` | `Report` struct |
| `ReportSection` | `Section` struct |
| `Theme` | `Theme` struct |
| `HtmlRenderer` | `HtmlRenderer` struct |
| `JsonRenderer` | removed — JSON output not needed |
| `ReportBuilder` | removed — callers use struct literals |
| Extension methods | removed — callers call methods directly or construct structs inline |
| `CultureInfo` | `language.Tag` from `golang.org/x/text` OR plain `string` locale |

> **Note on culture/locale:** The C# version used `CultureInfo` for number and date formatting. In Go, use the standard `time` package for date formatting (format strings like `"2006-01-02"`) and `strconv`/`fmt` for numbers with a configurable locale string (e.g. `"en-US"`, `"de-DE"`) defaulting to `""` (invariant/Go default). Avoid pulling in heavy i18n dependencies unless formatting fidelity demands it.

## Work Units (Spec Files)

Each spec below is a self-contained agent task. They must be completed in order because later tasks depend on types defined in earlier ones.

| File | Task | Depends on |
|---|---|---|
| `01-core-types.md` | Define `Report`, `Section`, `Element` interface, `Theme` | — |
| `02-elements.md` | Implement `NumberTile`, `DateTile`, `FreeText`, `Table`, `Canvas` | 01 |
| `03-charts.md` | Implement `BarChart`, `LineChart`, `PieChart`, `StackedBarChart` | 01 |
| `04-html-renderer.md` | Implement `HtmlRenderer` with embedded Chart.js, full CSS, all element rendering | 01, 02, 03 |
| `05-tests.md` | Write unit tests for all packages | 01–04 |
