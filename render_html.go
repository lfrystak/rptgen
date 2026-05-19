package rptgen

import (
	"bufio"
	"fmt"
	"html"
	"html/template"
	"io"
	"strings"
)

// HtmlRenderer renders a Report to a self-contained HTML document string.
type HtmlRenderer struct{}

// docData is the context passed to the HTML document template.
type docData struct {
	Title        string
	CSS          template.CSS
	LogoURL      string
	GeneratedAt  string
	Footer       string
	Sections     []template.HTML
	ChartJS      template.JS
	ChartScripts []template.JS
}

// docTemplate is the outer HTML document structure. User-supplied strings (Title,
// LogoURL, GeneratedAt, Footer) are auto-escaped by html/template. Section content
// and chart scripts are pre-rendered and marked safe via template.HTML / template.JS.
var docTemplate = template.Must(template.New("doc").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}}</title>
  <style>
{{.CSS}}</style>
</head>
<body>
  <div class="report">
    <header class="report-header">
{{- if .LogoURL}}
      <img src="{{.LogoURL}}" alt="logo">
{{- end}}
      <h1 class="report-title">{{.Title}}</h1>
{{- if .GeneratedAt}}
      <span class="report-generated">Generated: {{.GeneratedAt}}</span>
{{- end}}
    </header>

{{range .Sections}}{{.}}
{{end}}
{{- if .Footer}}
    <footer class="report-footer">{{.Footer}}</footer>
{{end -}}
  </div>

{{- if .ChartScripts}}
  <script>{{.ChartJS}}</script>
  <script>
{{- range .ChartScripts}}
  {{.}}
{{- end}}
  </script>
{{- end}}
</body>
</html>`))

// Render converts report to a complete HTML document and writes it to w.
// If theme is nil, DefaultTheme() is used.
func (h HtmlRenderer) Render(w io.Writer, report *Report, theme *Theme) error {
	if theme == nil {
		theme = DefaultTheme()
	}

	gen := newIDGen()
	var chartScripts []template.JS
	sections := make([]template.HTML, 0, len(report.Sections))

	for _, section := range report.Sections {
		sectionHTML, scripts, err := renderSectionHTML(section, theme, gen)
		if err != nil {
			return err
		}
		sections = append(sections, template.HTML(sectionHTML))
		for _, s := range scripts {
			chartScripts = append(chartScripts, template.JS(s))
		}
	}

	generatedAt := ""
	if !report.GeneratedAt.IsZero() {
		generatedAt = report.GeneratedAt.Format("2006-01-02 15:04:05")
	}

	data := docData{
		Title:        report.Title,
		CSS:          template.CSS(generateCSS(theme)),
		LogoURL:      report.LogoURL,
		GeneratedAt:  generatedAt,
		Footer:       report.Footer,
		Sections:     sections,
		ChartJS:      template.JS(chartJSSource),
		ChartScripts: chartScripts,
	}

	bw := bufio.NewWriter(w)
	if err := docTemplate.Execute(bw, data); err != nil {
		return err
	}
	return bw.Flush()
}

// RenderString renders report to a string. It is a convenience wrapper around Render
// for callers that need the HTML as a string rather than streaming to a writer.
func (h HtmlRenderer) RenderString(report *Report, theme *Theme) (string, error) {
	var buf strings.Builder
	if err := h.Render(&buf, report, theme); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// renderSectionHTML pre-renders a section to an HTML string and collects chart init scripts.
func renderSectionHTML(section *Section, theme *Theme, gen *idGen) (string, []string, error) {
	var b strings.Builder
	var scripts []string

	ctx := &HTMLRenderContext{
		Theme:        theme,
		SectionTitle: section.Title,
		idGen:        gen,
	}

	colTemplate := columnWidthsToCSS(section.ColumnWidths)
	b.WriteString("    <section class=\"report-section\">\n")
	if section.Title != "" {
		fmt.Fprintf(&b, "      <h2 class=\"section-title\">%s</h2>\n", html.EscapeString(section.Title))
	}
	fmt.Fprintf(&b, "      <div class=\"section-grid\" style=\"--col-template: %s\">\n", colTemplate)
	for _, elem := range section.Elements {
		b.WriteString("        <div class=\"element-wrapper\">\n")
		rendered, elemScripts, err := renderElement(elem, ctx)
		if err != nil {
			return "", nil, err
		}
		b.WriteString(rendered)
		scripts = append(scripts, elemScripts...)
		b.WriteString("        </div>\n")
	}
	b.WriteString("      </div>\n")
	b.WriteString("    </section>\n")
	return b.String(), scripts, nil
}

// columnWidthsToCSS converts a slice of proportional widths to a CSS grid-template-columns value.
func columnWidthsToCSS(widths []int) string {
	if len(widths) == 0 {
		return "1fr"
	}
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = fmt.Sprintf("%dfr", w)
	}
	return strings.Join(parts, " ")
}
