package rptgen

import (
	"encoding/json"
	"fmt"
	"html"
	"sort"
	"strings"
)

var defaultChartColors = []string{
	"#2563eb", "#10b981", "#f59e0b", "#ef4444",
	"#8b5cf6", "#06b6d4", "#f97316", "#84cc16",
}

// HtmlRenderer renders a Report to a self-contained HTML document string.
type HtmlRenderer struct{}

// Render converts report to a complete HTML document using the provided theme.
// If theme is nil, DefaultTheme() is used.
func (h HtmlRenderer) Render(report *Report, theme *Theme) (string, error) {
	if theme == nil {
		theme = DefaultTheme()
	}

	var chartScripts []string
	var b strings.Builder

	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString("<html lang=\"en\">\n<head>\n")
	b.WriteString("  <meta charset=\"UTF-8\">\n")
	b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	fmt.Fprintf(&b, "  <title>%s</title>\n", html.EscapeString(report.Title))
	b.WriteString("  <style>\n")
	b.WriteString(generateCSS(theme))
	b.WriteString("  </style>\n")
	b.WriteString("</head>\n<body>\n")
	b.WriteString("  <div class=\"report\">\n")

	// Header
	b.WriteString("    <header class=\"report-header\">\n")
	if report.LogoURL != "" {
		fmt.Fprintf(&b, "      <img src=\"%s\" alt=\"logo\">\n", html.EscapeString(report.LogoURL))
	}
	fmt.Fprintf(&b, "      <h1 class=\"report-title\">%s</h1>\n", html.EscapeString(report.Title))
	if !report.GeneratedAt.IsZero() {
		fmt.Fprintf(&b, "      <span class=\"report-generated\">Generated: %s</span>\n",
			html.EscapeString(report.GeneratedAt.Format("2006-01-02 15:04:05")))
	}
	b.WriteString("    </header>\n\n")

	// Sections
	for _, section := range report.Sections {
		colTemplate := columnWidthsToCSS(section.ColumnWidths)
		b.WriteString("    <section class=\"report-section\">\n")
		if section.Title != "" {
			fmt.Fprintf(&b, "      <h2 class=\"section-title\">%s</h2>\n", html.EscapeString(section.Title))
		}
		fmt.Fprintf(&b, "      <div class=\"section-grid\" style=\"--col-template: %s\">\n", colTemplate)
		for _, elem := range section.Elements {
			b.WriteString("        <div class=\"element-wrapper\">\n")
			rendered, scripts := renderElement(elem, report, theme)
			b.WriteString(rendered)
			chartScripts = append(chartScripts, scripts...)
			b.WriteString("        </div>\n")
		}
		b.WriteString("      </div>\n")
		b.WriteString("    </section>\n\n")
	}

	// Footer
	if report.Footer != "" {
		fmt.Fprintf(&b, "    <footer class=\"report-footer\">%s</footer>\n", html.EscapeString(report.Footer))
	}

	b.WriteString("  </div>\n\n")

	b.WriteString("  <script src=\"https://cdn.jsdelivr.net/npm/chart.js\"></script>\n")

	// Chart init scripts
	if len(chartScripts) > 0 {
		b.WriteString("  <script>\n")
		for _, s := range chartScripts {
			b.WriteString(s)
			b.WriteByte('\n')
		}
		b.WriteString("  </script>\n")
	}

	b.WriteString("</body>\n</html>")
	return b.String(), nil
}

// columnWidthsToCSS converts a slice of proportional widths to a CSS grid-template-columns value.
func columnWidthsToCSS(widths []int) string {
	if len(widths) == 0 {
		return "1fr"
	}
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = fmt.Sprintf("%dfr", w)
	}
	return strings.Join(parts, " ")
}

