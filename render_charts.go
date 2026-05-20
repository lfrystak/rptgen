package rptgen

import "sort"

// RenderHTML implements HTMLRenderer for BarChart.
func (e *BarChart) RenderHTML(ctx *HTMLRenderContext) (string, []string, error) {
	id := ctx.NextID(e.Title)
	script, err := renderBarChartScript(id, e, ctx.Theme)
	if err != nil {
		return "", nil, err
	}
	return RenderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
}

// RenderHTML implements HTMLRenderer for LineChart.
func (e *LineChart) RenderHTML(ctx *HTMLRenderContext) (string, []string, error) {
	id := ctx.NextID(e.Title)
	script, err := renderLineChartScript(id, e, ctx.Theme)
	if err != nil {
		return "", nil, err
	}
	return RenderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
}

// RenderHTML implements HTMLRenderer for PieChart.
func (e *PieChart) RenderHTML(ctx *HTMLRenderContext) (string, []string, error) {
	id := ctx.NextID(e.Title)
	script, err := renderPieChartScript(id, e, ctx.Theme)
	if err != nil {
		return "", nil, err
	}
	return RenderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
}

// RenderHTML implements HTMLRenderer for StackedBarChart.
func (e *StackedBarChart) RenderHTML(ctx *HTMLRenderContext) (string, []string, error) {
	id := ctx.NextID(e.Title)
	script, err := renderStackedBarChartScript(id, e, ctx.Theme)
	if err != nil {
		return "", nil, err
	}
	return RenderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
}

// RenderHTML implements HTMLRenderer for ScatterChart.
// This is the acceptance test for spec 005: a new chart type supplies its own rendering
// by implementing HTMLRenderer — no modification to renderElement required.
func (e *ScatterChart) RenderHTML(ctx *HTMLRenderContext) (string, []string, error) {
	colors := ctx.ChartColors()
	points := make([]scatterPoint, len(e.Points))
	for i, p := range e.Points {
		points[i] = scatterPoint(p)
	}
	ratio := 2.0
	cfg := scatterConfig{
		Type: "scatter",
		Data: scatterData{
			Datasets: []scatterDataset{{
				Label:           e.Title,
				Data:            points,
				BackgroundColor: colors[0],
			}},
		},
		Options: chartOptions{Responsive: true, AspectRatio: &ratio},
	}
	applyChartOptions(&cfg.Options, e.Options, e.Title, true, false)
	id := ctx.NextID(e.Title)
	script, err := ChartInitScript(id, cfg)
	if err != nil {
		return "", nil, err
	}
	return RenderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
}

// scatterConfig is ScatterChart's own Chart.js config. It uses a different data shape
// ({x,y} objects instead of flat float64 slices) and must not reuse chartConfig.
type scatterConfig struct {
	Type    string       `json:"type"`
	Data    scatterData  `json:"data"`
	Options chartOptions `json:"options"`
}

type scatterData struct {
	Datasets []scatterDataset `json:"datasets"`
}

type scatterDataset struct {
	Label           string         `json:"label,omitempty"`
	Data            []scatterPoint `json:"data"`
	BackgroundColor string         `json:"backgroundColor,omitempty"`
}

type scatterPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func renderBarChartScript(id string, e *BarChart, theme *Theme) (string, error) {
	colors := chartColors(theme)
	labels := make([]string, len(e.Data))
	data := make([]float64, len(e.Data))
	bgColors := make([]string, len(e.Data))
	for i, dp := range e.Data {
		labels[i] = dp.Label
		data[i] = dp.Value
		bgColors[i] = colors[i%len(colors)]
	}

	indexAxis := ""
	if e.IsHorizontal {
		indexAxis = "y"
	}

	ratio := 2.0
	cfg := chartConfig{
		Type: "bar",
		Data: chartData{
			Labels:   labels,
			Datasets: []chartDataset{{Data: data, BackgroundColor: bgColors}},
		},
		Options: chartOptions{
			Responsive:  true,
			AspectRatio: &ratio,
			IndexAxis:   indexAxis,
			Plugins:     &chartPlugins{Legend: &chartLegend{Display: false}},
		},
	}
	applyChartOptions(&cfg.Options, e.Options, e.Title, true, e.IsHorizontal)
	return ChartInitScript(id, cfg)
}

