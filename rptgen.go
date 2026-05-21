package rptgen

import (
	"io"
	"time"
)

// Element is the interface all report elements implement.
type Element interface {
	ElementType() string
}

// Report is the top-level document container.
type Report struct {
	Title       string
	Lang        string // BCP-47 language tag for the HTML lang attribute, e.g. "en", "de". Defaults to "en".
	Sections    []*Section
	GeneratedAt time.Time
	Footer      string
	LogoURL     string
	LogoAlt     string // alt text for the logo image; defaults to "<Title> logo" when empty.
}

// NewReport returns a Report with GeneratedAt set to the current time.
func NewReport(title string) *Report {
	return &Report{
		Title:       title,
		GeneratedAt: time.Now(),
		Sections:    []*Section{},
	}
}

// AddSection appends s to the report and returns r for optional chaining.
func (r *Report) AddSection(s *Section) *Report {
	r.Sections = append(r.Sections, s)
	return r
}

// Section groups elements into a layout row within a Report.
// ColumnWidths holds proportional widths (e.g. [1,2] = 33%/67%); nil or empty = single column.
//
// Section is a top-level row managed by Report and is not itself an Element.
// Use Canvas when you need a nestable sub-grid element that sits inside a Section column.
type Section struct {
	Title        string
	Elements     []Element
	ColumnWidths []int
}

// NewSection returns a Section with the given title and optional proportional column widths.
// NewSection("Stats", 1, 2) creates a section with a 33%/67% two-column grid.
// Omit columnWidths (or pass none) for a single-column layout.
func NewSection(title string, columnWidths ...int) *Section {
	return &Section{
		Title:        title,
		ColumnWidths: columnWidths,
	}
}

// AddElement appends e to the section and returns s for optional chaining.
func (s *Section) AddElement(e Element) *Section {
	s.Elements = append(s.Elements, e)
	return s
}

// EqualColumns returns a ColumnWidths slice of n equal proportional units.
// EqualColumns(4) is shorthand for []int{1, 1, 1, 1}.
func EqualColumns(n int) []int {
	if n <= 0 {
		return nil
	}
	cols := make([]int, n)
	for i := range cols {
		cols[i] = 1
	}
	return cols
}

// Theme controls the visual appearance of a rendered report.
// Empty string fields mean "use default"; HtmlRenderer applies defaults at render time.
type Theme struct {
	PrimaryColor     string
	SecondaryColor   string
	BackgroundColor  string // page/body background
	CardColor        string // element card/tile background
	TextColor        string
	AccentColor      string
	FontFamily       string
	BorderRadius     string
	ChartColors      []string
	ShadowIntensity  string // "none"|"subtle"|"medium"|"strong"
	EnableAnimations bool
	EnableGradients  bool
}

// DefaultTheme returns a Theme populated with the standard design defaults.
func DefaultTheme() *Theme {
	return &Theme{
		PrimaryColor:     "#2563eb",
		SecondaryColor:   "#64748b",
		BackgroundColor:  "#f1f5f9",
		CardColor:        "#ffffff",
		TextColor:        "#1e293b",
		AccentColor:      "#10b981",
		FontFamily:       `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`,
		BorderRadius:     "0.5rem",
		ShadowIntensity:  "medium",
		EnableAnimations: true,
	}
}

// htmlRenderContext is the rendering environment passed to each element's renderHTML method.
type htmlRenderContext struct {
	Theme        *Theme
	SectionTitle string
	idGen        *idGen
}

// nextID returns a stable, unique HTML element ID for a chart canvas with the given
// title within the current section. Call it once per chart element.
func (ctx *htmlRenderContext) nextID(elementTitle string) string {
	return ctx.idGen.next(ctx.SectionTitle, elementTitle)
}

// chartColors returns the active color palette: the theme's ChartColors when set,
// otherwise the package defaults.
func (ctx *htmlRenderContext) chartColors() []string {
	return chartColors(ctx.Theme)
}

// htmlRenderer is the render-dispatch interface for HTML output.
// Every Element type in this package implements it; renderElement uses it for dispatch.
type htmlRenderer interface {
	renderHTML(ctx *htmlRenderContext) (html string, scripts []string, err error)
}

// Renderer generates a report document from a Report and an optional Theme.
// Implementations must call DefaultTheme() when theme is nil.
type Renderer interface {
	Render(w io.Writer, report *Report, theme *Theme) error
}

// Ptr returns a pointer to v. Use it to set optional pointer fields in
// ChartOptions and similar structs without declaring a temporary variable:
//
//	opts.ShowTooltips = rptgen.Ptr(false)
//	opts.YMin = rptgen.Ptr(0.0)
func Ptr[T any](v T) *T { return &v }
