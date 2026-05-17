# Spec 04 — HTML Renderer

**File(s) to create:** `html.go`, `assets/chartjs/chart.umd.min.js`  
**Depends on:** spec 01, 02, 03  
**Referenced by:** spec 06 (tests)

---

## Goal

Implement `HtmlRenderer`, which converts a `*Report` into a complete, self-contained HTML document string. Chart.js must be embedded in the binary via `go:embed` so reports work fully offline.

---

## Chart.js Asset

Download Chart.js v4.4.0 UMD minified build from the CDN and save it to:

```
assets/chartjs/chart.umd.min.js
```

Embed it in the binary:

```go
//go:embed assets/chartjs/chart.umd.min.js
var chartJSSource string
```

The HTML renderer inlines this script in a `<script>` tag — no CDN request needed at view time.

> **How to get the file:** `curl -sL https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js -o assets/chartjs/chart.umd.min.js`

---

## HtmlRenderer

```go
type HtmlRenderer struct{}

func (h *HtmlRenderer) Render(report *Report, theme *Theme) (string, error) { ... }
```

If `theme` is `nil`, use `DefaultTheme()`.

---

## Output Structure

The renderer produces a single HTML string with this structure:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{Report.Title}}</title>
  <style>{{generated CSS}}</style>
</head>
<body>
  <div class="report">
    <header class="report-header">
      <!-- optional logo, title, generated-at timestamp -->
    </header>

    {{for each Section}}
    <section class="report-section">
      <h2 class="section-title">{{Section.Title}}</h2>
      <div class="section-grid" style="--col-template: {{computed grid-template-columns}}">
        {{for each Element in Section}}
        <div class="element-wrapper">{{rendered element HTML}}</div>
        {{end}}
      </div>
    </section>
    {{end}}

    {{if Report.Footer}}
    <footer class="report-footer">{{Report.Footer}}</footer>
    {{end}}
  </div>

  <script>{{chartJSSource}}</script>
  <script>{{chart initialization scripts}}</script>
</body>
</html>
```

---

## CSS

Generate CSS from Theme fields. Required styles:

**Layout:**
- `.report` — max-width 1400px, centered, padding
- `.report-header` — flex row, space-between, bottom border
- `.report-header img` — logo sizing (max-height 48px)
- `.section-grid` — CSS Grid using `--col-template` custom property for `grid-template-columns`
- Responsive: single-column below 768px
- `.element-wrapper` — card-style box with border-radius, optional shadow

**Shadows** (by `Theme.ShadowIntensity`):
- `"none"` → no shadow
- `"subtle"` → `0 1px 2px rgba(0,0,0,0.05)`
- `"medium"` → `0 4px 6px rgba(0,0,0,0.07)` (default)
- `"strong"` → `0 10px 15px rgba(0,0,0,0.1)`

**Animations** (if `Theme.EnableAnimations`):
- CSS `@keyframes fadeIn` with `animation: fadeIn 0.3s ease`

**Print:**
- `@media print` — hide shadows, remove animations, force single-column

**Element-specific:**
- `.tile` — centered text, large value font
- `.tile-title` — muted secondary label
- `.tile-value` — primary color, large font (2rem+)
- `.tile-subtitle` — small muted text below value
- `.chart-container` — fixed height (350px default), relative position (required by Chart.js)
- `.table-wrapper` — overflow-x auto for wide tables
- `table` — full width, border-collapse
- `th` — primary background, white text
- `td` — alternating row colors using `:nth-child(even)`
- `.free-text` — normal prose styling

**Tooltip:** If an element has a non-empty `Tooltip`, wrap it in a container with a CSS `[data-tooltip]` attribute and use `:hover::after` to show it — no JS required.

**Gradients** (if `Theme.EnableGradients`):
- Apply a subtle linear gradient to the `.report-header` background using `PrimaryColor`

---

## Section Grid Layout

Translate `Section.ColumnWidths` to a CSS `grid-template-columns` value:

- `nil` or `[]` → `"1fr"` (single column)
- `[1, 1, 1]` → `"1fr 1fr 1fr"` (three equal columns)
- `[1, 2]` → `"1fr 2fr"` (one third / two thirds)
- `[2, 1, 1]` → `"2fr 1fr 1fr"`

The formula: join each weight with `"fr"` suffix and space-separate.

The same logic applies to `Canvas.ColumnWidths` for nested grids.

---

## Element Rendering

### NumberTile

```html
<div class="element tile number-tile" {{if Tooltip}}data-tooltip="{{Tooltip}}"{{end}}>
  <div class="tile-title">{{Title}}</div>
  <div class="tile-value">{{FormatValue(report.Locale)}}</div>
  {{if Subtitle}}<div class="tile-subtitle">{{Subtitle}}</div>{{end}}
