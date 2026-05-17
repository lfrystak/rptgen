# Spec 05 — Tests

**File(s) to create:** `rptgen_test.go`, `elements_test.go`, `charts_test.go`, `html_test.go`  
**Depends on:** spec 01–04 (all implementations complete)

---

## Goal

Write unit tests covering the core behavior of the library. Tests live alongside source files in `package rptgen` (whitebox — they can access unexported helpers if needed).

---

## Test Files and Coverage

### `elements_test.go`

**NumberTile.FormatValue:**
- Empty format → raw float string
- `"N2"` → 2 decimal places
- `"C0"` → currency-prefixed, no decimals
- `"P1"` → percentage with 1 decimal
- Custom format like `"%.4f"` → direct sprintf

**DateTile.FormatValue:**
- Empty format → default layout `"2006-01-02 15:04:05"`
- `"2006-01-02"` → date-only output
- Zero `time.Time` → empty string

**Table constructors:**
- `NewTable` with a two-row input: columns inferred and sorted alphabetically
- `NewTableWithColumns`: columns in the supplied order
- `NewTableFromColumns`: row count matches slice length, panics on mismatched lengths

---

### `charts_test.go`

**LineChart constructors:**
- `NewLineChartSingle` wraps data as single series with chart title as series name
- `NewLineChart` with two series preserves series order

**ElementType strings:**
- Assert `ElementType()` for all eight element types returns the exact expected string

---

### `html_test.go`

**Smoke test — full report render:**

Build a report with at least one of every element type, call `HtmlRenderer{}.Render(report, nil)`, assert:
- No error returned
- Output starts with `<!DOCTYPE html>`
- Contains the report title
- Contains `new Chart(` (Chart.js initialization present)
- Does NOT contain `cdn.jsdelivr.net` (no CDN references)

**Theme application:**
- Render with a theme setting `PrimaryColor: "#ff0000"` — assert `#ff0000` appears in the output

**Column widths:**
- Section with `ColumnWidths: []int{1, 2}` → output contains `grid-template-columns` with `1fr 2fr`

**HTML escaping:**
- `FreeText{Content: "<script>alert(1)</script>", IsHTML: false}` → output contains `&lt;script&gt;`, not the raw tag
- `FreeText{Content: "<b>bold</b>", IsHTML: true}` → output contains `<b>bold</b>` verbatim

**Nil theme:**
- `Render(report, nil)` must not panic and must apply default colors

---

### `rptgen_test.go`

**Report construction:**
- `NewReport("X")` → `Title == "X"`, `GeneratedAt` not zero, `Sections` empty (not nil)
- `AddSection` appends and returns same pointer

**Section.AddElement:**
- Appends element; returns same pointer (for chaining)

**DefaultTheme:**
- All expected defaults are set (primary color, shadow intensity, animations enabled)

**gridTemplateColumns helper** (or equivalent internal function):
- `nil` / `[]int{}` → `"1fr"`
- `[1,1,1]` → `"1fr 1fr 1fr"`
- `[2,1]` → `"2fr 1fr"`

---

## Running Tests

```
go test ./... -v
```

All tests must pass. Aim for no `t.Skip` or incomplete tests — only test what is implemented.

---

## Acceptance Criteria

- `go test ./...` passes with zero failures
- All element types covered by at least one test
- HTML escaping tested for both `IsHTML: false` and `IsHTML: true`
- No CDN references in rendered HTML (Chart.js embedded)
