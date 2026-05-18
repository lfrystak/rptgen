package rptgen

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
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

// Render converts report to a complete HTML document using the provided theme.
// If theme is nil, DefaultTheme() is used.
func (h HtmlRenderer) Render(report *Report, theme *Theme) (string, error) {
	if theme == nil {
		theme = DefaultTheme()
	}

	gen := newIDGen()
	var chartScripts []template.JS
	sections := make([]template.HTML, 0, len(report.Sections))

	for _, section := range report.Sections {
		sectionHTML, scripts := renderSectionHTML(section, theme, gen)
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

	var buf bytes.Buffer
	if err := docTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// renderSectionHTML pre-renders a section to an HTML string and collects chart init scripts.
func renderSectionHTML(section *Section, theme *Theme, gen *idGen) (string, []string) {
	var b strings.Builder
	var scripts []string

	colTemplate := columnWidthsToCSS(section.ColumnWidths)
	b.WriteString("    <section class=\"report-section\">\n")
	if section.Title != "" {
		fmt.Fprintf(&b, "      <h2 class=\"section-title\">%s</h2>\n", html.EscapeString(section.Title))
	}
	fmt.Fprintf(&b, "      <div class=\"section-grid\" style=\"--col-template: %s\">\n", colTemplate)
	for _, elem := range section.Elements {
		b.WriteString("        <div class=\"element-wrapper\">\n")
		rendered, elemScripts := renderElement(elem, theme, gen, section.Title)
		b.WriteString(rendered)
		scripts = append(scripts, elemScripts...)
		b.WriteString("        </div>\n")
	}
	b.WriteString("      </div>\n")
	b.WriteString("    </section>\n")
	return b.String(), scripts
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
