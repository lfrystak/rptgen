package rptgen

import (
	"strings"
	"testing"
	"time"
)

func buildFullReport() *Report {
	r := NewReport("Smoke Test Report")

	section := &Section{Title: "All Elements"}

	section.AddElement(&NumberTile{
		BaseElement: newBaseElement(),
		Title:       "Revenue",
		Value:       99999.99,
		Format:      "C2",
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
	if strings.Contains(out, "cdn.jsdelivr.net") {
		t.Error("output must not reference cdn.jsdelivr.net")
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