func renderLineChartScript(id string, e *LineChart, theme *Theme) (string, error) {
	colors := chartColors(theme)

	// Collect all unique labels in order of first appearance across all series.
	seen := map[string]bool{}
	var labels []string
	for _, s := range e.Series {
		for _, dp := range s.Points {
			if !seen[dp.Label] {
				seen[dp.Label] = true
				labels = append(labels, dp.Label)
			}
		}
	}

	falseVal := false
	tension := 0.4
	datasets := make([]chartDataset, len(e.Series))
	for i, s := range e.Series {
		lookup := make(map[string]float64, len(s.Points))
		for _, dp := range s.Points {
			lookup[dp.Label] = dp.Value
		}
		data := make([]float64, len(labels))
		for j, lbl := range labels {
			data[j] = lookup[lbl]
		}
		color := colors[i%len(colors)]
		ds := chartDataset{
			Label:           s.Name,
			Data:            data,
			BackgroundColor: color,
			BorderColor:     color,
			Fill:            &falseVal,
			Tension:         &tension,
		}
		if !e.ShowPoints {
			pointFalse := false
			ds.PointStyle = &pointFalse
		}
		datasets[i] = ds
	}

	ratio := 2.0
	cfg := chartConfig{
		Type: "line",
		Data: chartData{Labels: labels, Datasets: datasets},
		Options: chartOptions{
			Responsive:  true,
			AspectRatio: &ratio,
			Plugins:     &chartPlugins{Legend: &chartLegend{Display: len(e.Series) > 1}},
		},
	}
	applyChartOptions(&cfg.Options, e.Options, e.Title, true, false)
	return ChartInitScript(id, cfg)
}

func renderPieChartScript(id string, e *PieChart, theme *Theme) (string, error) {
	colors := chartColors(theme)
	labels := make([]string, len(e.Data))
	data := make([]float64, len(e.Data))
	bgColors := make([]string, len(e.Data))
	for i, dp := range e.Data {
		labels[i] = dp.Label
		data[i] = dp.Value
		bgColors[i] = colors[i%len(colors)]
	}

	chartType := "pie"
	if e.IsDonut {
		chartType = "doughnut"
	}

	ratio := 2.0
	cfg := chartConfig{
		Type: chartType,
		Data: chartData{
			Labels:   labels,
			Datasets: []chartDataset{{Data: data, BackgroundColor: bgColors}},
		},
		Options: chartOptions{Responsive: true, AspectRatio: &ratio},
	}
	applyChartOptions(&cfg.Options, e.Options, e.Title, false, false)
	return ChartInitScript(id, cfg)
}

func renderStackedBarChartScript(id string, e *StackedBarChart, theme *Theme) (string, error) {
	colors := chartColors(theme)

	// Collect all series names in order of first appearance.
	seen := map[string]bool{}
	var seriesNames []string
	for _, row := range e.Series {
		for name := range row.Values {
			if !seen[name] {
				seen[name] = true
				seriesNames = append(seriesNames, name)
			}
		}
	}
	sort.Strings(seriesNames)

	labels := make([]string, len(e.Series))
	for i, row := range e.Series {
		labels[i] = row.Category
	}

	datasets := make([]chartDataset, len(seriesNames))
	for si, name := range seriesNames {
		data := make([]float64, len(e.Series))
		for ri, row := range e.Series {
			data[ri] = row.Values[name]
		}
		datasets[si] = chartDataset{
			Label:           name,
			Data:            data,
			BackgroundColor: colors[si%len(colors)],
		}
	}

	stacked := &chartAxis{Stacked: true}
	indexAxis := ""
	if e.IsHorizontal {
		indexAxis = "y"
	}

	ratio := 2.0
	cfg := chartConfig{
		Type: "bar",
		Data: chartData{Labels: labels, Datasets: datasets},
		Options: chartOptions{
			Responsive:  true,
			AspectRatio: &ratio,
			IndexAxis:   indexAxis,
			Scales:      &chartScales{X: stacked, Y: stacked},
		},
	}
	applyChartOptions(&cfg.Options, e.Options, e.Title, true, e.IsHorizontal)
	return ChartInitScript(id, cfg)
}
