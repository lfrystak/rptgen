package rptgen

import "sort"

// DataPoint is an ordered label-value pair used in single-series chart data.
// Use a []DataPoint literal to preserve the insertion order that will appear on
// the chart axis. If alphabetical order is acceptable, use DataPointsFromMap.
type DataPoint struct {
	Label string
	Value float64
}

// DataPointsFromMap converts a map to a []DataPoint sorted alphabetically by
// label. Use this only when alphabetical order is acceptable; prefer a
// []DataPoint literal when the display order matters (e.g. months, quarters).
func DataPointsFromMap(m map[string]float64) []DataPoint {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pts := make([]DataPoint, len(keys))
	for i, k := range keys {
		pts[i] = DataPoint{Label: k, Value: m[k]}
	}
	return pts
}

// ChartOptions holds common Chart.js options shared by every chart type.
// Zero values preserve Chart.js defaults — a chart with all-zero ChartOptions
// produces identical output to one with no Options field set.
//
// Axis-related fields (XAxisTitle, YAxisTitle, YMin, YMax) are silently
// ignored for non-cartesian charts (PieChart, any future polar/radial type).
// For horizontal bar charts, YMin/YMax are applied to the value axis (X in
// Chart.js terms); XAxisTitle and YAxisTitle always map to their literal axes.
type ChartOptions struct {
	// LegendPosition sets the legend placement.
	// Accepted values: "top", "bottom", "left", "right".
	// Use "none" to hide the legend entirely.
	// Empty string (zero value) preserves the per-chart-type default.
	LegendPosition string

	// XAxisTitle sets the label for the X axis (cartesian charts only).
	XAxisTitle string

	// YAxisTitle sets the label for the Y axis (cartesian charts only).
	// For horizontal bar charts this is the category axis; the value axis is X.
	YAxisTitle string

	// YMin and YMax clamp the value-axis range (cartesian charts only).
	// For horizontal bar charts the value axis is X.
	YMin *float64
	YMax *float64

	// ShowTooltips controls whether Chart.js data-point tooltips are shown.
	// Nil preserves the Chart.js default (enabled).
	// Set to a pointer to false to disable: opts.ShowTooltips = rptgen.Ptr(false)
	ShowTooltips *bool

	// AspectRatio overrides the chart's width-to-height ratio.
	// Nil uses the per-chart-type default (2.0 for most built-in types).
	AspectRatio *float64

	// ShowChartTitle, when true, additionally emits a Chart.js native title
	// block from ChartBase.Title, appearing inside the canvas area. The HTML
	// <h3> heading above the chart card is always rendered regardless.
	ShowChartTitle bool
}

// ChartBase is embedded by all chart types and holds common chart fields.
type ChartBase struct {
	Title   string
	Tooltip string       // optional hover tooltip on the chart card (not on data points)
	Options ChartOptions // optional Chart.js display options; zero value = Chart.js defaults
}

// BarChart displays a single-series bar chart.
// Data is a slice to preserve the caller's insertion order on the chart axis.
type BarChart struct {
	ChartBase
	Data         []DataPoint
	IsHorizontal bool
}

func (b *BarChart) ElementType() string { return "BarChart" }

func NewBarChart(title string, data []DataPoint) *BarChart {
	return &BarChart{
		ChartBase: ChartBase{Title: title},
		Data:      data,
	}
}

// LineSeries holds a named data series for a LineChart in categorical mode.
// Points is a slice to preserve the caller's insertion order on the chart axis.
type LineSeries struct {
	Name   string
	Points []DataPoint
}

// XYLineSeries holds a named data series of numeric X/Y points for a LineChart
// in XY mode (created with [NewLineChartXY]). Unlike [LineSeries], the X axis
// is linear so point spacing reflects actual X values rather than category order.
type XYLineSeries struct {
	Name   string
	Points []XYPoint
}

