package rptgen

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// updateGolden regenerates testdata/*.html when passed as -update.
var updateGolden = flag.Bool("update", false, "regenerate golden test files in testdata/")

// checkGolden compares got against the committed golden file, or writes it when -update is set.
func checkGolden(t *testing.T, name, got string) {
	t.Helper()
	const dir = "testdata"
	path := filepath.Join(dir, name)
	if *updateGolden {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		t.Logf("updated %s", path)
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s (run with -update to create): %v", path, err)
	}
	if string(want) != got {
		wl := strings.Split(string(want), "\n")
		gl := strings.Split(got, "\n")
		for i := 0; i < len(gl) && i < len(wl); i++ {
			if gl[i] != wl[i] {
				t.Errorf("golden %s: first diff at line %d:\ngot:  %q\nwant: %q", name, i+1, gl[i], wl[i])
				return
			}
		}
		t.Errorf("golden %s: line count differs (got %d, want %d)", name, len(gl), len(wl))
	}
}

// buildDeterministicReport returns buildFullReport() with a fixed timestamp so golden tests are stable.
func buildDeterministicReport() *Report {
	r := buildFullReport()
	r.GeneratedAt = time.Time{} // suppress non-deterministic timestamp
	return r
}

// parseChartScript extracts and parses the Chart.js config from a chartInitScript output.
func parseChartScript(t *testing.T, script string) chartConfig {
	t.Helper()
	const open = "new Chart(ctx,"
	i := strings.Index(script, open)
	if i < 0 {
		t.Fatalf("parseChartScript: 'new Chart(ctx,' not found in: %s", script)
	}
	rest := script[i+len(open):]
	j := strings.LastIndex(rest, ");})();")
	if j < 0 {
		t.Fatalf("parseChartScript: closing ');})();' not found")
	}
	var cfg chartConfig
	if err := json.Unmarshal([]byte(rest[:j]), &cfg); err != nil {
		t.Fatalf("parseChartScript: %v — JSON: %s", err, rest[:j])
	}
	return cfg
}

// customTestElement is a dummy Element used to exercise the unknown-type branch in renderElement.
type customTestElement struct{ BaseElement }

func (c *customTestElement) ElementType() string { return "CustomTestElement" }

func buildFullReport() *Report {
	r := NewReport("Smoke Test Report")

	section := &Section{Title: "All Elements"}

	section.AddElement(&NumberTile{
		BaseElement:  newBaseElement(),
		Title:        "Revenue",
		Value:        99999.99,
		Format:       "%.2f",
		Prefix:       "$",
		ThousandsSep: true,
	})
	section.AddElement(&DateTile{
		BaseElement: newBaseElement(),
		Title:       "As Of",
		Value:       time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	})
	section.AddElement(&FreeText{
		BaseElement: newBaseElement(),
		Content:     "Hello world",
		IsHTML:      false,
	})
	section.AddElement(NewTable("Data", []map[string]any{
		{"col1": "a", "col2": "b"},
	}))
	section.AddElement(NewBarChart("Bar", map[string]float64{"X": 1, "Y": 2}))
	section.AddElement(NewLineChartSingle("Line", map[string]float64{"Jan": 10, "Feb": 20}))
	section.AddElement(NewPieChart("Pie", map[string]float64{"A": 30, "B": 70}))
	section.AddElement(NewStackedBarChart("Stacked", []StackedBarSeries{
		{Category: "Q1", Values: map[string]float64{"S1": 10, "S2": 20}},
	}))

	canvas := NewCanvas(1, 1)
	canvas.AddElement(&NumberTile{
		BaseElement: newBaseElement(),
		Title:       "Inside Canvas",
		Value:       42,
	})
	section.AddElement(canvas)

	r.AddSection(section)
	return r
}

