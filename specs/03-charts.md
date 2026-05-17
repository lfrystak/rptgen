# Spec 03 — Charts

**File(s) to create:** `charts.go`  
**Depends on:** spec 01 (core types)  
**Referenced by:** spec 04 (HTML renderer), spec 05 (JSON renderer), spec 06 (tests)

---

## Goal

Implement the four chart element types. Each embeds `BaseElement` and `ChartBase`, and implements the `Element` interface. Chart data is held as plain Go structs; the HTML renderer is responsible for serialising it to Chart.js config JSON.

---

## ChartBase

Embedded by all chart types:

```go
type ChartBase struct {
    BaseElement
    Title   string
    Tooltip string  // optional hover tooltip (on the chart card, not on data points)
}
```

---

## BarChart

```go
type BarChart struct {
    ChartBase
    Data         map[string]float64  // label → value
    IsHorizontal bool
}

func (b *BarChart) ElementType() string { return "BarChart" }
```

Constructor:

```go
func NewBarChart(title string, data map[string]float64) *BarChart { ... }
```

> The C# version had overloads for `int` values. In Go, callers can convert `map[string]int` to `map[string]float64` themselves. Keep constructors simple.

---

## LineChart

Supports both single-series and multi-series data.

```go
type LineSeries struct {
    Name   string
    Points map[string]float64  // label → value
}

type LineChart struct {
    ChartBase
    Series     []LineSeries
    ShowPoints bool  // default: true
}

func (l *LineChart) ElementType() string { return "LineChart" }
```

Constructors:

```go
// NewLineChart creates a multi-series line chart.
func NewLineChart(title string, series []LineSeries) *LineChart { ... }

// NewLineChartSingle wraps a single data series using the chart title as the series name.
func NewLineChartSingle(title string, points map[string]float64) *LineChart { ... }
```

> The C# version stored series as `Dictionary<string, Dictionary<string, double>>`. The Go version uses a slice of `LineSeries` structs to preserve series order — map iteration order is not guaranteed in Go.

---

## PieChart

```go
type PieChart struct {
    ChartBase
    Data    map[string]float64  // label → value
    IsDonut bool
}

func (p *PieChart) ElementType() string { return "PieChart" }
```

Constructor:

```go
func NewPieChart(title string, data map[string]float64) *PieChart { ... }
```

---

## StackedBarChart

```go
type StackedBarSeries struct {
    Category string
    Values   map[string]float64  // series name → value
}

type StackedBarChart struct {
    ChartBase
    Series       []StackedBarSeries
    IsHorizontal bool
}

func (s *StackedBarChart) ElementType() string { return "StackedBarChart" }
```

Constructor:

```go
func NewStackedBarChart(title string, series []StackedBarSeries) *StackedBarChart { ... }
```

> The C# version offered a simplified single-series overload (wrapping `map[string]int` into a series named "Value"). In Go, callers can construct `[]StackedBarSeries` directly. No need for multiple constructors.

---

## Chart.js Data Shape Reference

When the HTML renderer serialises chart data to JSON it must produce Chart.js-compatible config. The shapes are:

**BarChart / PieChart:**
```json
{
  "type": "bar",
  "data": {
    "labels": ["A", "B", "C"],
    "datasets": [{ "label": "Title", "data": [1.0, 2.0, 3.0] }]
  }
}
```

**LineChart (multi-series):**
```json
{
  "type": "line",
  "data": {
    "labels": ["Jan", "Feb"],
    "datasets": [
      { "label": "Series1", "data": [10, 20] },
      { "label": "Series2", "data": [5, 15] }
    ]
  }
}
```

**StackedBarChart:**
```json
{
  "type": "bar",
  "data": {
    "labels": ["Cat1", "Cat2"],
    "datasets": [
      { "label": "SeriesA", "data": [10, 20] },
      { "label": "SeriesB", "data": [5, 15] }
    ]
  },
  "options": { "scales": { "x": { "stacked": true }, "y": { "stacked": true } } }
}
```

> The renderer (spec 04) owns the full Chart.js config generation. This spec only defines the Go data types. Document these shapes here so the renderer spec author has a reference.

---

## Acceptance Criteria

- All four chart types compile and implement `Element`
- `ElementType()` returns exact strings: `"BarChart"`, `"LineChart"`, `"PieChart"`, `"StackedBarChart"`
- Series/data ordering is deterministic: `LineSeries` and `StackedBarSeries` use slices, not maps
- `go build ./...` and `go vet ./...` pass
