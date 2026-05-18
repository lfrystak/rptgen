package rptgen

import (
	"fmt"
	"strings"
)

// slugify converts a string to a lowercase hyphen-separated identifier.
// Runs of non-alphanumeric characters collapse to a single hyphen; leading/trailing hyphens are trimmed.
func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	prevDash := true
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
		} else if !prevDash {
			b.WriteByte('-')
			prevDash = true
		}
	}
	return strings.TrimRight(b.String(), "-")
}

// idGen produces unique, deterministic IDs for chart canvas elements within a single render.
type idGen struct{ seen map[string]int }

func newIDGen() *idGen { return &idGen{seen: make(map[string]int)} }

func (g *idGen) next(sectionTitle, chartTitle string) string {
	s, t := slugify(sectionTitle), slugify(chartTitle)
	var base string
	switch {
	case s != "" && t != "":
		base = s + "-" + t
	case t != "":
		base = t
	case s != "":
		base = s
	default:
		base = "chart"
	}
	g.seen[base]++
	if g.seen[base] == 1 {
		return base
	}
	return fmt.Sprintf("%s-%d", base, g.seen[base])
}