func TestHtmlSmokeTest(t *testing.T) {
	r := buildFullReport()
	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.HasPrefix(out, "<!DOCTYPE html>") {
		t.Error("output must start with <!DOCTYPE html>")
	}
	if !strings.Contains(out, r.Title) {
		t.Error("output must contain the report title")
	}
	if !strings.Contains(out, "new Chart(") {
		t.Error("output must contain 'new Chart('")
	}
	if !strings.Contains(out, "cdn.jsdelivr.net") {
		t.Error("output must reference cdn.jsdelivr.net for Chart.js")
	}
}

func TestHtmlThemeApplication(t *testing.T) {
	r := NewReport("Theme Test")
	r.AddSection(&Section{})

	th := DefaultTheme()
	th.PrimaryColor = "#ff0000"

	out, err := HtmlRenderer{}.Render(r, th)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(out, "#ff0000") {
		t.Error("output must contain the custom primary color #ff0000")
	}
}

func TestHtmlColumnWidths(t *testing.T) {
	r := NewReport("Col Test")
	r.AddSection(&Section{
		ColumnWidths: []int{1, 2},
	})

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(out, "1fr 2fr") {
		t.Error("output must contain '1fr 2fr' for ColumnWidths [1,2]")
	}
}

func TestHtmlEscapingFreeTextNotHTML(t *testing.T) {
	r := NewReport("Escape Test")
	section := &Section{}
	section.AddElement(&FreeText{
		BaseElement: newBaseElement(),
		Content:     "<script>alert(1)</script>",
		IsHTML:      false,
	})
	r.AddSection(section)

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Error("output must contain escaped &lt;script&gt;")
	}
	if strings.Contains(out, "<script>alert(1)</script>") {
		t.Error("output must not contain raw unescaped <script> tag")
	}
}

func TestHtmlEscapingFreeTextIsHTML(t *testing.T) {
	r := NewReport("HTML Test")
	section := &Section{}
	section.AddElement(&FreeText{
		BaseElement: newBaseElement(),
		Content:     "<b>bold</b>",
		IsHTML:      true,
	})
	r.AddSection(section)

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(out, "<b>bold</b>") {
		t.Error("output must contain verbatim <b>bold</b> for IsHTML: true")
	}
}

func TestHtmlNilThemeDoesNotPanic(t *testing.T) {
	r := buildFullReport()
	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(out, DefaultTheme().PrimaryColor) {
		t.Error("nil theme must apply default primary color")
	}
}

// --- golden file tests ---

