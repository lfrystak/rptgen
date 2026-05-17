package main

import (
	"log"
	"os"
	"time"

	"github.com/lfrystak/rptgen"
)

func main() {
	// Sample data — simulating data extracted from an API
	salesByRegion := map[string]float64{
		"North America": 125000,
		"Europe":        98000,
		"Asia Pacific":  142000,
		"Latin America": 67000,
	}

	monthlyRevenue := map[string]float64{
		"January":  45000,
		"February": 52000,
		"March":    48000,
		"April":    55000,
		"May":      61000,
		"June":     58000,
	}

	productMix := map[string]float64{
		"Enterprise":   45,
		"Professional": 30,
		"Starter":      25,
	}

	quarterlyGrowth := []rptgen.StackedBarSeries{
		{Category: "Q1", Values: map[string]float64{"Revenue": 145000, "Costs": 95000}},
		{Category: "Q2", Values: map[string]float64{"Revenue": 174000, "Costs": 102000}},
	}

	topCustomers := []map[string]any{
		{"Customer": "Acme Corp", "Revenue": 45000, "Status": "Active"},
		{"Customer": "TechStart Inc", "Revenue": 38000, "Status": "Active"},
		{"Customer": "Global Solutions", "Revenue": 35000, "Status": "Active"},
		{"Customer": "Innovation Labs", "Revenue": 28000, "Status": "Pending"},
	}

	// Build the report
	report := rptgen.NewReport("Q2 2024 Business Performance Report")
	report.LogoURL = "https://via.placeholder.com/150x50/2563eb/ffffff?text=MyCompany"
	report.Locale = "en-US"
	report.Footer = "Confidential - Internal Use Only"

	// Section: Key Metrics
	kpis := &rptgen.Section{Title: "Key Metrics", ColumnWidths: []int{1, 1, 1, 1}}
	kpis.AddElement(&rptgen.NumberTile{
		Title:   "Total Revenue",
		Value:   432000,
		Format:  "C0",
		Tooltip: "Total revenue from all regions and product lines for Q2 2024",
	})
	kpis.AddElement(&rptgen.NumberTile{
		Title:   "New Customers",
		Value:   47,
		Format:  "N0",
		Tooltip: "Number of new customer accounts opened during this quarter",
	})
	kpis.AddElement(&rptgen.NumberTile{
		Title:    "Growth Rate",
		Value:    0.198,
		Format:   "P1",
		Subtitle: "↑ vs Q1",
		Tooltip:  "Year-over-year growth compared to the same quarter last year",
	})
	kpis.AddElement(&rptgen.NumberTile{
		Title:    "Customer Satisfaction",
		Value:    4.7,
		Format:   "N1",
		Subtitle: "out of 5.0",
	})
	report.AddSection(kpis)

	// Section: Timeline
	timeline := &rptgen.Section{Title: "Timeline", ColumnWidths: []int{1, 1, 1}}
	timeline.AddElement(&rptgen.DateTile{
		Title:  "Quarter Start",
		Value:  time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		Format: "January 02, 2006",
	})
	timeline.AddElement(&rptgen.DateTile{
		Title:  "Quarter End",
		Value:  time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
		Format: "January 02, 2006",
	})
	timeline.AddElement(&rptgen.DateTile{
		Title:    "Report Generated",
		Value:    time.Date(2024, 7, 15, 14, 30, 0, 0, time.UTC),
		Format:   "2006-01-02 15:04",
		Subtitle: "Generated date",
	})
	report.AddSection(timeline)

	// Section: Revenue Analysis
	revenue := &rptgen.Section{Title: "Revenue Analysis", ColumnWidths: []int{1, 1}}
	revenue.AddElement(func() rptgen.Element {
		c := rptgen.NewBarChart("Revenue by Region", salesByRegion)
		c.Tooltip = "Revenue distribution across our four main geographic markets for Q2 2024"
		return c
	}())
	revenue.AddElement(func() rptgen.Element {
		c := rptgen.NewLineChartSingle("Monthly Revenue Trend", monthlyRevenue)
		c.Tooltip = "Month-over-month revenue progression showing steady growth throughout the quarter"
		return c
	}())
	report.AddSection(revenue)

	// Section: Product & Growth
	productGrowth := &rptgen.Section{Title: "Product & Growth", ColumnWidths: []int{1, 1}}
	productGrowth.AddElement(func() rptgen.Element {
		c := rptgen.NewPieChart("Product Mix (%)", productMix)
		c.IsDonut = true
		c.Tooltip = "Percentage breakdown of revenue by product tier - Enterprise, Professional, and Starter packages"
		return c
	}())
	productGrowth.AddElement(func() rptgen.Element {
		c := rptgen.NewStackedBarChart("Quarterly Performance", quarterlyGrowth)
		c.Tooltip = "Comparison of revenue vs costs for Q1 and Q2, showing improved profit margins"
		return c
	}())
	report.AddSection(productGrowth)

	// Section: Top Customers
	customers := &rptgen.Section{Title: "Top Customers"}
	customers.AddElement(rptgen.NewTableWithColumns(
		"Q2 Top Revenue Contributors",
		topCustomers,
		[]string{"Customer", "Revenue", "Status"},
	))
	report.AddSection(customers)

	// Section: Canvas Example - Mixed Layout (1:2 ratio)
	canvasSection := &rptgen.Section{Title: "Canvas Example - Mixed Layout", ColumnWidths: []int{1, 2}}
	canvas := rptgen.NewCanvas(1, 1)
	canvas.AddElement(&rptgen.NumberTile{
		Title:    "Canvas Demo",
		Value:    100,
		Format:   "N0",
		Subtitle: "This is in a canvas!",
	})
	canvas.AddElement(rptgen.NewBarChart("Sample Chart", map[string]float64{"A": 10, "B": 20, "C": 15}))
	canvas.AddElement(&rptgen.NumberTile{Title: "Another Tile", Value: 42, Format: "N0"})
	canvas.AddElement(&rptgen.DateTile{
		Title:  "Today",
		Value:  time.Now(),
		Format: "January 02, 2006",
	})
	canvasSection.AddElement(canvas)
	canvasSection.AddElement(func() rptgen.Element {
		c := rptgen.NewPieChart("Sample Pie", map[string]float64{"X": 30, "Y": 40, "Z": 30})
		c.IsDonut = true
		return c
	}())
	report.AddSection(canvasSection)

	// Section: Summary
	summary := &rptgen.Section{Title: "Summary"}
	summary.AddElement(&rptgen.FreeText{
		Content: `Q2 2024 showed strong performance across all regions with total revenue of $432,000, representing a 19.8% increase over Q1. Asia Pacific continues to be our strongest market, while Latin America presents significant growth opportunities.

The Enterprise product tier dominates our revenue mix at 45%, indicating strong traction in the high-value segment. Customer satisfaction remains high at 4.7/5.0.

Key focus areas for Q3:
• Expand sales team in Latin America
• Launch new features for Professional tier
• Increase customer retention initiatives`,
	})
	report.AddSection(summary)

	// Render default theme
	write("report.html", report, nil)

	// Render custom theme (emerald green + indigo)
	customTheme := rptgen.DefaultTheme()
	customTheme.PrimaryColor = "#059669"
	customTheme.SecondaryColor = "#6366f1"
	customTheme.BackgroundColor = "#ffffff"
	customTheme.TextColor = "#111827"
	customTheme.FontFamily = "Georgia, serif"
	write("report-custom.html", report, customTheme)
}

func write(path string, report *rptgen.Report, theme *rptgen.Theme) {
	html, err := rptgen.HtmlRenderer{}.Render(report, theme)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(html), 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", path)
}
