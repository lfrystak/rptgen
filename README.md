# RPTGEN: A Go Report Generator

[![CI](https://github.com/lfrystak/rptgen/actions/workflows/ci.yml/badge.svg)](https://github.com/lfrystak/rptgen/actions/workflows/ci.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=rptgen&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=rptgen)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lfrystak/rptgen)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/github/license/lfrystak/rptgen)](LICENSE.txt)
[![Latest Release](https://img.shields.io/github/v/release/lfrystak/rptgen)](https://github.com/lfrystak/rptgen/releases/latest)

![image](img/rptgen_logo.svg)


`rptgen` is a Go library for building structured, self-contained HTML reports. You compose a report from typed elements — tiles, tables, charts, and free text — arrange them into multi-column sections, and render everything to a single HTML file with no external dependencies.

Reports are themed via a `Theme` struct that controls colors, fonts, shadows, and animations. A built-in default theme is provided; custom themes only require overriding the fields you care about.

## Installation

```sh
go get github.com/lfrystak/rptgen
```

## Minimal Example

```go
package main

import (
    "log"
    "os"

    "github.com/lfrystak/rptgen"
)

func main() {
    report := rptgen.NewReport("Monthly Summary")
    report.Footer = "Internal use only"

    section := &rptgen.Section{Title: "KPIs", ColumnWidths: []int{1, 1}}
    section.AddElement(&rptgen.NumberTile{Title: "Revenue", Value: 98000, Format: "%.0f", Prefix: "$", ThousandsSep: true})
    section.AddElement(&rptgen.NumberTile{Title: "Growth", Value: 12.0, Format: "%.1f%%"})
    report.AddSection(section)

    html, err := rptgen.HtmlRenderer{}.Render(report, nil) // nil = default theme
    if err != nil {
        log.Fatal(err)
    }
    os.WriteFile("report.html", []byte(html), 0644)
}
```

## Report Structure

### `Report`

The top-level document container. Created with `NewReport(title)`.

| Field         | Type       | Description                                      |
|---------------|------------|--------------------------------------------------|
| `Title`       | `string`   | Report heading displayed in the HTML header.     |
| `Sections`    | `[]*Section` | Ordered list of content sections.              |
| `GeneratedAt` | `time.Time` | Timestamp shown in the header. Set automatically by `NewReport`. |
| `Footer`      | `string`   | Text shown at the bottom of the page.            |
| `LogoURL`     | `string`   | URL of a logo image rendered above the title.    |

```go
report := rptgen.NewReport("Q2 Report")
report.LogoURL = "https://example.com/logo.png"
report.Footer  = "Confidential"
```

### `Section`

A section groups elements into a CSS grid row.

| Field          | Type       | Description                                                                 |
|----------------|------------|-----------------------------------------------------------------------------|
| `Title`        | `string`   | Section heading. Empty = no heading.                                        |
| `Elements`     | `[]Element` | Ordered list of elements in the section.                                   |
| `ColumnWidths` | `[]int`    | Proportional column widths. `[]int{1, 2}` = 33 %/67 %. Nil = one column. |

Use `EqualColumns(n)` to create n equal-width columns without spelling out the slice:

```go
// Two equal columns — both are equivalent
section := &rptgen.Section{Title: "Revenue", ColumnWidths: rptgen.EqualColumns(2)}
section := &rptgen.Section{Title: "Revenue", ColumnWidths: []int{1, 1}}

section.AddElement(chart1)
section.AddElement(chart2)
report.AddSection(section)
```

### Elements

#### `NumberTile`

Displays a single numeric metric.

| Field          | Type      | Description                                                              |
|----------------|-----------|--------------------------------------------------------------------------|
| `Title`        | `string`  | Label above the value.                                                   |
| `Value`        | `float64` | The numeric value to display.                                            |
| `Format`       | `string`  | `fmt.Sprintf` format string, e.g. `"%.2f"`, `"%.1f%%"`. Empty = raw number. |
| `Prefix`       | `string`  | String prepended to the formatted value, e.g. `"$"`, `"€"`.            |
| `ThousandsSep` | `bool`    | Insert comma thousands separators into the integer part.                 |
| `Subtitle`     | `string`  | Small caption below the value (e.g. `"↑ vs last month"`).               |
| `Tooltip`      | `string`  | Hover text on the tile card.                                             |

```go
&rptgen.NumberTile{Title: "Revenue", Value: 98000, Format: "%.0f", Prefix: "$", ThousandsSep: true}
&rptgen.NumberTile{Title: "Growth",  Value: 12.0,  Format: "%.1f%%", Subtitle: "↑ vs Q1"}
&rptgen.NumberTile{Title: "Score",   Value: 4.7,   Format: "%.1f",   Tooltip: "Out of 5"}
```

#### `DateTile`

Displays a date or datetime metric.

| Field      | Type        | Description                                                          |
|------------|-------------|----------------------------------------------------------------------|
| `Title`    | `string`    | Label above the value.                                               |
| `Value`    | `time.Time` | The time value. Zero value renders as empty.                         |
| `Format`   | `string`    | Go time layout string. Empty = `"2006-01-02 15:04:05"`.             |
| `Subtitle` | `string`    | Small caption below the date.                                        |
| `Tooltip`  | `string`    | Hover text on the tile card.                                         |

```go
&rptgen.DateTile{
    Title:  "Quarter End",
    Value:  time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
    Format: "January 02, 2006",
}
```

#### `FreeText`

Displays a block of plain text or raw HTML.

| Field     | Type     | Description                                                         |
|-----------|----------|---------------------------------------------------------------------|
| `Content` | `string` | The text (or HTML) to render.                                       |
| `IsHTML`  | `bool`   | If `true`, `Content` is injected as-is. If `false`, it is escaped. |

> **Security note:** When `IsHTML` is `true`, `Content` is written into the document verbatim
> with **no escaping**. If the value originates from user input or any untrusted source you
> must sanitize it (e.g. with [bluemonday](https://github.com/microcosm-cc/bluemonday)) before
> passing it to `FreeText`. Failing to do so allows stored XSS.

```go
&rptgen.FreeText{Content: "Plain paragraph text."}
&rptgen.FreeText{Content: "<p>Rich <strong>HTML</strong> content.</p>", IsHTML: true}
```

#### `Table`

Displays tabular data with a header row.

| Field     | Type              | Description                                               |
|-----------|-------------------|-----------------------------------------------------------|
| `Title`   | `string`          | Caption above the table.                                  |
| `Columns` | `[]string`        | Ordered column names used as headers and row-map keys.    |
| `Rows`    | `[]map[string]any` | Row data. Each map key must match a column name.         |

Three constructors are available:

```go
// Infer columns from first row keys (sorted alphabetically)
rptgen.NewTable("Sales", rows)

// Explicit column order
rptgen.NewTableWithColumns("Sales", rows, []string{"Customer", "Revenue", "Status"})

// Column-oriented input
rptgen.NewTableFromColumns("Sales", map[string][]any{
    "Customer": {"Acme", "TechCo"},
    "Revenue":  {45000, 38000},
})
```

#### `Canvas`

A flexible sub-grid container that nests elements inside a section column.

| Field          | Type       | Description                                                              |
|----------------|------------|--------------------------------------------------------------------------|
| `ColumnWidths` | `[]int`    | Proportional widths for columns within the canvas.                       |
| `Elements`     | `[]Element` | Elements placed inside the canvas grid.                                 |

```go
canvas := rptgen.NewCanvas(rptgen.EqualColumns(2)...) // two equal columns
canvas.AddElement(&rptgen.NumberTile{Title: "Users", Value: 500, Format: "%.0f", ThousandsSep: true})
canvas.AddElement(rptgen.NewBarChart("Trend", data))

section := &rptgen.Section{ColumnWidths: []int{1, 2}}
section.AddElement(canvas)       // occupies the first (narrower) column
section.AddElement(anotherChart) // occupies the second (wider) column
```

### Charts

Charts are rendered using Chart.js, embedded inline — no CDN calls required.

#### `BarChart`

Single-series vertical or horizontal bar chart.

| Field          | Type                  | Description                              |
|----------------|-----------------------|------------------------------------------|
| `Title`        | `string`              | Chart heading.                           |
| `Data`         | `map[string]float64`  | Label → value pairs.                     |
| `IsHorizontal` | `bool`                | Render bars horizontally if `true`.      |
| `Tooltip`      | `string`              | Hover text on the chart card.            |

```go
chart := rptgen.NewBarChart("Sales by Region", map[string]float64{
    "North America": 125000,
    "Europe":        98000,
})
chart.IsHorizontal = true
```

#### `LineChart`

Line chart with one or more named series.

| Field        | Type           | Description                                         |
|--------------|----------------|-----------------------------------------------------|
| `Title`      | `string`       | Chart heading.                                      |
| `Series`     | `[]LineSeries` | Ordered slice of series (preserves legend order).   |
| `ShowPoints` | `bool`         | Show data point dots. Default: `true`.              |
| `Tooltip`    | `string`       | Hover text on the chart card.                       |

Each `LineSeries` has:

| Field    | Type                 | Description              |
|----------|----------------------|--------------------------|
| `Name`   | `string`             | Series label in legend.  |
| `Points` | `map[string]float64` | Label → value pairs.     |

```go
// Single series
chart := rptgen.NewLineChartSingle("Monthly Revenue", map[string]float64{
    "Jan": 45000, "Feb": 52000, "Mar": 48000,
})

// Multiple series
chart := rptgen.NewLineChart("Revenue vs Costs", []rptgen.LineSeries{
    {Name: "Revenue", Points: map[string]float64{"Q1": 145000, "Q2": 174000}},
    {Name: "Costs",   Points: map[string]float64{"Q1": 95000,  "Q2": 102000}},
})
```

#### `PieChart`

Pie or donut chart.

| Field     | Type                 | Description                              |
|-----------|----------------------|------------------------------------------|
| `Title`   | `string`             | Chart heading.                           |
| `Data`    | `map[string]float64` | Label → value pairs.                     |
| `IsDonut` | `bool`               | Render as donut (hollow center) if `true`. |
| `Tooltip` | `string`             | Hover text on the chart card.            |

```go
chart := rptgen.NewPieChart("Product Mix", map[string]float64{
    "Enterprise":   45,
    "Professional": 30,
    "Starter":      25,
})
chart.IsDonut = true
```

#### `StackedBarChart`

Stacked bar chart where each bar segment represents a named series.

| Field          | Type                  | Description                                    |
|----------------|-----------------------|------------------------------------------------|
| `Title`        | `string`              | Chart heading.                                 |
| `Series`       | `[]StackedBarSeries`  | Ordered slice of categories (preserves order). |
| `IsHorizontal` | `bool`                | Render bars horizontally if `true`.            |
| `Tooltip`      | `string`              | Hover text on the chart card.                  |

Each `StackedBarSeries` has:

| Field      | Type                 | Description                                    |
|------------|----------------------|------------------------------------------------|
| `Category` | `string`             | The bar label (e.g. `"Q1"`).                   |
| `Values`   | `map[string]float64` | Series name → value for each segment.          |

```go
chart := rptgen.NewStackedBarChart("Quarterly Performance", []rptgen.StackedBarSeries{
    {Category: "Q1", Values: map[string]float64{"Revenue": 145000, "Costs": 95000}},
    {Category: "Q2", Values: map[string]float64{"Revenue": 174000, "Costs": 102000}},
})
```

## Theming

`Theme` controls the visual appearance of the rendered report.

```go
theme := rptgen.DefaultTheme()
theme.PrimaryColor   = "#059669" // override specific fields
theme.FontFamily     = "Georgia, serif"

html, err := rptgen.HtmlRenderer{}.Render(report, theme)
```

Pass `nil` as the theme to use `DefaultTheme()` with no overrides.

| Field             | Type       | Default                      | Description                                 |
|-------------------|------------|------------------------------|---------------------------------------------|
| `PrimaryColor`    | `string`   | `#2563eb`                    | Headings, accents.                          |
| `SecondaryColor`  | `string`   | `#64748b`                    | Secondary text and borders.                 |
| `BackgroundColor` | `string`   | `#ffffff`                    | Page background.                            |
| `TextColor`       | `string`   | `#1e293b`                    | Body text.                                  |
| `AccentColor`     | `string`   | `#10b981`                    | Highlight elements.                         |
| `FontFamily`      | `string`   | System UI stack              | CSS `font-family` value.                    |
| `BorderRadius`    | `string`   | `0.5rem`                     | Card corner radius.                         |
| `ChartColors`     | `[]string` | Eight-color palette          | Colors cycled through chart series.         |
| `ShadowIntensity` | `string`   | `"medium"`                   | `"none"`, `"subtle"`, `"medium"`, `"strong"`. |
| `EnableAnimations`| `bool`     | `true`                       | CSS entry animations on cards.              |
| `EnableGradients` | `bool`     | `false`                      | Gradient fills on chart bars.               |

## Renderer

`HtmlRenderer` is the built-in renderer. It produces a fully self-contained HTML document — all CSS and Chart.js JavaScript are embedded inline.

```go
html, err := rptgen.HtmlRenderer{}.Render(report, theme)
```

Custom renderers can be implemented by satisfying the `Renderer` interface:

```go
type Renderer interface {
    Render(report *Report, theme *Theme) (string, error)
}
```