func TestGoldenFullReport(t *testing.T) {
	out, err := HtmlRenderer{}.Render(buildDeterministicReport(), nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	checkGolden(t, "full_report.html", out)
}

func TestGoldenCustomTheme(t *testing.T) {
	th := DefaultTheme()
	th.PrimaryColor = "#7c3aed"
	th.EnableGradients = true
	th.EnableAnimations = false
	th.ShadowIntensity = "strong"
	out, err := HtmlRenderer{}.Render(buildDeterministicReport(), th)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	checkGolden(t, "custom_theme.html", out)
}

// --- bar chart JSON ---

func TestBarChartScriptLabelsAndData(t *testing.T) {
	chart := NewBarChart("Sales", map[string]float64{
		"Bananas": 30, "Apples": 50, "Cherries": 20,
	})
	cfg := parseChartScript(t, renderBarChartScript("id", chart, DefaultTheme()))

	if cfg.Type != "bar" {
		t.Errorf("Type: got %q, want bar", cfg.Type)
	}
	wantLabels := []string{"Apples", "Bananas", "Cherries"}
	if len(cfg.Data.Labels) != len(wantLabels) {
		t.Fatalf("Labels: got %v, want %v", cfg.Data.Labels, wantLabels)
	}
	for i, lbl := range wantLabels {
		if cfg.Data.Labels[i] != lbl {
			t.Errorf("Labels[%d]: got %q, want %q", i, cfg.Data.Labels[i], lbl)
		}
	}
	if len(cfg.Data.Datasets) != 1 {
		t.Fatalf("Datasets: got %d, want 1", len(cfg.Data.Datasets))
	}
	// data values must follow sorted label order: Apples=50, Bananas=30, Cherries=20
	wantData := []float64{50, 30, 20}
	for i, v := range wantData {
		if cfg.Data.Datasets[0].Data[i] != v {
			t.Errorf("Data[%d]: got %v, want %v", i, cfg.Data.Datasets[0].Data[i], v)
		}
	}
	if cfg.Options.Plugins == nil || cfg.Options.Plugins.Legend == nil || cfg.Options.Plugins.Legend.Display {
		t.Error("Legend.Display must be false for single-series bar chart")
	}
	if cfg.Options.IndexAxis != "" {
		t.Errorf("IndexAxis: got %q, want empty for vertical chart", cfg.Options.IndexAxis)
	}
}

func TestBarChartScriptHorizontal(t *testing.T) {
	chart := &BarChart{
		ChartBase:    ChartBase{BaseElement: newBaseElement(), Title: "H"},
		Data:         map[string]float64{"A": 1},
		IsHorizontal: true,
	}
	cfg := parseChartScript(t, renderBarChartScript("id", chart, DefaultTheme()))
	if cfg.Options.IndexAxis != "y" {
		t.Errorf("IndexAxis: got %q, want y for horizontal bar", cfg.Options.IndexAxis)
	}
}

// --- pie/donut chart JSON ---

func TestPieChartScriptTypeToggle(t *testing.T) {
	pie := NewPieChart("P", map[string]float64{"A": 40, "B": 60})
	cfgPie := parseChartScript(t, renderPieChartScript("id", pie, DefaultTheme()))
	if cfgPie.Type != "pie" {
		t.Errorf("pie Type: got %q, want pie", cfgPie.Type)
	}

	donut := &PieChart{
		ChartBase: ChartBase{BaseElement: newBaseElement(), Title: "D"},
		Data:      map[string]float64{"A": 40, "B": 60},
		IsDonut:   true,
	}
	cfgDonut := parseChartScript(t, renderPieChartScript("id", donut, DefaultTheme()))
	if cfgDonut.Type != "doughnut" {
		t.Errorf("donut Type: got %q, want doughnut", cfgDonut.Type)
	}
}

// --- line chart JSON ---

func TestLineChartScriptMultiSeriesLabels(t *testing.T) {
	lc := NewLineChart("Trend", []LineSeries{
		{Name: "Alpha", Points: map[string]float64{"Q1": 10, "Q2": 20}},
		{Name: "Beta", Points: map[string]float64{"Q2": 5, "Q3": 15}},
	})
	cfg := parseChartScript(t, renderLineChartScript("id", lc, DefaultTheme()))

	// legend visible for multi-series
	if cfg.Options.Plugins == nil || cfg.Options.Plugins.Legend == nil || !cfg.Options.Plugins.Legend.Display {
		t.Error("Legend.Display must be true for multi-series line chart")
	}
	// labels: Q1, Q2 from Alpha (sorted), then Q3 from Beta (new)
	wantLabels := []string{"Q1", "Q2", "Q3"}
	if len(cfg.Data.Labels) != len(wantLabels) {
		t.Fatalf("Labels: got %v, want %v", cfg.Data.Labels, wantLabels)
	}
	for i, lbl := range wantLabels {
		if cfg.Data.Labels[i] != lbl {
			t.Errorf("Labels[%d]: got %q, want %q", i, cfg.Data.Labels[i], lbl)
		}
	}
	if len(cfg.Data.Datasets) != 2 {
		t.Fatalf("Datasets: got %d, want 2", len(cfg.Data.Datasets))
	}
	if cfg.Data.Datasets[0].Label != "Alpha" || cfg.Data.Datasets[1].Label != "Beta" {
		t.Errorf("Dataset labels: got [%q %q], want [Alpha Beta]",
			cfg.Data.Datasets[0].Label, cfg.Data.Datasets[1].Label)
	}
}

func TestLineChartScriptShowPointsFalse(t *testing.T) {
	lc := NewLineChartSingle("S", map[string]float64{"Jan": 1})
	lc.ShowPoints = false
	cfg := parseChartScript(t, renderLineChartScript("id", lc, DefaultTheme()))

	if len(cfg.Data.Datasets) != 1 {
		t.Fatalf("Datasets: got %d, want 1", len(cfg.Data.Datasets))
	}
	ds := cfg.Data.Datasets[0]
	if ds.PointStyle == nil {
		t.Fatal("PointStyle must be set (not nil) when ShowPoints=false")
	}
	if *ds.PointStyle {
		t.Error("PointStyle must be false when ShowPoints=false")
	}
}

// --- stacked bar chart JSON ---

func TestStackedBarChartScriptAxesAndOrder(t *testing.T) {
	s := NewStackedBarChart("SB", []StackedBarSeries{
		{Category: "Q1", Values: map[string]float64{"North": 10, "South": 20}},
		{Category: "Q2", Values: map[string]float64{"North": 15, "South": 25}},
	})
	cfg := parseChartScript(t, renderStackedBarChartScript("id", s, DefaultTheme()))

	if cfg.Type != "bar" {
		t.Errorf("Type: got %q, want bar", cfg.Type)
	}
	// categories in original Series slice order
	wantLabels := []string{"Q1", "Q2"}
	if len(cfg.Data.Labels) != len(wantLabels) {
		t.Fatalf("Labels: got %v, want %v", cfg.Data.Labels, wantLabels)
	}
	for i, lbl := range wantLabels {
		if cfg.Data.Labels[i] != lbl {
			t.Errorf("Labels[%d]: got %q, want %q", i, cfg.Data.Labels[i], lbl)
		}
	}
	if cfg.Options.Scales == nil {
		t.Fatal("Scales must not be nil for stacked chart")
	}
	if cfg.Options.Scales.X == nil || !cfg.Options.Scales.X.Stacked {
		t.Error("X axis must be stacked")
	}
	if cfg.Options.Scales.Y == nil || !cfg.Options.Scales.Y.Stacked {
		t.Error("Y axis must be stacked")
	}
	// series names sorted alphabetically → North, South
	if len(cfg.Data.Datasets) != 2 {
		t.Fatalf("Datasets: got %d, want 2", len(cfg.Data.Datasets))
	}
	if cfg.Data.Datasets[0].Label != "North" || cfg.Data.Datasets[1].Label != "South" {
		t.Errorf("Dataset labels: got [%q %q], want [North South]",
			cfg.Data.Datasets[0].Label, cfg.Data.Datasets[1].Label)
	}
	// North=[10,15], South=[20,25]
	if cfg.Data.Datasets[0].Data[0] != 10 || cfg.Data.Datasets[0].Data[1] != 15 {
		t.Errorf("North data: got %v, want [10 15]", cfg.Data.Datasets[0].Data)
	}
	if cfg.Data.Datasets[1].Data[0] != 20 || cfg.Data.Datasets[1].Data[1] != 25 {
		t.Errorf("South data: got %v, want [20 25]", cfg.Data.Datasets[1].Data)
	}
}

func TestStackedBarChartScriptHorizontal(t *testing.T) {
	s := &StackedBarChart{
		ChartBase:    ChartBase{BaseElement: newBaseElement(), Title: "H"},
		Series:       []StackedBarSeries{{Category: "C1", Values: map[string]float64{"S1": 1}}},
		IsHorizontal: true,
	}
	cfg := parseChartScript(t, renderStackedBarChartScript("id", s, DefaultTheme()))
	if cfg.Options.IndexAxis != "y" {
		t.Errorf("IndexAxis: got %q, want y for horizontal stacked bar", cfg.Options.IndexAxis)
	}
}

// --- idGen ---

func TestIDGenAllBranches(t *testing.T) {
	g := newIDGen()
	cases := []struct {
		section, chart string
		want           string
	}{
		{"My Section", "My Chart", "my-section-my-chart"},
		{"", "Just Chart", "just-chart"},
		{"Only Section", "", "only-section"},
		{"", "", "chart"},
	}
	for _, tc := range cases {
		got := g.next(tc.section, tc.chart)
		if got != tc.want {
			t.Errorf("next(%q,%q): got %q, want %q", tc.section, tc.chart, got, tc.want)
		}
	}
}

func TestIDGenCollision(t *testing.T) {
	g := newIDGen()
	id1 := g.next("Section", "Sales Chart")
	id2 := g.next("Section", "Sales Chart")
	id3 := g.next("Section", "Sales Chart")

	if id1 != "section-sales-chart" {
		t.Errorf("first: got %q, want section-sales-chart", id1)
	}
	if id2 != "section-sales-chart-2" {
		t.Errorf("second (collision): got %q, want section-sales-chart-2", id2)
	}
	if id3 != "section-sales-chart-3" {
		t.Errorf("third (collision): got %q, want section-sales-chart-3", id3)
	}
}

// --- shadowCSS ---

func TestShadowCSS(t *testing.T) {
	cases := []struct {
		intensity string
		want      string
	}{
		{"none", "none"},
		{"subtle", "0 1px 2px rgba(0,0,0,0.05)"},
		{"medium", "0 4px 6px rgba(0,0,0,0.07)"},
		{"strong", "0 10px 15px rgba(0,0,0,0.1)"},
		{"", "0 4px 6px rgba(0,0,0,0.07)"},
		{"unknown", "0 4px 6px rgba(0,0,0,0.07)"},
	}
	for _, tc := range cases {
		got := shadowCSS(tc.intensity)
		if got != tc.want {
			t.Errorf("shadowCSS(%q): got %q, want %q", tc.intensity, got, tc.want)
		}
	}
}

// --- tooltipIcon ---

func TestTooltipIconEmpty(t *testing.T) {
	if got := tooltipIcon(""); got != "" {
		t.Errorf("empty tooltip: got %q, want empty string", got)
	}
}

func TestTooltipIconNonEmpty(t *testing.T) {
	got := tooltipIcon("click here for info")
	if !strings.Contains(got, `data-tooltip="click here for info"`) {
		t.Errorf("tooltip: missing data-tooltip attr; got %q", got)
	}
	if !strings.Contains(got, "tooltip-icon") {
		t.Errorf("tooltip: missing tooltip-icon class; got %q", got)
	}
}

func TestTooltipIconEscapesHTML(t *testing.T) {
	got := tooltipIcon(`<b>"danger"</b>`)
	if strings.Contains(got, "<b>") {
		t.Error("tooltip: unescaped < in data-tooltip attribute")
	}
	if !strings.Contains(got, "&lt;b&gt;") {
		t.Error("tooltip: expected &lt;b&gt; in escaped output")
	}
}

// --- chartColors ---

func TestChartColorsDefault(t *testing.T) {
	colors := chartColors(&Theme{})
	if len(colors) != len(defaultChartColors) {
		t.Fatalf("default colors len: got %d, want %d", len(colors), len(defaultChartColors))
	}
	for i, c := range defaultChartColors {
		if colors[i] != c {
			t.Errorf("colors[%d]: got %q, want %q", i, colors[i], c)
		}
	}
}

func TestChartColorsThemeOverride(t *testing.T) {
	custom := []string{"#aabbcc", "#ddeeff"}
	colors := chartColors(&Theme{ChartColors: custom})
	if len(colors) != 2 || colors[0] != "#aabbcc" || colors[1] != "#ddeeff" {
		t.Errorf("override colors: got %v, want %v", colors, custom)
	}
}

// --- security regression (spec 002): chart data must not break out of <script> ---

func TestChartLabelInjectionEscaped(t *testing.T) {
	r := NewReport("Sec Test")
	section := &Section{Title: "S"}
	section.AddElement(NewBarChart("</script>alert(1)//", map[string]float64{
		"</script>alert(1)": 99,
	}))
	r.AddSection(section)

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	// Find the <script> block containing chart init code.
	scriptOpen := strings.Index(out, "<script>\n")
	scriptClose := strings.LastIndex(out, "</script>")
	if scriptOpen < 0 || scriptClose <= scriptOpen {
		t.Fatal("no <script> block found")
	}
	scriptBlock := out[scriptOpen : scriptClose+len("</script>")]
	// The script block must contain exactly one </script> — its own closing tag.
	if strings.Count(scriptBlock, "</script>") != 1 {
		t.Errorf("chart <script> block contains %d </script> occurrences; injection not escaped",
			strings.Count(scriptBlock, "</script>"))
	}
}

// --- render coverage: footer, logo URL, tile subtitles, CSS branches ---

func TestRenderReportFooterAndLogo(t *testing.T) {
	r := NewReport("F")
	r.Footer = "Confidential footer"
	r.LogoURL = "https://example.com/logo.png"
	r.AddSection(&Section{})

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "Confidential footer") {
		t.Error("output must contain footer text")
	}
	if !strings.Contains(out, "https://example.com/logo.png") {
		t.Error("output must contain logo URL in img src")
	}
}

