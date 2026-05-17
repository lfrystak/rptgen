# Spec 01 — Core Types

**File(s) to create:** `rptgen.go`  
**Depends on:** nothing  
**Referenced by:** all other specs

---

## Goal

Define the foundational types that the rest of the library builds on:
- `Element` interface
- `BaseElement` struct (shared fields)
- `Report` struct
- `Section` struct
- `Theme` struct
- `Renderer` interface

These types must compile on their own (`go build ./...` passes) even before the element and renderer implementations exist. Stub out the renderer files if needed.

---

## Element Interface

```go
type Element interface {
    elementID() string   // unexported — only called by renderers
    ElementType() string
}
```

An unexported `elementID()` keeps IDs internal; renderers call it when generating Chart.js canvas IDs and HTML `id` attributes.

---

## BaseElement

All concrete element types embed this:

```go
type BaseElement struct {
    id string
}

func newBaseElement() BaseElement {
    return BaseElement{id: uuid.New().String()}  // or generate without external dep (see note)
}

func (b BaseElement) elementID() string { return b.id }
```

> **Note on UUID:** Do NOT pull in `github.com/google/uuid` or similar. Instead, generate a short random ID using `crypto/rand` + `encoding/hex`, e.g. 8 random bytes → 16-char hex string. This keeps the module dependency-free.

---

## Report

```go
type Report struct {
    Title       string
    Sections    []*Section
    GeneratedAt time.Time  // set to time.Now() when not provided
    Footer      string     // optional
    LogoURL     string     // optional
    Locale      string     // e.g. "en-US", "de-DE", "" = invariant (default)
}
```

Constructor:

```go
func NewReport(title string) *Report {
    return &Report{
        Title:       title,
        GeneratedAt: time.Now(),
    }
}
```

Method to add a section:

```go
func (r *Report) AddSection(s *Section) *Report { ... }  // returns r for optional chaining
```

---

## Section

```go
type Section struct {
    Title        string
    Elements     []Element
    ColumnWidths []int  // proportional widths, e.g. [1,2] = 33%/67%. Nil = single column.
}
```

> `ColumnWidths` replaces the C# `Columns int + List<int>? ColumnWidths` duality. A nil or empty slice means one column; a slice `[1,1,1]` means three equal columns; `[1,2]` means two columns at 1/3 and 2/3 width.

Method:

```go
func (s *Section) AddElement(e Element) *Section { ... }  // returns s for optional chaining
```

---

## Theme

All fields have zero-value defaults (empty string = "use default"). The `HtmlRenderer` applies defaults at render time, not at struct construction.

```go
type Theme struct {
    PrimaryColor    string   // default: "#2563eb"
    SecondaryColor  string   // default: "#64748b"
    BackgroundColor string   // default: "#ffffff"
    TextColor       string   // default: "#1e293b"
    AccentColor     string   // default: "#10b981"
    FontFamily      string   // default: system font stack
    BorderRadius    string   // default: "0.5rem"
    ChartColors     []string // default: built-in palette
    ShadowIntensity string   // "none"|"subtle"|"medium"|"strong"; default: "medium"
    EnableAnimations bool    // default: true
    EnableGradients  bool    // default: false
}

func DefaultTheme() *Theme {
    return &Theme{
        PrimaryColor:     "#2563eb",
        SecondaryColor:   "#64748b",
        BackgroundColor:  "#ffffff",
        TextColor:        "#1e293b",
        AccentColor:      "#10b981",
        FontFamily:       `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`,
        BorderRadius:     "0.5rem",
        ShadowIntensity:  "medium",
        EnableAnimations: true,
    }
}
```

---

## Renderer Interface

```go
type Renderer interface {
    Render(report *Report, theme *Theme) (string, error)
}
```

`theme` may be `nil`; implementations must call `DefaultTheme()` if `nil` is passed.

---

## Acceptance Criteria

- `go build ./...` passes with no errors
- `go vet ./...` passes
- `Report`, `Section`, `Theme`, `Element`, `Renderer` are all exported
- `BaseElement.elementID()` is unexported
- No external dependencies added to `go.mod`