</div>
```

### DateTile

Same structure as NumberTile, using `FormatValue()`.

### FreeText

```html
<div class="element free-text">
  {{if IsHTML}}{{Content (raw)}}{{else}}<p>{{html.EscapeString(Content)}}</p>{{end}}
</div>
```

### Table

```html
<div class="element table-wrapper">
  <h3 class="element-title">{{Title}}</h3>
  <table>
    <thead><tr>{{for col}}<th>{{col}}</th>{{end}}</tr></thead>
    <tbody>
      {{for row}}<tr>{{for col}}<td>{{row[col]}}</td>{{end}}</tr>{{end}}
    </tbody>
  </table>
</div>
```

Cell values: call `fmt.Sprint(value)` to convert `any` to string.

### Canvas

Render as a nested CSS grid (same column logic as Section):

```html
<div class="element canvas-grid" style="--col-template: {{computed}}">
  {{for each child element}}
  <div class="element-wrapper">{{rendered child}}</div>
  {{end}}
</div>
```

### Charts

Each chart gets a unique `<canvas>` element and an initialization script:

```html
<div class="element chart-container" {{if Tooltip}}data-tooltip="{{Tooltip}}"{{end}}>
  <h3 class="element-title">{{Title}}</h3>
  <canvas id="{{element.elementID()}}"></canvas>
</div>
```

Chart init script (collected and emitted as one `<script>` block at end of body):

```js
(function() {
  var ctx = document.getElementById('{{id}}').getContext('2d');
  new Chart(ctx, {{chartConfigJSON}});
})();
```

**Chart config generation** — use `encoding/json` to marshal a config struct to JSON. Define unexported Go structs for Chart.js config shapes (see spec 03 for reference shapes).

Apply `Theme.ChartColors` as the `backgroundColor` array in datasets if non-empty; otherwise use a built-in 8-color palette:

```go
var defaultChartColors = []string{
    "#2563eb", "#10b981", "#f59e0b", "#ef4444",
    "#8b5cf6", "#06b6d4", "#f97316", "#84cc16",
}
```

For `LineChart`, additionally set `borderColor` (same color), `fill: false`, and `tension: 0.4`.

---

## HTML Generation Approach

Use Go's `html/template` package for the outer document structure. For element rendering, write helper functions that return `template.HTML` so the template does not double-escape already-safe HTML.

Alternatively, use `strings.Builder` directly throughout — both are acceptable. The important constraint is that **user-supplied text content (titles, tile values, footer, free text when `IsHTML=false`) must be HTML-escaped**.

---

## Acceptance Criteria

- `HtmlRenderer{}.Render(report, nil)` returns valid HTML (parseable by `golang.org/x/net/html` or a basic string check)
- `chart.umd.min.js` is embedded — generated HTML contains `new Chart(` without any CDN URL
- All six element types render without error
- Theme fields are applied: primary color in tile values, chart colors, shadow, animations
- `Section.ColumnWidths` and `Canvas.ColumnWidths` produce correct `grid-template-columns`
- User text is HTML-escaped; `FreeText{IsHTML: true}` is inserted raw
- `go build ./...` and `go vet ./...` pass
