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
	BaseElement
	Title        string
	Value        float64
	Format       string // fmt.Sprintf format, e.g. "%.2f", "%.0f", "%.1f%%". Empty = raw float.
	Prefix       string // prepended to the formatted value, e.g. "$", "€"
	ThousandsSep bool   // insert comma thousands separator into the integer part
	Subtitle     string
	Tooltip      string
}

func (n *NumberTile) ElementType() string { return "NumberTile" }

func (n *NumberTile) FormatValue() string {
	var s string
	if n.Format == "" {
		s = strconv.FormatFloat(n.Value, 'f', -1, 64)
	} else {
		s = fmt.Sprintf(n.Format, n.Value)
	}
	if n.ThousandsSep {
		s = addThousandsSep(s)
	}
	if n.Prefix != "" {
		s = n.Prefix + s
	}
	return s
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
	BaseElement
	Title    string
	Value    time.Time // zero value = not set
	Format   string    // Go time layout string, e.g. "2006-01-02". Empty = "2006-01-02 15:04:05"
	Subtitle string
	Tooltip  string
}

func (d *DateTile) ElementType() string { return "DateTile" }

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
	BaseElement
	Content string
	IsHTML  bool // if true, Content is rendered as raw HTML (not escaped)
}

func (f *FreeText) ElementType() string { return "FreeText" }

// Table displays tabular data.
type Table struct {
	BaseElement
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
		BaseElement: newBaseElement(),
		Title:       title,
		Columns:     cols,
		Rows:        rows,
	}
}

// NewTableWithColumns uses the provided column order explicitly.
func NewTableWithColumns(title string, rows []map[string]any, columns []string) *Table {
	return &Table{
		BaseElement: newBaseElement(),
		Title:       title,
		Columns:     columns,
		Rows:        rows,
	}
}

// NewTableFromColumns converts column-oriented data to row-oriented.
// Panics if column slices have different lengths.
func NewTableFromColumns(title string, columns map[string][]any) *Table {
	rowCount := -1
	colNames := make([]string, 0, len(columns))
	for name, vals := range columns {
		colNames = append(colNames, name)
		if rowCount == -1 {
			rowCount = len(vals)
		} else if len(vals) != rowCount {
			panic("rptgen: NewTableFromColumns: column slices have different lengths")
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
		BaseElement: newBaseElement(),
		Title:       title,
		Columns:     colNames,
		Rows:        rows,
	}
}

// Canvas is a flexible sub-grid container that holds other elements.
type Canvas struct {
	BaseElement
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
		BaseElement:  newBaseElement(),
		ColumnWidths: columnWidths,
	}
}