// renderElement dispatches to the correct element renderer and returns HTML + any chart init scripts.
func renderElement(elem Element, report *Report, theme *Theme) (string, []string) {
	switch e := elem.(type) {
	case *NumberTile:
		return renderNumberTile(e, report), nil
	case *DateTile:
		return renderDateTile(e), nil
	case *FreeText:
		return renderFreeText(e), nil
	case *Table:
		return renderTable(e), nil
	case *Canvas:
		return renderCanvas(e, report, theme)
	case *BarChart:
		script := renderBarChartScript(e, theme)
		return renderChartContainer(e.elementID(), e.Title, e.Tooltip), []string{script}
	case *LineChart:
		script := renderLineChartScript(e, theme)
		return renderChartContainer(e.elementID(), e.Title, e.Tooltip), []string{script}
	case *PieChart:
		script := renderPieChartScript(e, theme)
		return renderChartContainer(e.elementID(), e.Title, e.Tooltip), []string{script}
	case *StackedBarChart:
		script := renderStackedBarChartScript(e, theme)
		return renderChartContainer(e.elementID(), e.Title, e.Tooltip), []string{script}
	default:
		return fmt.Sprintf("<div><!-- unknown element: %s --></div>\n", html.EscapeString(elem.ElementType())), nil
	}
}

func renderNumberTile(e *NumberTile, report *Report) string {
	var b strings.Builder
	b.WriteString("          <div class=\"element tile number-tile\">\n")
	b.WriteString(tooltipIcon(e.Tooltip))
	fmt.Fprintf(&b, "            <div class=\"tile-title\">%s</div>\n", html.EscapeString(e.Title))
	fmt.Fprintf(&b, "            <div class=\"tile-value\">%s</div>\n", html.EscapeString(e.FormatValue(report.Locale)))
	if e.Subtitle != "" {
		fmt.Fprintf(&b, "            <div class=\"tile-subtitle\">%s</div>\n", html.EscapeString(e.Subtitle))
	}
	b.WriteString("          </div>\n")
	return b.String()
}

func renderDateTile(e *DateTile) string {
	var b strings.Builder
	b.WriteString("          <div class=\"element tile date-tile\">\n")
	b.WriteString(tooltipIcon(e.Tooltip))
	fmt.Fprintf(&b, "            <div class=\"tile-title\">%s</div>\n", html.EscapeString(e.Title))
	fmt.Fprintf(&b, "            <div class=\"tile-value\">%s</div>\n", html.EscapeString(e.FormatValue()))
	if e.Subtitle != "" {
		fmt.Fprintf(&b, "            <div class=\"tile-subtitle\">%s</div>\n", html.EscapeString(e.Subtitle))
	}
	b.WriteString("          </div>\n")
	return b.String()
}

func renderFreeText(e *FreeText) string {
	var b strings.Builder
	b.WriteString("          <div class=\"element free-text\">\n")
	if e.IsHTML {
		b.WriteString("            ")
		b.WriteString(e.Content)
		b.WriteByte('\n')
	} else {
		fmt.Fprintf(&b, "            <p>%s</p>\n", html.EscapeString(e.Content))
	}
	b.WriteString("          </div>\n")
	return b.String()
}

func renderTable(e *Table) string {
	var b strings.Builder
	b.WriteString("          <div class=\"element table-wrapper\">\n")
	if e.Title != "" {
		fmt.Fprintf(&b, "            <h3 class=\"element-title\">%s</h3>\n", html.EscapeString(e.Title))
	}
	b.WriteString("            <table>\n              <thead><tr>")
	for _, col := range e.Columns {
		fmt.Fprintf(&b, "<th>%s</th>", html.EscapeString(col))
	}
	b.WriteString("</tr></thead>\n              <tbody>\n")
	for _, row := range e.Rows {
		b.WriteString("                <tr>")
		for _, col := range e.Columns {
			val := ""
			if v, ok := row[col]; ok {
				val = fmt.Sprint(v)
			}
			fmt.Fprintf(&b, "<td>%s</td>", html.EscapeString(val))
		}
		b.WriteString("</tr>\n")
	}
	b.WriteString("              </tbody>\n            </table>\n          </div>\n")
	return b.String()
}

