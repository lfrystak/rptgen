# Spec 02 — Elements

**File(s) to create:** `elements.go`  
**Depends on:** spec 01 (core types)  
**Referenced by:** spec 04 (HTML renderer), spec 05 (JSON renderer), spec 06 (tests)

---

## Goal

Implement all non-chart report elements. Each struct embeds `BaseElement` and implements the `Element` interface.

---

## NumberTile

Displays a single numeric metric.

```go
type NumberTile struct {
    BaseElement
    Title    string
    Value    float64
    Format   string  // Go fmt verb or pattern, e.g. "%.2f", "currency", "percent". Empty = default float formatting.
    Subtitle string  // optional secondary label
    Tooltip  string  // optional tooltip text shown on hover
}

func (n *NumberTile) ElementType() string { return "NumberTile" }

// FormatValue returns the display string for Value using Format and locale.
// locale is the Report.Locale string (empty = invariant).
func (n *NumberTile) FormatValue(locale string) string { ... }
```

**Format conventions** (keep it simple — no i18n library needed):
- `""` → `strconv.FormatFloat(value, 'f', -1, 64)` (Go default)
- `"N"` or `"N2"` → fixed 2 decimal places (or N digits if specified): `fmt.Sprintf("%.2f", value)`
- `"C"` or `"C0"` → prefix with `$` (or locale symbol), commas for thousands (basic implementation)
- `"P"` or `"P1"` → multiply by 100, append `%`
- any other string → treat as a `fmt.Sprintf` format pattern directly

> **Implementation note:** The C# version used `CultureInfo` for full locale-aware formatting. The Go version should produce reasonable output without pulling in `golang.org/x/text`. A currency prefix of `$` is acceptable; the key is that the pattern system works. More sophisticated locale formatting can be added later.

---

## DateTile

Displays a date or datetime metric.

```go
type DateTile struct {
    BaseElement
    Title    string
    Value    time.Time // zero value = not set
    Format   string    // Go time layout string, e.g. "2006-01-02". Empty = "2006-01-02 15:04:05"
    Subtitle string
    Tooltip  string
}

func (d *DateTile) ElementType() string { return "DateTile" }

func (d *DateTile) FormatValue() string {
    if d.Value.IsZero() { return "" }
    layout := d.Format
    if layout == "" { layout = "2006-01-02 15:04:05" }
    return d.Value.Format(layout)
}
```

> The C# version had separate `DateTime` and `DateOnly` overloads. In Go a single `time.Time` covers both. Callers wanting date-only output pass `Format: "2006-01-02"`.

---

## FreeText

Displays a block of text or raw HTML.

```go
type FreeText struct {
    BaseElement
    Content string
    IsHTML  bool  // if true, Content is rendered as raw HTML (not escaped)
}

func (f *FreeText) ElementType() string { return "FreeText" }
```

> When `IsHTML` is false the renderer must HTML-escape `Content`. When true it is inserted verbatim — callers are responsible for sanitizing.

---

## Table

Displays tabular data.

```go
type Table struct {
    BaseElement
    Title   string
    Columns []string                   // ordered column names
    Rows    []map[string]any           // each row is a map from column name to value
}

func (t *Table) ElementType() string { return "Table" }
```

Constructors:

```go
// NewTable infers Columns from the keys of the first row.
func NewTable(title string, rows []map[string]any) *Table { ... }

// NewTableWithColumns uses the provided column order explicitly.
func NewTableWithColumns(title string, rows []map[string]any, columns []string) *Table { ... }

// NewTableFromColumns converts column-oriented data to row-oriented.
// Input: map of columnName → slice of values (all slices must be the same length).
func NewTableFromColumns(title string, columns map[string][]any) *Table { ... }
```

> `map[string]any` is used instead of `Dictionary<string, object>`. Column order in `NewTable` is derived from the first row's key order — but since Go map iteration is non-deterministic, `NewTable` should sort column names alphabetically for consistent output. `NewTableWithColumns` gives callers full control.

---

## Canvas

A flexible sub-grid container that holds other elements.

```go
type Canvas struct {
    BaseElement
    ColumnWidths []int     // same semantics as Section.ColumnWidths
    Elements     []Element
}

func (c *Canvas) ElementType() string { return "Canvas" }

func (c *Canvas) AddElement(e Element) *Canvas { ... }  // returns c for chaining
```

Constructor:

```go
func NewCanvas(columnWidths ...int) *Canvas {
    return &Canvas{ColumnWidths: columnWidths}
}
```

---

## Acceptance Criteria

- All five types (`NumberTile`, `DateTile`, `FreeText`, `Table`, `Canvas`) compile and implement `Element`
- `ElementType()` returns the exact strings listed above (renderers switch on these)
- `NumberTile.FormatValue` handles `""`, `"N"/"N2"`, `"C"/"C0"`, `"P"/"P1"`, and arbitrary `fmt.Sprintf` patterns
- `DateTile.FormatValue` uses Go time layouts; zero `time.Time` returns `""`
- `Table.Columns` is always deterministically ordered (sorted when inferred from map)
- `NewTableFromColumns` panics if column slices are different lengths
- `go build ./...` and `go vet ./...` pass
