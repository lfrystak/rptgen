package rptgen

import (
	"encoding/json"
	"fmt"
)

var defaultChartColors = []string{
	"#2563eb", "#10b981", "#f59e0b", "#ef4444",
	"#8b5cf6", "#06b6d4", "#f97316", "#84cc16",
}

func chartColors(theme *Theme) []string {
	if len(theme.ChartColors) > 0 {
		return theme.ChartColors
	}
	return defaultChartColors
}

// chartInitScript wraps a Chart.js config JSON in an IIFE.
// config may be any marshallable value.
//
// SECURITY: json.Marshal HTML-escapes <, >, and & to <, >, & by default,
// making the JSON safe to embed inside a <script> element. Do NOT use json.Encoder with
// SetEscapeHTML(false) here — it would reintroduce XSS via user-controlled label strings.
//
// id MUST be produced by ctx.nextID (output charset [a-z0-9-]), which ensures it is safe
// inside the single-quoted JS string getElementById('...'). Passing an arbitrary string
// directly as id would allow JS injection.
func chartInitScript(id string, config any) (string, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("rptgen: chartInitScript: marshal failed: %w", err)
	}
	return fmt.Sprintf("(function(){var ctx=document.getElementById('%s').getContext('2d');new Chart(ctx,%s);})();",
		id, string(configJSON)), nil
}

// --- Shared Chart.js config model ---
//
// These structs cover the common subset of Chart.js options used by the built-in
// chart types. New chart types should define their own config structs instead of
// extending these, so that adding a chart type requires no modification here.

type chartConfig struct {
	Type    string       `json:"type"`
	Data    chartData    `json:"data"`
	Options chartOptions `json:"options"`
}

type chartData struct {
	Labels   []string       `json:"labels"`
	Datasets []chartDataset `json:"datasets"`
}

type chartDataset struct {
	Label           string   `json:"label,omitempty"`
	Data            []float64 `json:"data"`
	BackgroundColor any      `json:"backgroundColor,omitempty"`
	BorderColor     any      `json:"borderColor,omitempty"`
	BorderWidth     *float64 `json:"borderWidth,omitempty"`
	Fill            *bool    `json:"fill,omitempty"`
	Tension         *float64 `json:"tension,omitempty"`
	PointStyle      *bool    `json:"pointStyle,omitempty"`
}

// chartOptions holds Chart.js options common to the built-in chart types.
// Extra is an optional escape hatch: its key-value pairs are merged into the
// JSON output, allowing new chart types to inject options (e.g. custom scales)
// without modifying this struct.
type chartOptions struct {
	Responsive  bool           `json:"responsive"`
	AspectRatio *float64       `json:"aspectRatio,omitempty"`
	IndexAxis   string         `json:"indexAxis,omitempty"`
	Plugins     *chartPlugins  `json:"plugins,omitempty"`
	Scales      *chartScales   `json:"scales,omitempty"`
	Extra       map[string]any `json:"-"`
}

// MarshalJSON merges Extra fields into the standard JSON output of chartOptions.
func (o chartOptions) MarshalJSON() ([]byte, error) {
	type Alias chartOptions
	raw, err := json.Marshal(Alias(o))
	if err != nil || len(o.Extra) == 0 {
		return raw, err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	for k, v := range o.Extra {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		m[k] = b
	}
	return json.Marshal(m)
}

type chartPlugins struct {
	Legend  *chartLegend  `json:"legend,omitempty"`
	Title   *chartTitle   `json:"title,omitempty"`
	Tooltip *chartTooltip `json:"tooltip,omitempty"`
}

type chartLegend struct {
	Display  bool   `json:"display"`
	Position string `json:"position,omitempty"`
}

type chartTitle struct {
	Display bool   `json:"display"`
	Text    string `json:"text,omitempty"`
}

type chartTooltip struct {
	Enabled bool `json:"enabled"`
}

type chartScales struct {
	X *chartAxis `json:"x,omitempty"`
	Y *chartAxis `json:"y,omitempty"`
}

type chartAxis struct {
	Type    string          `json:"type,omitempty"` // e.g. "linear" for numeric XY axes
	Stacked bool            `json:"stacked,omitempty"`
	Title   *chartAxisTitle `json:"title,omitempty"`
	Min     *float64        `json:"min,omitempty"`
	Max     *float64        `json:"max,omitempty"`
}

type chartAxisTitle struct {
	Display bool   `json:"display"`
	Text    string `json:"text,omitempty"`
}

// applyChartOptions merges user-visible ChartOptions into o.
// title is the chart's title, used when opts.ShowChartTitle is true.
// isCartesian must be true for bar, line, scatter, etc.; false for pie/doughnut.
// isHorizontal must be true when the value axis is X (horizontal bar charts).
func applyChartOptions(o *chartOptions, opts ChartOptions, title string, isCartesian, isHorizontal bool) {
	if opts.AspectRatio != nil {
		o.AspectRatio = opts.AspectRatio
	}

	needPlugins := opts.LegendPosition != "" || opts.ShowChartTitle || (opts.ShowTooltips != nil && !*opts.ShowTooltips)
	if needPlugins && o.Plugins == nil {
		o.Plugins = &chartPlugins{}
	}

	if opts.LegendPosition != "" {
		applyLegendOption(o.Plugins, opts.LegendPosition)
	}
	if opts.ShowChartTitle && title != "" {
		o.Plugins.Title = &chartTitle{Display: true, Text: title}
	}
	if opts.ShowTooltips != nil && !*opts.ShowTooltips {
		o.Plugins.Tooltip = &chartTooltip{Enabled: false}
	}

	if isCartesian {
		applyCartesianOptions(o, opts, isHorizontal)
	}
}

func applyLegendOption(plugins *chartPlugins, position string) {
	if plugins.Legend == nil {
		plugins.Legend = &chartLegend{Display: true}
	}
	if position == "none" {
		plugins.Legend.Display = false
		plugins.Legend.Position = ""
	} else {
		plugins.Legend.Display = true
		plugins.Legend.Position = position
	}
}

func applyCartesianOptions(o *chartOptions, opts ChartOptions, isHorizontal bool) {
	if opts.XAxisTitle != "" {
		ensureScaleX(o).Title = &chartAxisTitle{Display: true, Text: opts.XAxisTitle}
	}
	if opts.YAxisTitle != "" {
		ensureScaleY(o).Title = &chartAxisTitle{Display: true, Text: opts.YAxisTitle}
	}
	if opts.YMin != nil || opts.YMax != nil {
		applyYBounds(o, opts.YMin, opts.YMax, isHorizontal)
	}
}

// applyYBounds applies min/max to the value axis: Y for normal orientation, X for horizontal.
func applyYBounds(o *chartOptions, min, max *float64, isHorizontal bool) {
	if isHorizontal {
		ax := ensureScaleX(o)
		ax.Min = min
		ax.Max = max
	} else {
		ax := ensureScaleY(o)
		ax.Min = min
		ax.Max = max
	}
}

func ensureScaleX(o *chartOptions) *chartAxis {
	if o.Scales == nil {
		o.Scales = &chartScales{}
	}
	if o.Scales.X == nil {
		o.Scales.X = &chartAxis{}
	}
	return o.Scales.X
}

func ensureScaleY(o *chartOptions) *chartAxis {
	if o.Scales == nil {
		o.Scales = &chartScales{}
	}
	if o.Scales.Y == nil {
		o.Scales.Y = &chartAxis{}
	}
	return o.Scales.Y
}