func renderCanvas(e *Canvas, report *Report, theme *Theme) (string, []string) {
	colTemplate := columnWidthsToCSS(e.ColumnWidths)
	var b strings.Builder
	var allScripts []string
	fmt.Fprintf(&b, "          <div class=\"element canvas-grid\" style=\"--col-template: %s\">\n", colTemplate)
	for _, child := range e.Elements {
		b.WriteString("            <div class=\"element-wrapper\">\n")
		rendered, scripts := renderElement(child, report, theme)
		b.WriteString(rendered)
		allScripts = append(allScripts, scripts...)
		b.WriteString("            </div>\n")
	}
	b.WriteString("          </div>\n")
	return b.String(), allScripts
}

func renderChartContainer(id, title, tooltip string) string {
	var b strings.Builder
	b.WriteString("          <div class=\"element chart-container\">\n")
	b.WriteString(tooltipIcon(tooltip))
	if title != "" {
		fmt.Fprintf(&b, "            <h3 class=\"element-title\">%s</h3>\n", html.EscapeString(title))
	}
	fmt.Fprintf(&b, "            <canvas id=\"%s\"></canvas>\n", html.EscapeString(id))
	b.WriteString("          </div>\n")
	return b.String()
}

func tooltipIcon(tooltip string) string {
	if tooltip == "" {
		return ""
	}
	return fmt.Sprintf("            <span class=\"tooltip-icon\" data-tooltip=\"%s\">&#9432;</span>\n", html.EscapeString(tooltip))
}

func chartColors(theme *Theme) []string {
	if len(theme.ChartColors) > 0 {
		return theme.ChartColors
	}
	return defaultChartColors
}

// sortedKeys returns map keys in sorted order for deterministic chart rendering.
func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// chartInitScript wraps a Chart.js config JSON in an IIFE.
func chartInitScript(id string, config any) string {
	configJSON, err := json.Marshal(config)
	if err != nil {
		configJSON = []byte("{}")
	}
	return fmt.Sprintf("(function(){var ctx=document.getElementById('%s').getContext('2d');new Chart(ctx,%s);})();",
		id, string(configJSON))
}

// --- Chart config structs ---

type chartConfig struct {
	Type    string      `json:"type"`
	Data    chartData   `json:"data"`
	Options chartOptions `json:"options"`
}

type chartData struct {
	Labels   []string       `json:"labels"`
	Datasets []chartDataset `json:"datasets"`
}

type chartDataset struct {
	Label           string    `json:"label,omitempty"`
	Data            []float64 `json:"data"`
	BackgroundColor any       `json:"backgroundColor,omitempty"`
	BorderColor     any       `json:"borderColor,omitempty"`
	Fill            *bool     `json:"fill,omitempty"`
	Tension         *float64  `json:"tension,omitempty"`
	PointStyle      *bool     `json:"pointStyle,omitempty"`
}

type chartOptions struct {
	Responsive bool             `json:"responsive"`
	IndexAxis  string           `json:"indexAxis,omitempty"`
	Plugins    *chartPlugins    `json:"plugins,omitempty"`
	Scales     *chartScales     `json:"scales,omitempty"`
}

type chartPlugins struct {
	Legend *chartLegend `json:"legend,omitempty"`
}

type chartLegend struct {
	Display bool `json:"display"`
}

type chartScales struct {
	X *chartAxis `json:"x,omitempty"`
	Y *chartAxis `json:"y,omitempty"`
}

type chartAxis struct {
	Stacked bool `json:"stacked,omitempty"`
}

func renderBarChartScript(e *BarChart, theme *Theme) string {
	colors := chartColors(theme)
	labels := sortedKeys(e.Data)
	data := make([]float64, len(labels))
	bgColors := make([]string, len(labels))
	for i, lbl := range labels {
		data[i] = e.Data[lbl]
		bgColors[i] = colors[i%len(colors)]
	}

	chartType := "bar"
	indexAxis := ""
	if e.IsHorizontal {
		indexAxis = "y"
	}

	cfg := chartConfig{
		Type: chartType,
		Data: chartData{
			Labels:   labels,
			Datasets: []chartDataset{{Data: data, BackgroundColor: bgColors}},
		},
		Options: chartOptions{
			Responsive: true,
			IndexAxis:  indexAxis,
		},
	}
	return chartInitScript(e.elementID(), cfg)
}

