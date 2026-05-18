package rptgen

import (
	"os/exec"
	"testing"
	"time"
)

func TestNewReport(t *testing.T) {
	before := time.Now()
	r := NewReport("My Report")
	after := time.Now()

	if r.Title != "My Report" {
		t.Errorf("Title: got %q, want %q", r.Title, "My Report")
	}
	if r.GeneratedAt.IsZero() {
		t.Error("GeneratedAt must not be zero")
	}
	if r.GeneratedAt.Before(before) || r.GeneratedAt.After(after) {
		t.Error("GeneratedAt not within expected range")
	}
	if r.Sections == nil {
		t.Error("Sections must be non-nil")
	}
	if len(r.Sections) != 0 {
		t.Errorf("Sections: got len %d, want 0", len(r.Sections))
	}
}

func TestAddSection(t *testing.T) {
	r := NewReport("R")
	s1 := &Section{Title: "S1"}
	s2 := &Section{Title: "S2"}

	returned := r.AddSection(s1).AddSection(s2)
	if returned != r {
		t.Error("AddSection must return same *Report for chaining")
	}
	if len(r.Sections) != 2 {
		t.Errorf("Sections: got len %d, want 2", len(r.Sections))
	}
	if r.Sections[0] != s1 || r.Sections[1] != s2 {
		t.Error("Sections order mismatch")
	}
}

func TestSectionAddElement(t *testing.T) {
	s := &Section{}
	e1 := &NumberTile{BaseElement: newBaseElement(), Title: "A"}
	e2 := &NumberTile{BaseElement: newBaseElement(), Title: "B"}

	returned := s.AddElement(e1).AddElement(e2)
	if returned != s {
		t.Error("AddElement must return same *Section for chaining")
	}
	if len(s.Elements) != 2 {
		t.Errorf("Elements: got len %d, want 2", len(s.Elements))
	}
	if s.Elements[0] != e1 || s.Elements[1] != e2 {
		t.Error("Elements order mismatch")
	}
}

func TestDefaultTheme(t *testing.T) {
	th := DefaultTheme()

	if th.PrimaryColor != "#2563eb" {
		t.Errorf("PrimaryColor: got %q", th.PrimaryColor)
	}
	if th.ShadowIntensity != "medium" {
		t.Errorf("ShadowIntensity: got %q", th.ShadowIntensity)
	}
	if !th.EnableAnimations {
		t.Error("EnableAnimations must be true")
	}
	if th.BackgroundColor == "" {
		t.Error("BackgroundColor must not be empty")
	}
	if th.TextColor == "" {
		t.Error("TextColor must not be empty")
	}
}

func TestEqualColumns(t *testing.T) {
	cases := []struct {
		n    int
		want []int
	}{
		{0, nil},
		{-1, nil},
		{1, []int{1}},
		{3, []int{1, 1, 1}},
		{4, []int{1, 1, 1, 1}},
	}
	for _, tc := range cases {
		got := EqualColumns(tc.n)
		if len(got) != len(tc.want) {
			t.Errorf("EqualColumns(%d): got len %d, want len %d", tc.n, len(got), len(tc.want))
			continue
		}
		for i, v := range got {
			if v != 1 {
				t.Errorf("EqualColumns(%d)[%d]: got %d, want 1", tc.n, i, v)
			}
		}
	}
}

func TestExampleBuilds(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "/dev/null", "./example/...")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("example/... failed to build: %v\n%s", err, out)
	}
}

func TestColumnWidthsToCSS(t *testing.T) {
	cases := []struct {
		widths []int
		want   string
	}{
		{nil, "1fr"},
		{[]int{}, "1fr"},
		{[]int{1, 1, 1}, "1fr 1fr 1fr"},
		{[]int{2, 1}, "2fr 1fr"},
	}
	for _, tc := range cases {
		got := columnWidthsToCSS(tc.widths)
		if got != tc.want {
			t.Errorf("columnWidthsToCSS(%v): got %q, want %q", tc.widths, got, tc.want)
		}
	}
}