// LineChart displays a line chart with one or more series.
//
// A LineChart operates in one of two mutually exclusive modes determined by the
// constructor used:
//
//   - Categorical mode ([NewLineChart], [NewLineChartSingle]): X axis labels are
//     strings. Points are spaced evenly regardless of their numeric value. Use
//     this for named categories such as months, quarters, or product names.
//
//   - XY mode ([NewLineChartXY]): X axis is numeric (linear). Points are
//     positioned proportionally to their X value, so uneven sampling is
//     rendered correctly. Use this for mathematical functions, sensor readings,
//     or any data where X is a continuous number.
//
// Both modes honour all [ChartOptions] (axis titles, legend, tooltips, aspect
// ratio, etc.) and emit a Chart.js "line" chart under the hood.
// [ScatterChart] remains the right choice when you want individual data points
// without a connecting line.
type LineChart struct {
	ChartBase
	Series     []LineSeries   // categorical mode: string labels on X axis
	XYSeries   []XYLineSeries // XY mode: numeric X axis; mutually exclusive with Series
	ShowPoints bool           // default: true
}

func (l *LineChart) ElementType() string { return "LineChart" }

// NewLineChart creates a multi-series line chart with categorical string labels on the X axis.
func NewLineChart(title string, series []LineSeries) *LineChart {
	return &LineChart{
		ChartBase:  ChartBase{Title: title},
		Series:     series,
		ShowPoints: true,
	}
}

// NewLineChartSingle creates a single-series line chart with categorical string labels on the X axis.
// The chart title is used as the series name.
func NewLineChartSingle(title string, points []DataPoint) *LineChart {
	return &LineChart{
		ChartBase:  ChartBase{Title: title},
		Series:     []LineSeries{{Name: title, Points: points}},
		ShowPoints: true,
	}
}

// NewLineChartXY creates a single-series line chart with a numeric (linear) X axis.
// Use this instead of [NewLineChartSingle] when X values are floats rather than
// category labels — for example when plotting mathematical functions or continuous
// measurements. Points are positioned proportionally to their X value, so uneven
// sampling is rendered correctly.
//
// Internally the chart uses Chart.js type "line" with scales.x.type "linear".
// This is the same engine as the categorical LineChart but with a numeric axis.
// [ScatterChart] (dots only, no connecting line) remains the right choice when
// you want to show individual data points without a line.
func NewLineChartXY(title string, points []XYPoint) *LineChart {
	return &LineChart{
		ChartBase:  ChartBase{Title: title},
		XYSeries:   []XYLineSeries{{Name: title, Points: points}},
		ShowPoints: true,
	}
}

// PieChart displays a pie or donut chart.
// Data is a slice to preserve the caller's insertion order on the chart legend.
type PieChart struct {
	ChartBase
	Data    []DataPoint
	IsDonut bool
}

func (p *PieChart) ElementType() string { return "PieChart" }

func NewPieChart(title string, data []DataPoint) *PieChart {
	return &PieChart{
		ChartBase: ChartBase{Title: title},
		Data:      data,
	}
}

// StackedBarSeries holds one category's values across all series in a StackedBarChart.
type StackedBarSeries struct {
	Category string
	Values   map[string]float64 // series name → value
}

// StackedBarChart displays a stacked bar chart.
// Series is a slice to guarantee deterministic category order.
type StackedBarChart struct {
	ChartBase
	Series       []StackedBarSeries
	IsHorizontal bool
}

func (s *StackedBarChart) ElementType() string { return "StackedBarChart" }

func NewStackedBarChart(title string, series []StackedBarSeries) *StackedBarChart {
	return &StackedBarChart{
		ChartBase: ChartBase{Title: title},
		Series:    series,
	}
}

// XYPoint is a numeric X/Y coordinate pair used by [ScatterChart] and the XY mode of [LineChart].
type XYPoint struct {
	X float64
	Y float64
}

// ScatterChart displays a scatter plot. Each point is an independent {x, y} coordinate;
// unlike bar/line charts there are no shared axis labels.
//
// ScatterChart is also the spec-005 acceptance-test element: it adds a new chart type
// by implementing HTMLRenderer on its own struct without touching renderElement.
type ScatterChart struct {
	ChartBase
	Points []XYPoint
}

func (s *ScatterChart) ElementType() string { return "ScatterChart" }

// NewScatterChart creates a scatter chart with the given X/Y data points.
func NewScatterChart(title string, points []XYPoint) *ScatterChart {
	return &ScatterChart{
		ChartBase: ChartBase{Title: title},
		Points:    points,
	}
}