func TestNumberTileWithSubtitle(t *testing.T) {
	r := NewReport("N")
	section := &Section{}
	section.AddElement(&NumberTile{
		BaseElement: newBaseElement(),
		Title:       "Metric",
		Value:       42,
		Subtitle:    "vs last quarter",
	})
	r.AddSection(section)

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "vs last quarter") {
		t.Error("output must render number tile subtitle")
	}
}

func TestDateTileWithSubtitle(t *testing.T) {
	r := NewReport("D")
	section := &Section{}
	section.AddElement(&DateTile{
		BaseElement: newBaseElement(),
		Title:       "Date",
		Value:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Subtitle:    "fiscal year end",
	})
	r.AddSection(section)

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "fiscal year end") {
		t.Error("output must render date tile subtitle")
	}
}

func TestEnableGradientsCSSProducesGradient(t *testing.T) {
	th := DefaultTheme()
	th.EnableGradients = true
	r := NewReport("G")
	r.AddSection(&Section{})

	out, err := HtmlRenderer{}.Render(r, th)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "linear-gradient") {
		t.Error("output must contain linear-gradient in CSS when EnableGradients=true")
	}
}

func TestEnableAnimationsOffOmitsFadeIn(t *testing.T) {
	th := DefaultTheme()
	th.EnableAnimations = false
	r := NewReport("A")
	r.AddSection(&Section{})

	out, err := HtmlRenderer{}.Render(r, th)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if strings.Contains(out, "@keyframes fadeIn") {
		t.Error("output must not include fadeIn keyframe when EnableAnimations=false")
	}
}

// --- unknown element type ---

func TestRenderElementUnknownType(t *testing.T) {
	r := NewReport("U")
	section := &Section{}
	section.AddElement(&customTestElement{BaseElement: newBaseElement()})
	r.AddSection(section)

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "unknown element") {
		t.Error("output must contain 'unknown element' comment for unrecognized element types")
	}
}

// --- table: non-string any cell values ---

func TestTableNonStringAnyCell(t *testing.T) {
	r := NewReport("T")
	section := &Section{}
	section.AddElement(NewTableWithColumns("Test", []map[string]any{
		{"count": 42, "ratio": 3.14},
	}, []string{"count", "ratio"}))
	r.AddSection(section)

	out, err := HtmlRenderer{}.Render(r, nil)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, ">42<") {
		t.Error("output must contain int cell value 42")
	}
	if !strings.Contains(out, ">3.14<") {
		t.Error("output must contain float cell value 3.14")
	}
}
