package rptgen

import (
	"encoding/json"
	"fmt"
	"sort"
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
// config may be any marshallable value. New chart types (scatter, radar, bubble, etc.)
// should define their own config structs and pass them here rather than modifying the
// shared chartConfig/chartDataset/chartOptions structs below.
func chartInitScript(id string, config any) string {
	configJSON, err := json.Marshal(config)
	if err != nil {
		configJSON = []byte("{}")
	}
	return fmt.Sprintf("(function(){var ctx=document.getElementById('%s').getContext('2d');new Chart(ctx,%s);})();",
		id, string(configJSON))
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
	Label           string    `json:"label,omitempty"`
	Data            []float64 `json:"data"`
	BackgroundColor any       `json:"backgroundColor,omitempty"`
	BorderColor     any       `json:"borderColor,omitempty"`
	Fill            *bool     `json:"fill,omitempty"`
	Tension         *float64  `json:"tension,omitempty"`
	PointStyle      *bool     `json:"pointStyle,omitempty"`
}

// chartOptions holds Chart.js options common to the built-in chart types.
// Extra is an optional escape hatch: its key-value pairs are merged into the
// JSON output, allowing new chart types to inject options (e.g. custom scales)
// without modifying this struct.
type chartOptions struct {
	Responsive  bool          `json:"responsive"`
	AspectRatio *float64      `json:"aspectRatio,omitempty"`
	IndexAxis   string        `json:"indexAxis,omitempty"`
	Plugins     *chartPlugins `json:"plugins,omitempty"`
	Scales      *chartScales  `json:"scales,omitempty"`
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
