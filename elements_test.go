package rptgen

import (
	"strings"
	"testing"
	"time"
)

// --- NumberTile.FormatValue ---

func TestNumberTileFormatValue(t *testing.T) {
	cases := []struct {
		name string
		tile NumberTile
		want string
	}{
		{"empty format", NumberTile{Value: 42.5}, "42.5"},
		{"fmt format two decimals", NumberTile{Value: 1234.5, Format: "%.2f"}, "1234.50"},
		{"thousands separator", NumberTile{Value: 1234.5, Format: "%.2f", ThousandsSep: true}, "1,234.50"},
		{"currency with prefix and thousands", NumberTile{Value: 1500, Format: "%.0f", Prefix: "$", ThousandsSep: true}, "$1,500"},
		{"percentage via fmt pattern", NumberTile{Value: 75.3, Format: "%.1f%%"}, "75.3%"},
		{"custom sprintf", NumberTile{Value: 3.14159, Format: "%.4f"}, "3.1416"},
		{"negative with thousands", NumberTile{Value: -1234.5, Format: "%.2f", ThousandsSep: true}, "-1,234.50"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.tile.BaseElement = newBaseElement()
			got := tc.tile.FormatValue()
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
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
	tbl, err := NewTableFromColumns("T", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tbl.Rows) != 3 {
		t.Errorf("Rows: got len %d, want 3", len(tbl.Rows))
	}
}

func TestNewTableFromColumnsErrorOnMismatch(t *testing.T) {
	_, err := NewTableFromColumns("T", map[string][]any{
		"x": {1, 2},
		"y": {1},
	})
	if err == nil {
		t.Fatal("expected error on mismatched column lengths, got nil")
	}
	if !strings.Contains(err.Error(), "different lengths") {
		t.Errorf("error should mention 'different lengths', got: %v", err)
	}
}

func TestNewTableFromColumnsRowValues(t *testing.T) {
	data := map[string][]any{
		"name":  {"Alice", "Bob"},
		"score": {95, 87},
	}
	tbl, err := NewTableFromColumns("T", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// columns sorted alphabetically: name, score
	if len(tbl.Columns) != 2 || tbl.Columns[0] != "name" || tbl.Columns[1] != "score" {
		t.Fatalf("Columns: got %v, want [name score]", tbl.Columns)
	}
	if len(tbl.Rows) != 2 {
		t.Fatalf("Rows: got %d, want 2", len(tbl.Rows))
	}
	if tbl.Rows[0]["name"] != "Alice" {
		t.Errorf("row[0][name]: got %v, want Alice", tbl.Rows[0]["name"])
	}
	if tbl.Rows[0]["score"] != 95 {
		t.Errorf("row[0][score]: got %v, want 95", tbl.Rows[0]["score"])
	}
	if tbl.Rows[1]["name"] != "Bob" {
		t.Errorf("row[1][name]: got %v, want Bob", tbl.Rows[1]["name"])
	}
	if tbl.Rows[1]["score"] != 87 {
		t.Errorf("row[1][score]: got %v, want 87", tbl.Rows[1]["score"])
	}
}

func TestNewTableFromColumnsEmpty(t *testing.T) {
	tbl, err := NewTableFromColumns("T", map[string][]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tbl.Rows) != 0 {
		t.Errorf("Rows: got %d, want 0", len(tbl.Rows))
	}
	if len(tbl.Columns) != 0 {
		t.Errorf("Columns: got %d, want 0", len(tbl.Columns))
	}
}

func TestNewTableFromColumnsNonStringAny(t *testing.T) {
	data := map[string][]any{
		"int_val":   {42, 100},
		"float_val": {3.14, 2.71},
	}
	tbl, err := NewTableFromColumns("T", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tbl.Rows) != 2 {
		t.Fatalf("Rows: got %d, want 2", len(tbl.Rows))
	}
	if tbl.Rows[0]["int_val"] != 42 {
		t.Errorf("row[0][int_val]: got %v (%T), want 42", tbl.Rows[0]["int_val"], tbl.Rows[0]["int_val"])
	}
	if tbl.Rows[0]["float_val"] != 3.14 {
		t.Errorf("row[0][float_val]: got %v (%T), want 3.14", tbl.Rows[0]["float_val"], tbl.Rows[0]["float_val"])
	}
}
