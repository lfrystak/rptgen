package rptgen_test

import (
	"fmt"
	"io"

	"github.com/lfrystak/rptgen"
)

// ExampleHtmlRenderer_Render shows the minimal steps to build and render a report.
func ExampleHtmlRenderer_Render() {
	report := rptgen.NewReport("Sales Summary")

	section := rptgen.NewSection("Key Metrics", rptgen.EqualColumns(2)...)
	section.AddElement(rptgen.NewNumberTile("Revenue", 432_000))
	section.AddElement(rptgen.NewNumberTile("New Customers", 47))
	report.AddSection(section)

	section2 := rptgen.NewSection("Trend")
	section2.AddElement(rptgen.NewLineChartSingle("Monthly Revenue", []rptgen.DataPoint{
		{Label: "Jan", Value: 45_000},
		{Label: "Feb", Value: 52_000},
		{Label: "Mar", Value: 48_000},
	}))
	report.AddSection(section2)

	if err := (rptgen.HtmlRenderer{}).Render(io.Discard, report, nil); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output: ok
}
