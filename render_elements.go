package rptgen

import (
	"fmt"
	"html"
	"strings"
)

// renderElement dispatches to the correct element renderer and returns HTML + any chart init scripts.
func renderElement(elem Element, theme *Theme, gen *idGen, sectionTitle string) (string, []string, error) {
	switch e := elem.(type) {
	case *NumberTile:
		return renderNumberTile(e), nil, nil
	case *DateTile:
		return renderDateTile(e), nil, nil
	case *FreeText:
		return renderFreeText(e), nil, nil
	case *Table:
		return renderTable(e), nil, nil
	case *Canvas:
		return renderCanvas(e, theme, gen, sectionTitle)
	case *BarChart:
		id := gen.next(sectionTitle, e.Title)
		script, err := renderBarChartScript(id, e, theme)
		if err != nil {
			return "", nil, err
		}
		return renderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
	case *LineChart:
		id := gen.next(sectionTitle, e.Title)
		script, err := renderLineChartScript(id, e, theme)
		if err != nil {
			return "", nil, err
		}
		return renderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
	case *PieChart:
		id := gen.next(sectionTitle, e.Title)
		script, err := renderPieChartScript(id, e, theme)
		if err != nil {
			return "", nil, err
		}
		return renderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
	case *StackedBarChart:
		id := gen.next(sectionTitle, e.Title)
		script, err := renderStackedBarChartScript(id, e, theme)
		if err != nil {
			return "", nil, err
		}
		return renderChartContainer(id, e.Title, e.Tooltip), []string{script}, nil
	default:
		return "", nil, fmt.Errorf("rptgen: unknown element type %q", elem.ElementType())
	}
}

func renderNumberTile(e *NumberTile) string {
	var b strings.Builder
	b.WriteString("          <div class=\"element tile number-tile\">\n")
	b.WriteString(tooltipIcon(e.Tooltip))
	fmt.Fprintf(&b, "            <div class=\"tile-title\">%s</div>\n", html.EscapeString(e.Title))
	fmt.Fprintf(&b, "            <div class=\"tile-value\">%s</div>\n", html.EscapeString(e.FormatValue()))
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

func renderCanvas(e *Canvas, theme *Theme, gen *idGen, sectionTitle string) (string, []string, error) {
	colTemplate := columnWidthsToCSS(e.ColumnWidths)
	var b strings.Builder
	var allScripts []string
	fmt.Fprintf(&b, "          <div class=\"element canvas-grid\" style=\"--col-template: %s\">\n", colTemplate)
	for _, child := range e.Elements {
		b.WriteString("            <div class=\"element-wrapper\">\n")
		rendered, scripts, err := renderElement(child, theme, gen, sectionTitle)
		if err != nil {
			return "", nil, err
		}
		b.WriteString(rendered)
		allScripts = append(allScripts, scripts...)
		b.WriteString("            </div>\n")
	}
	b.WriteString("          </div>\n")
	return b.String(), allScripts, nil
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
