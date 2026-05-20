package rptgen

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed assets/report.css.tmpl
var cssTmplSrc string

var cssTemplate = template.Must(template.New("css").Parse(cssTmplSrc))

// cssData extends Theme with computed values for the CSS template.
type cssData struct {
	*Theme
	Shadow string
}

// generateCSS renders the CSS template with theme values applied.
func generateCSS(t *Theme) string {
	data := &cssData{Theme: t, Shadow: shadowCSS(t.ShadowIntensity)}
	var b strings.Builder
	if err := cssTemplate.Execute(&b, data); err != nil {
		return ""
	}
	return b.String()
}

func shadowCSS(intensity string) string {
	switch intensity {
	case "none":
		return "none"
	case "subtle":
		return "0 1px 2px rgba(0,0,0,0.05)"
	case "strong":
		return "0 10px 15px rgba(0,0,0,0.1)"
	default: // "medium" or empty
		return "0 4px 6px rgba(0,0,0,0.07)"
	}
}
