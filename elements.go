package rptgen

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NumberTile displays a single numeric metric.
type NumberTile struct {
	Title  string
	Value  float64
	Format string // fmt.Sprintf format for a float64, e.g. "%.2f", "%.0f", "%.1f%%". Empty = raw decimal.
	// Only float verbs produce valid output (f, F, e, E, g, G, x, X, b); other verbs cause fmt to emit
	// its %!verb(...) error sentinel, which FormatValue surfaces verbatim so the problem is visible.
	Prefix       string // prepended to the formatted value, e.g. "$", "€"
	ThousandsSep bool   // insert comma thousands separator; only effective with decimal (f/F-verb) formats
	Subtitle     string
	Tooltip      string
}

func (n *NumberTile) ElementType() string { return "NumberTile" }

// NewNumberTile returns a NumberTile with the given title and value.
// Optional display fields (Format, Prefix, ThousandsSep, Subtitle, Tooltip) can be set on the returned pointer.
func NewNumberTile(title string, value float64) *NumberTile {
	return &NumberTile{
		Title: title,
		Value: value,
	}
}

func (n *NumberTile) FormatValue() string {
	var s string
	if n.Format == "" {
		s = strconv.FormatFloat(n.Value, 'f', -1, 64)
	} else {
		s = fmt.Sprintf(n.Format, n.Value)
		// Surface fmt's %! error sentinel immediately; further processing would only mangle it.
		if strings.Contains(s, "%!") {
			return s
		}
	}
	if n.ThousandsSep && isDecimalOutput(s) {
		s = addThousandsSep(s)
	}
	if n.Prefix != "" {
		s = n.Prefix + s
	}
	return s
}

// isDecimalOutput reports whether s looks like a plain decimal number (optional leading minus,
// digits, optional dot + digits, optional trailing percent). Only such strings are safe to pass
// to addThousandsSep; scientific, hex, and other non-decimal fmt outputs are not.
func isDecimalOutput(s string) bool {
	trimmed := strings.TrimSuffix(s, "%")
	trimmed = strings.TrimPrefix(trimmed, "-")
	if trimmed == "" {
		return false
	}
	for _, ch := range trimmed {
		if (ch < '0' || ch > '9') && ch != '.' {
			return false
		}
	}
	return true
}

// addThousandsSep inserts comma thousands separators into the integer part of a formatted number string.
func addThousandsSep(s string) string {
	neg := false
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}
	dot := strings.Index(s, ".")
	intPart, fracPart := s, ""
	if dot >= 0 {
		intPart, fracPart = s[:dot], s[dot:]
	}
	var b strings.Builder
	for i, ch := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	result := b.String() + fracPart
	if neg {
		return "-" + result
	}
	return result
}

// DateTile displays a date or datetime metric.
type DateTile struct {
	Title    string
	Value    time.Time // zero value = not set
	Format   string    // Go time layout string, e.g. "2006-01-02". Empty = "2006-01-02 15:04:05"
	Subtitle string
	Tooltip  string
}

func (d *DateTile) ElementType() string { return "DateTile" }

// NewDateTile returns a DateTile with the given title and time value.
// Optional display fields (Format, Subtitle, Tooltip) can be set on the returned pointer.
func NewDateTile(title string, value time.Time) *DateTile {
	return &DateTile{
		Title: title,
		Value: value,
	}
}

func (d *DateTile) FormatValue() string {
	if d.Value.IsZero() {
		return ""
	}
	layout := d.Format
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return d.Value.Format(layout)
}

// FreeText displays a block of text or raw HTML.
type FreeText struct {
	Content string
	// SECURITY: when IsHTML is true, Content is injected verbatim into the document without
	// any escaping. The caller is responsible for ensuring the value is safe HTML — i.e. it
	// must be sanitized before being passed here if it originates from untrusted input.
	IsHTML bool
}

func (f *FreeText) ElementType() string { return "FreeText" }

// NewFreeText returns a FreeText element with the given plain-text content.
// Set IsHTML = true on the returned pointer to inject Content as raw HTML (caller must ensure it is safe).
func NewFreeText(content string) *FreeText {
	return &FreeText{Content: content}
}

// Table displays tabular data.
type Table struct {
	Title   string
	Columns []string         // ordered column names
	Rows    []map[string]any // each row maps column name to value
}

func (t *Table) ElementType() string { return "Table" }

// NewTable infers Columns from the keys of the first row, sorted alphabetically.
func NewTable(title string, rows []map[string]any) *Table {
	var cols []string
	if len(rows) > 0 {
		for k := range rows[0] {
			cols = append(cols, k)
		}
		sort.Strings(cols)
	}
	return &Table{
		Title:   title,
		Columns: cols,
		Rows:    rows,
	}
}

// NewTableWithColumns uses the provided column order explicitly.
func NewTableWithColumns(title string, rows []map[string]any, columns []string) *Table {
	return &Table{
		Title:   title,
		Columns: columns,
		Rows:    rows,
	}
}

// NewTableFromColumns converts column-oriented data to row-oriented.
// Returns an error if column slices have different lengths.
func NewTableFromColumns(title string, columns map[string][]any) (*Table, error) {
	rowCount := -1
	colNames := make([]string, 0, len(columns))
	for name, vals := range columns {
		colNames = append(colNames, name)
		if rowCount == -1 {
			rowCount = len(vals)
		} else if len(vals) != rowCount {
			return nil, fmt.Errorf("rptgen: NewTableFromColumns: column slices have different lengths")
		}
	}
	sort.Strings(colNames)

	if rowCount < 0 {
		rowCount = 0
	}
	rows := make([]map[string]any, rowCount)
	for i := range rows {
		rows[i] = make(map[string]any, len(colNames))
		for _, col := range colNames {
			rows[i][col] = columns[col][i]
		}
	}

	return &Table{
		Title:   title,
		Columns: colNames,
		Rows:    rows,
	}, nil
}

// Canvas is a nestable sub-grid element that can be placed inside a Section column.
// It uses the same proportional ColumnWidths semantics as Section, but unlike Section
// it is itself an Element and can be nested at any depth.
type Canvas struct {
	ColumnWidths []int
	Elements     []Element
}

func (c *Canvas) ElementType() string { return "Canvas" }

// AddElement appends e to the canvas and returns c for chaining.
func (c *Canvas) AddElement(e Element) *Canvas {
	c.Elements = append(c.Elements, e)
	return c
}

// NewCanvas returns a Canvas with the given proportional column widths.
func NewCanvas(columnWidths ...int) *Canvas {
	return &Canvas{
		ColumnWidths: columnWidths,
	}
}
