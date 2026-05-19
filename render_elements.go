package rptgen

import (
	"fmt"
	"html"
	"strings"
)

// renderElement dispatches to the element's own RenderHTML implementation.
// Any Element that does not implement HTMLRenderer returns an error.
func renderElement(elem Element, ctx *HTMLRenderContext) (string, []string, error) {
	hr, ok := elem.(HTMLRenderer)
	if !ok {
		return "", nil, fmt.Errorf("rptgen: unknown element type %q", elem.ElementType())
	}
	return hr.RenderHTML(ctx)
}

// RenderHTML implements HTMLRenderer for NumberTile.
func (e *NumberTile) RenderHTML(_ *HTMLRenderContext) (string, []string, error) {
	var b strings.Builder
	b.WriteString("          <div class=\"element tile number-tile\">\n")
	b.WriteString(tooltipIcon(e.Tooltip))
	fmt.Fprintf(&b, "            <div class=\"tile-title\">%s</div>\n", html.EscapeString(e.Title))
	fmt.Fprintf(&b, "            <div class=\"tile-value\">%s</div>\n", html.EscapeString(e.FormatValue()))
	if e.Subtitle != "" {
		fmt.Fprintf(&b, "            <div class=\"tile-subtitle\">%s</div>\n", html.EscapeString(e.Subtitle))
	}
	b.WriteString("          </div>\n")
	return b.String(), nil, nil
}

// RenderHTML implements HTMLRenderer for DateTile.
func (e *DateTile) RenderHTML(_ *HTMLRenderContext) (string, []string, error) {
	var b strings.Builder
	b.WriteString("          <div class=\"element tile date-tile\">\n")
	b.WriteString(tooltipIcon(e.Tooltip))
	fmt.Fprintf(&b, "            <div class=\"tile-title\">%s</div>\n", html.EscapeString(e.Title))
	fmt.Fprintf(&b, "            <div class=\"tile-value\">%s</div>\n", html.EscapeString(e.FormatValue()))
	if e.Subtitle != "" {
		fmt.Fprintf(&b, "            <div class=\"tile-subtitle\">%s</div>\n", html.EscapeString(e.Subtitle))
	}
	b.WriteString("          </div>\n")
	return b.String(), nil, nil
}

// RenderHTML implements HTMLRenderer for FreeText.
func (e *FreeText) RenderHTML(_ *HTMLRenderContext) (string, []string, error) {
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
	return b.String(), nil, nil
}

// RenderHTML implements HTMLRenderer for Table.
func (e *Table) RenderHTML(_ *HTMLRenderContext) (string, []string, error) {
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
	return b.String(), nil, nil
}

// RenderHTML implements HTMLRenderer for Canvas.
func (e *Canvas) RenderHTML(ctx *HTMLRenderContext) (string, []string, error) {
	colTemplate := columnWidthsToCSS(e.ColumnWidths)
	var b strings.Builder
	var allScripts []string
	fmt.Fprintf(&b, "          <div class=\"element canvas-grid\" style=\"--col-template: %s\">\n", colTemplate)
	for _, child := range e.Elements {
		b.WriteString("            <div class=\"element-wrapper\">\n")
		rendered, scripts, err := renderElement(child, ctx)
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

// RenderChartContainer returns the HTML card wrapper for a Chart.js canvas element.
// Custom chart elements call this to produce the standard chart card HTML, then supply
// the chart init script separately.
func RenderChartContainer(id, title, tooltip string) string {
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
