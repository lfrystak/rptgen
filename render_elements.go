package rptgen

import (
	"fmt"
	"html"
	"strings"
)

const closingDiv = "          </div>\n"

// renderElement dispatches to the element's own renderHTML implementation.
// Any Element that does not implement htmlRenderer returns an error.
func renderElement(elem Element, ctx *htmlRenderContext) (string, []string, error) {
	hr, ok := elem.(htmlRenderer)
	if !ok {
		return "", nil, fmt.Errorf("rptgen: unknown element type %q", elem.ElementType())
	}
	return hr.renderHTML(ctx)
}

func (e *NumberTile) renderHTML(_ *htmlRenderContext) (string, []string, error) {
	return renderTileHTML("number-tile", e.Tooltip, e.Title, e.FormatValue(), e.Subtitle), nil, nil
}

func (e *DateTile) renderHTML(_ *htmlRenderContext) (string, []string, error) {
	return renderTileHTML("date-tile", e.Tooltip, e.Title, e.FormatValue(), e.Subtitle), nil, nil
}

func renderTileHTML(cssClass, tooltip, title, value, subtitle string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "          <div class=\"element tile %s\">\n", cssClass)
	b.WriteString(tooltipIcon(tooltip))
	fmt.Fprintf(&b, "            <div class=\"tile-title\">%s</div>\n", html.EscapeString(title))
	fmt.Fprintf(&b, "            <div class=\"tile-value\">%s</div>\n", html.EscapeString(value))
	if subtitle != "" {
		fmt.Fprintf(&b, "            <div class=\"tile-subtitle\">%s</div>\n", html.EscapeString(subtitle))
	}
	b.WriteString(closingDiv)
	return b.String()
}

func (e *FreeText) renderHTML(_ *htmlRenderContext) (string, []string, error) {
	var b strings.Builder
	b.WriteString("          <div class=\"element free-text\">\n")
	if e.IsHTML {
		b.WriteString("            ")
		b.WriteString(e.Content)
		b.WriteByte('\n')
	} else {
		fmt.Fprintf(&b, "            <p>%s</p>\n", html.EscapeString(e.Content))
	}
	b.WriteString(closingDiv)
	return b.String(), nil, nil
}

func (e *Table) renderHTML(_ *htmlRenderContext) (string, []string, error) {
	var b strings.Builder
	b.WriteString("          <div class=\"element table-wrapper\">\n")
	b.WriteString("            <table>\n")
	if e.Title != "" {
		fmt.Fprintf(&b, "              <caption class=\"element-title\">%s</caption>\n", html.EscapeString(e.Title))
	}
	b.WriteString("              <thead><tr>")
	for _, col := range e.Columns {
		fmt.Fprintf(&b, "<th scope=\"col\">%s</th>", html.EscapeString(col))
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

func (e *Canvas) renderHTML(ctx *htmlRenderContext) (string, []string, error) {
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
	b.WriteString(closingDiv)
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
	b.WriteString(closingDiv)
	return b.String()
}

func tooltipIcon(tooltip string) string {
	if tooltip == "" {
		return ""
	}
	esc := html.EscapeString(tooltip)
	return fmt.Sprintf("            <span class=\"tooltip-icon\" data-tooltip=\"%s\" title=\"%s\" aria-label=\"%s\" tabindex=\"0\">&#9432;</span>\n", esc, esc, esc)
}
