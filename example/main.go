package main

import (
	"log"
	"os"
	"time"

	"github.com/lfrystak/rptgen"
)

const (
	dateFormat = "2. 1. 2006"
)

func main() {
	// Sample data — simulating data extracted from an API
	salesByRegion := []rptgen.DataPoint{
		{Label: "Asia Pacific", Value: 142000},
		{Label: "North America", Value: 125000},
		{Label: "Europe", Value: 98000},
		{Label: "Latin America", Value: 67000},
	}

	monthlyRevenue := []rptgen.DataPoint{
		{Label: "January", Value: 45000},
		{Label: "February", Value: 52000},
		{Label: "March", Value: 48000},
		{Label: "April", Value: 55000},
		{Label: "May", Value: 61000},
		{Label: "June", Value: 58000},
	}

	productMix := []rptgen.DataPoint{
		{Label: "Enterprise", Value: 45},
		{Label: "Professional", Value: 30},
		{Label: "Starter", Value: 25},
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
	report.LogoURL = "https://raw.githubusercontent.com/lfrystak/rptgen/refs/heads/main/img/rptgen_logo_compact.svg"
	report.Footer = "Confidential - Internal Use Only"

	// Section: Key Metrics
	kpis := rptgen.NewSection("Key Metrics", rptgen.EqualColumns(4)...)

	revenue := rptgen.NewNumberTile("Total Revenue", 432000)
	revenue.Format = "%.0f"
	revenue.Prefix = "$ "
	revenue.ThousandsSep = true
	revenue.Tooltip = "Total revenue from all regions and product lines for Q2 2024"
	kpis.AddElement(revenue)

	customers := rptgen.NewNumberTile("New Customers", 47)
	customers.Format = "%.0f"
	customers.Tooltip = "Number of new customer accounts opened during this quarter"
	kpis.AddElement(customers)

	growth := rptgen.NewNumberTile("Growth Rate", 19.8)
	growth.Format = "%.1f%%"
	growth.Subtitle = "↑ vs Q1"
	growth.Tooltip = "Year-over-year growth compared to the same quarter last year"
	kpis.AddElement(growth)

	satisfaction := rptgen.NewNumberTile("Customer Satisfaction", 4.7)
	satisfaction.Format = "%.1f"
	satisfaction.Subtitle = "out of 5.0"
	kpis.AddElement(satisfaction)

	report.AddSection(kpis)

	// Section: Timeline
	timeline := rptgen.NewSection("Timeline", 1, 1, 1)

	qStart := rptgen.NewDateTile("Quarter Start", time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC))
	qStart.Format = dateFormat
	timeline.AddElement(qStart)

	qEnd := rptgen.NewDateTile("Quarter End", time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC))
	qEnd.Format = dateFormat
	timeline.AddElement(qEnd)

	generated := rptgen.NewDateTile("Report Generated", time.Date(2024, 7, 15, 14, 30, 0, 0, time.UTC))
	generated.Format = "2006-01-02 15:04"
	generated.Subtitle = "Generated date"
	timeline.AddElement(generated)

	report.AddSection(timeline)

	// Section: Revenue Analysis
	revenueSection := rptgen.NewSection("Revenue Analysis", 1, 1)
	revenueSection.AddElement(func() rptgen.Element {
		c := rptgen.NewBarChart("Revenue by Region", salesByRegion)
		c.Tooltip = "Revenue distribution across our four main geographic markets for Q2 2024"
		return c
	}())
	revenueSection.AddElement(func() rptgen.Element {
		c := rptgen.NewLineChartSingle("Monthly Revenue Trend", monthlyRevenue)
		c.Tooltip = "Month-over-month revenue progression showing steady growth throughout the quarter"
		return c
	}())
	report.AddSection(revenueSection)

	// Section: Product & Growth
	productGrowth := rptgen.NewSection("Product & Growth", 1, 1)
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
	topCustomersSection := rptgen.NewSection("Top Customers")
	topCustomersSection.AddElement(rptgen.NewTableWithColumns(
		"Q2 Top Revenue Contributors",
		topCustomers,
		[]string{"Customer", "Revenue", "Status"},
	))
	report.AddSection(topCustomersSection)

	// Section: Canvas Example - Mixed Layout (2:1 ratio)
	canvasSection := rptgen.NewSection("Canvas Example - Mixed Layout", 3, 2)
	canvas := rptgen.NewCanvas(1, 3)

	canvasTile := rptgen.NewNumberTile("Canvas Demo", 100)
	canvasTile.Format = "%.0f"
	canvasTile.Subtitle = "This is in a canvas!"
	canvas.AddElement(canvasTile)

	canvas.AddElement(rptgen.NewBarChart("Sample Chart", []rptgen.DataPoint{{Label: "A", Value: 10}, {Label: "B", Value: 20}, {Label: "C", Value: 15}}))

	canvas.AddElement(rptgen.NewNumberTile("Another Tile", 42))

	today := rptgen.NewDateTile("Today", time.Now())
	today.Format = dateFormat
	canvas.AddElement(today)

	canvasSection.AddElement(canvas)
	canvasSection.AddElement(func() rptgen.Element {
		c := rptgen.NewPieChart("Sample Pie", []rptgen.DataPoint{{Label: "X", Value: 30}, {Label: "Y", Value: 40}, {Label: "Z", Value: 30}})
		c.IsDonut = true
		return c
	}())
	report.AddSection(canvasSection)

	// Section: Summary
	summary := rptgen.NewSection("Summary")
	summary.AddElement(rptgen.NewFreeText(`Q2 2024 showed strong performance across all regions with total revenue of $432,000, representing a 19.8% increase over Q1. Asia Pacific continues to be our strongest market, while Latin America presents significant growth opportunities.

The Enterprise product tier dominates our revenue mix at 45%, indicating strong traction in the high-value segment. Customer satisfaction remains high at 4.7/5.0.

Key focus areas for Q3:
• Expand sales team in Latin America
• Launch new features for Professional tier
• Increase customer retention initiatives`))
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
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	err = rptgen.HtmlRenderer{}.Render(f, report, theme)
	if err != nil {
		log.Fatal(err)
	}
	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", path)
}
