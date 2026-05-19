package rptgen

import "sort"

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
	return chartInitScript(id, cfg)
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
	return chartInitScript(id, cfg)
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
	return chartInitScript(id, cfg)
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
	return chartInitScript(id, cfg)
}
