package rptgen

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Element is the interface all report elements implement.
// elementID is unexported so only renderers within the package can call it.
type Element interface {
	elementID() string
	ElementType() string
}

// BaseElement holds the random ID shared by all concrete element types.
type BaseElement struct {
	id string
}

func newBaseElement() BaseElement {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic("rptgen: crypto/rand unavailable: " + err.Error())
	}
	return BaseElement{id: hex.EncodeToString(b)}
}

func (b BaseElement) elementID() string { return b.id }

// Report is the top-level document container.
type Report struct {
	Title       string
	Sections    []*Section
	GeneratedAt time.Time
	Footer      string
	LogoURL     string
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

// Section groups elements into a layout row.
// ColumnWidths holds proportional widths (e.g. [1,2] = 33%/67%).
// Nil or empty means a single column.
type Section struct {
	Title        string
	Elements     []Element
	ColumnWidths []int
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
	PrimaryColor    string
	SecondaryColor  string
	BackgroundColor string // page/body background
	CardColor       string // element card/tile background
	TextColor       string
	AccentColor     string
	FontFamily      string
	BorderRadius    string
	ChartColors     []string
	ShadowIntensity string // "none"|"subtle"|"medium"|"strong"
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

// Renderer generates a report document from a Report and an optional Theme.
// Implementations must call DefaultTheme() when theme is nil.
type Renderer interface {
	Render(report *Report, theme *Theme) (string, error)
}
