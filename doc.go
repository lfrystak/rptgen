// Package rptgen generates self-contained, themed HTML reports from structured Go data.
//
// # Overview
//
// rptgen lets you build a [Report] from typed [Element] values — numeric tiles, tables,
// date tiles, free text, and Chart.js charts — and render it to any [io.Writer] as a
// single, dependency-free HTML file (Chart.js is embedded at render time).
//
// # Quick start
//
//	report := rptgen.NewReport("Q2 Sales")
//
//	kpis := rptgen.NewSection("Key Metrics", rptgen.EqualColumns(2)...)
//	kpis.AddElement(rptgen.NewNumberTile("Revenue", 432_000))
//	kpis.AddElement(rptgen.NewNumberTile("New Customers", 47))
//	report.AddSection(kpis)
//
//	charts := rptgen.NewSection("Trend", 1, 2)
//	charts.AddElement(rptgen.NewPieChart("Product Mix", []rptgen.DataPoint{
//	    {Label: "Enterprise", Value: 45},
//	    {Label: "Professional", Value: 30},
//	    {Label: "Starter", Value: 25},
//	}))
//	// Categorical X axis (string labels, evenly spaced):
//	charts.AddElement(rptgen.NewLineChartSingle("Monthly Revenue", []rptgen.DataPoint{
//	    {Label: "Jan", Value: 45_000},
//	    {Label: "Feb", Value: 52_000},
//	    {Label: "Mar", Value: 48_000},
//	}))
//	// Numeric X axis (linear, proportional spacing) — for functions or continuous data:
//	charts.AddElement(rptgen.NewLineChartXY("Sin Wave", []rptgen.XYPoint{
//	    {X: 0.0, Y: 0.00}, {X: 0.5, Y: 0.48}, {X: 1.0, Y: 0.84},
//	}))
//	report.AddSection(charts)
//
//	f, _ := os.Create("report.html")
//	defer f.Close()
//	rptgen.HtmlRenderer{}.Render(f, report, nil) // nil theme uses DefaultTheme
//
// # Theming
//
// Pass a *[Theme] as the third argument to [HtmlRenderer.Render] to override colors,
// fonts, border radius, shadow intensity, and the eight-color chart palette. Start from
// [DefaultTheme] and override only the fields you need:
//
//	theme := rptgen.DefaultTheme()
//	theme.PrimaryColor = "#059669"
//	theme.FontFamily = "Georgia, serif"
//	rptgen.HtmlRenderer{}.Render(w, report, theme)
package rptgen
