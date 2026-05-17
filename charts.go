package rptgen

// ChartBase is embedded by all chart types and holds common chart fields.
type ChartBase struct {
	BaseElement
	Title   string
	Tooltip string // optional hover tooltip on the chart card (not on data points)
}

// BarChart displays a single-series bar chart.
type BarChart struct {
	ChartBase
	Data         map[string]float64 // label → value
	IsHorizontal bool
}

func (b *BarChart) ElementType() string { return "BarChart" }

func NewBarChart(title string, data map[string]float64) *BarChart {
	return &BarChart{
		ChartBase: ChartBase{BaseElement: newBaseElement(), Title: title},
		Data:      data,
	}
}

// LineSeries holds a named data series for a LineChart.
type LineSeries struct {
	Name   string
	Points map[string]float64 // label → value
}

// LineChart displays a line chart with one or more series.
// Series is a slice to guarantee deterministic rendering order.
type LineChart struct {
	ChartBase
	Series     []LineSeries
	ShowPoints bool // default: true
}

func (l *LineChart) ElementType() string { return "LineChart" }

// NewLineChart creates a multi-series line chart.
func NewLineChart(title string, series []LineSeries) *LineChart {
	return &LineChart{
		ChartBase:  ChartBase{BaseElement: newBaseElement(), Title: title},
		Series:     series,
		ShowPoints: true,
	}
}

// NewLineChartSingle wraps a single data series, using the chart title as the series name.
func NewLineChartSingle(title string, points map[string]float64) *LineChart {
	return &LineChart{
		ChartBase:  ChartBase{BaseElement: newBaseElement(), Title: title},
		Series:     []LineSeries{{Name: title, Points: points}},
		ShowPoints: true,
	}
}

// PieChart displays a pie or donut chart.
type PieChart struct {
	ChartBase
	Data    map[string]float64 // label → value
	IsDonut bool
}

func (p *PieChart) ElementType() string { return "PieChart" }

func NewPieChart(title string, data map[string]float64) *PieChart {
	return &PieChart{
		ChartBase: ChartBase{BaseElement: newBaseElement(), Title: title},
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
		ChartBase: ChartBase{BaseElement: newBaseElement(), Title: title},
		Series:    series,
	}
}