func renderLineChartScript(e *LineChart, theme *Theme) string {
	colors := chartColors(theme)

	// Collect all unique labels in order of first appearance across all series.
	seen := map[string]bool{}
	var labels []string
	for _, s := range e.Series {
		for _, lbl := range sortedKeys(s.Points) {
			if !seen[lbl] {
				seen[lbl] = true
				labels = append(labels, lbl)
			}
		}
	}

	falseVal := false
	tension := 0.4
	datasets := make([]chartDataset, len(e.Series))
	for i, s := range e.Series {
		data := make([]float64, len(labels))
		for j, lbl := range labels {
			data[j] = s.Points[lbl]
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

	cfg := chartConfig{
		Type: "line",
		Data: chartData{Labels: labels, Datasets: datasets},
		Options: chartOptions{
			Responsive: true,
			Plugins:    &chartPlugins{Legend: &chartLegend{Display: len(e.Series) > 1}},
		},
	}
	return chartInitScript(e.elementID(), cfg)
}

func renderPieChartScript(e *PieChart, theme *Theme) string {
	colors := chartColors(theme)
	labels := sortedKeys(e.Data)
	data := make([]float64, len(labels))
	bgColors := make([]string, len(labels))
	for i, lbl := range labels {
		data[i] = e.Data[lbl]
		bgColors[i] = colors[i%len(colors)]
	}

	chartType := "pie"
	if e.IsDonut {
		chartType = "doughnut"
	}

	cfg := chartConfig{
		Type: chartType,
		Data: chartData{
			Labels:   labels,
			Datasets: []chartDataset{{Data: data, BackgroundColor: bgColors}},
		},
		Options: chartOptions{Responsive: true},
	}
	return chartInitScript(e.elementID(), cfg)
}

func renderStackedBarChartScript(e *StackedBarChart, theme *Theme) string {
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

	cfg := chartConfig{
		Type: "bar",
		Data: chartData{Labels: labels, Datasets: datasets},
		Options: chartOptions{
			Responsive: true,
			IndexAxis:  indexAxis,
			Scales:     &chartScales{X: stacked, Y: stacked},
		},
	}
	return chartInitScript(e.elementID(), cfg)
}

// generateCSS produces the full CSS for the report, parameterized by Theme.
func generateCSS(t *Theme) string {
	shadow := shadowCSS(t.ShadowIntensity)

	var b strings.Builder

	fmt.Fprintf(&b, `
*, *::before, *::after { box-sizing: border-box; }
body {
  font-family: %s;
  background-color: %s;
  color: %s;
  margin: 0;
  padding: 1rem;
}
.report {
  max-width: 1400px;
  margin: 0 auto;
  padding: 1rem;
}
`, t.FontFamily, t.BackgroundColor, t.TextColor)

	// Header
	if t.EnableGradients {
		fmt.Fprintf(&b, `.report-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 0;
  margin-bottom: 2rem;
  border-bottom: 2px solid %s;
  background: linear-gradient(135deg, %s15 0%%, transparent 100%%);
}
`, t.PrimaryColor, t.PrimaryColor)
	} else {
		fmt.Fprintf(&b, `.report-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 0;
  margin-bottom: 2rem;
  border-bottom: 2px solid %s;
}
`, t.PrimaryColor)
	}

	fmt.Fprintf(&b, `.report-header img { max-height: 48px; }
.report-title { margin: 0; color: %s; font-size: 1.5rem; }
.report-generated { font-size: 0.8rem; color: %s; }
.report-footer {
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid %s;
  font-size: 0.85rem;
  color: %s;
}
`, t.PrimaryColor, t.SecondaryColor, t.SecondaryColor, t.SecondaryColor)

	// Section grid
	fmt.Fprintf(&b, `.report-section { margin-bottom: 2rem; }
.section-title {
  font-size: 1.1rem;
  color: %s;
  margin-bottom: 1rem;
  border-left: 4px solid %s;
  padding-left: 0.5rem;
}
.section-grid {
  display: grid;
  grid-template-columns: var(--col-template, 1fr);
  gap: 1rem;
}
.canvas-grid {
  display: grid;
  grid-template-columns: var(--col-template, 1fr);
  gap: 0.75rem;
  padding: 0;
}
`, t.TextColor, t.PrimaryColor)

	// Element wrapper / card
	fmt.Fprintf(&b, `.element-wrapper {
  background: %s;
  border-radius: %s;
  box-shadow: %s;
  position: relative;
}
`, t.CardColor, t.BorderRadius, shadow)

	// Tiles
	fmt.Fprintf(&b, `.tile {
  padding: 1.5rem;
  padding-top: 2rem;
  text-align: center;
  position: relative;
}
.tile-title {
  font-size: 0.85rem;
  color: %s;
  margin-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.tile-value {
  font-size: 2rem;
  font-weight: 700;
  color: %s;
  line-height: 1.2;
}
.tile-subtitle {
  font-size: 0.8rem;
  color: %s;
  margin-top: 0.25rem;
}
`, t.SecondaryColor, t.PrimaryColor, t.SecondaryColor)

	// Chart
	b.WriteString(`.chart-container {
  padding: 1rem;
  position: relative;
  height: 350px;
}
.element-title {
  font-size: 1rem;
  margin: 0 0 0.75rem;
}
`)

	// Table
	fmt.Fprintf(&b, `.table-wrapper {
  overflow-x: auto;
  padding: 1rem;
}
table {
  width: 100%%;
  border-collapse: collapse;
  font-size: 0.9rem;
}
th {
  background-color: %s;
  color: #ffffff;
  padding: 0.6rem 1rem;
  text-align: left;
  font-weight: 600;
}
td {
  padding: 0.5rem 1rem;
  border-bottom: 1px solid %s20;
}
tbody tr:nth-child(even) td {
  background-color: %s08;
}
tbody tr:hover td {
  background-color: %s15;
}
`, t.PrimaryColor, t.TextColor, t.TextColor, t.PrimaryColor)

	// Free text
	b.WriteString(`.free-text {
  padding: 1rem 1.5rem;
  line-height: 1.6;
}
.free-text p { margin: 0; }
`)

	// Tooltip icon
	b.WriteString(`.tooltip-icon {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  width: 1.1rem;
  height: 1.1rem;
  font-size: 1.1rem;
  line-height: 1;
  color: #9ca3af;
  cursor: help;
  z-index: 10;
}
.tooltip-icon:hover { color: #6b7280; }
.tooltip-icon::after {
  content: attr(data-tooltip);
  display: none;
  position: absolute;
  bottom: calc(100% + 6px);
  right: 0;
  min-width: 180px;
  max-width: 280px;
  background: rgba(0,0,0,0.85);
  color: #fff;
  padding: 0.4rem 0.6rem;
  border-radius: 4px;
  font-size: 0.78rem;
  white-space: normal;
  line-height: 1.4;
  z-index: 200;
  pointer-events: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.3);
}
.tooltip-icon:hover::after { display: block; }
`)

	// Animations
	if t.EnableAnimations {
		b.WriteString(`@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to   { opacity: 1; transform: translateY(0); }
}
.element-wrapper {
  animation: fadeIn 0.3s ease;
}
`)
	}

	// Responsive
	b.WriteString(`@media (max-width: 768px) {
  .section-grid, .canvas-grid { grid-template-columns: 1fr !important; }
}
`)

	// Print
	b.WriteString(`@media print {
  .element-wrapper { box-shadow: none !important; animation: none !important; }
  .section-grid, .canvas-grid { grid-template-columns: 1fr !important; }
}
`)

	return b.String()
}

func shadowCSS(intensity string) string {
	switch intensity {
	case "none":
		return "none"
	case "subtle":
		return "0 1px 2px rgba(0,0,0,0.05)"
	case "strong":
		return "0 10px 15px rgba(0,0,0,0.1)"
	default: // "medium" or empty
		return "0 4px 6px rgba(0,0,0,0.07)"
	}
}
