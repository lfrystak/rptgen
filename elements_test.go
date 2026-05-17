package rptgen

import (
	"testing"
	"time"
)

// --- NumberTile.FormatValue ---

func TestNumberTileFormatValue(t *testing.T) {
	cases := []struct {
		name   string
		value  float64
		format string
		want   string
	}{
		{"empty format", 42.5, "", "42.5"},
		{"N2 two decimals", 1234.5, "N2", "1234.50"},
		{"C0 currency no decimals", 1500, "C0", "$1,500"},
		{"P1 percentage 1 decimal", 0.753, "P1", "75.3%"},
		{"custom sprintf", 3.14159, "%.4f", "3.1416"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			n := &NumberTile{BaseElement: newBaseElement(), Value: tc.value, Format: tc.format}
			got := n.FormatValue("")
			if got != tc.want {
				t.Errorf("FormatValue(%q, %v): got %q, want %q", tc.format, tc.value, got, tc.want)
			}
		})
	}
}

// --- DateTile.FormatValue ---

func TestDateTileFormatValue(t *testing.T) {
	ts := time.Date(2024, 3, 15, 9, 5, 0, 0, time.UTC)

	t.Run("empty format uses default layout", func(t *testing.T) {
		d := &DateTile{BaseElement: newBaseElement(), Value: ts, Format: ""}
		got := d.FormatValue()
		if got != "2024-03-15 09:05:00" {
			t.Errorf("got %q", got)
		}
	})

	t.Run("date-only layout", func(t *testing.T) {
		d := &DateTile{BaseElement: newBaseElement(), Value: ts, Format: "2006-01-02"}
		got := d.FormatValue()
		if got != "2024-03-15" {
			t.Errorf("got %q", got)
		}
	})

	t.Run("zero time returns empty string", func(t *testing.T) {
		d := &DateTile{BaseElement: newBaseElement()}
		got := d.FormatValue()
		if got != "" {
			t.Errorf("got %q, want empty", got)
		}
	})
}

// --- Table constructors ---

func TestNewTable(t *testing.T) {
	rows := []map[string]any{
		{"banana": 1, "apple": 2},
		{"banana": 3, "apple": 4},
	}
	tbl := NewTable("T", rows)

	if tbl.Title != "T" {
		t.Errorf("Title: got %q", tbl.Title)
	}
	if len(tbl.Columns) != 2 {
		t.Fatalf("Columns: got len %d, want 2", len(tbl.Columns))
	}
	// columns sorted alphabetically
	if tbl.Columns[0] != "apple" || tbl.Columns[1] != "banana" {
		t.Errorf("Columns order: got %v, want [apple banana]", tbl.Columns)
	}
	if len(tbl.Rows) != 2 {
		t.Errorf("Rows: got len %d, want 2", len(tbl.Rows))
	}
}

func TestNewTableWithColumns(t *testing.T) {
	rows := []map[string]any{{"a": 1, "b": 2}}
	cols := []string{"b", "a"}
	tbl := NewTableWithColumns("T", rows, cols)

	if len(tbl.Columns) != 2 {
		t.Fatalf("Columns: got len %d, want 2", len(tbl.Columns))
	}
	if tbl.Columns[0] != "b" || tbl.Columns[1] != "a" {
		t.Errorf("Columns order: got %v, want [b a]", tbl.Columns)
	}
}

func TestNewTableFromColumns(t *testing.T) {
	data := map[string][]any{
		"x": {10, 20, 30},
		"y": {"a", "b", "c"},
	}
	tbl := NewTableFromColumns("T", data)

	if len(tbl.Rows) != 3 {
		t.Errorf("Rows: got len %d, want 3", len(tbl.Rows))
	}
}

func TestNewTableFromColumnsPanicsOnMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on mismatched column lengths")
		}
	}()
	NewTableFromColumns("T", map[string][]any{
		"x": {1, 2},
		"y": {1},
	})
}
