package rptgen

import "testing"

func TestNewLineChartSingle(t *testing.T) {
	pts := map[string]float64{"Jan": 10, "Feb": 20}
	lc := NewLineChartSingle("Sales", pts)

	if lc.Title != "Sales" {
		t.Errorf("Title: got %q", lc.Title)
	}
	if len(lc.Series) != 1 {
		t.Fatalf("Series: got len %d, want 1", len(lc.Series))
	}
	if lc.Series[0].Name != "Sales" {
		t.Errorf("Series[0].Name: got %q, want %q", lc.Series[0].Name, "Sales")
	}
	if len(lc.Series[0].Points) != 2 {
		t.Errorf("Series[0].Points: got len %d, want 2", len(lc.Series[0].Points))
	}
	if !lc.ShowPoints {
		t.Error("ShowPoints must be true")
	}
}

func TestNewLineChart(t *testing.T) {
	series := []LineSeries{
		{Name: "A", Points: map[string]float64{"Q1": 1}},
		{Name: "B", Points: map[string]float64{"Q1": 2}},
	}
	lc := NewLineChart("Trends", series)

	if len(lc.Series) != 2 {
		t.Fatalf("Series: got len %d, want 2", len(lc.Series))
	}
	if lc.Series[0].Name != "A" || lc.Series[1].Name != "B" {
		t.Errorf("Series order: got [%s %s], want [A B]", lc.Series[0].Name, lc.Series[1].Name)
	}
}

func TestElementTypeStrings(t *testing.T) {
	cases := []struct {
		elem Element
		want string
	}{
		{&NumberTile{BaseElement: newBaseElement()}, "NumberTile"},
		{&DateTile{BaseElement: newBaseElement()}, "DateTile"},
		{&FreeText{BaseElement: newBaseElement()}, "FreeText"},
		{&Table{BaseElement: newBaseElement()}, "Table"},
		{&Canvas{BaseElement: newBaseElement()}, "Canvas"},
		{NewBarChart("", nil), "BarChart"},
		{NewLineChart("", nil), "LineChart"},
		{NewPieChart("", nil), "PieChart"},
		{NewStackedBarChart("", nil), "StackedBarChart"},
	}
	for _, tc := range cases {
		got := tc.elem.ElementType()
		if got != tc.want {
			t.Errorf("%T.ElementType(): got %q, want %q", tc.elem, got, tc.want)
		}
	}
}
