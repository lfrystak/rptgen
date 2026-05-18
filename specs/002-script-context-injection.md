# 002 — User data can break out of the `<script>` block (output corruption / XSS)

**Problem:** Chart labels and series names are JSON-marshaled directly into an inline `<script>` element without HTML-context-safe escaping, so a value containing `</script>` (or `<!--`, `<![CDATA[`) corrupts the document and enables stored XSS.

## Background / context

`chartInitScript` (`html.go:303`) builds chart JavaScript by `json.Marshal`-ing a config struct and string-interpolating it into a script body:

```go
configJSON, err := json.Marshal(config)
...
return fmt.Sprintf("(function(){var ctx=document.getElementById('%s').getContext('2d');new Chart(ctx,%s);})();",
    id, string(configJSON))
```

The marshaled config embeds **user-controlled strings**: chart labels and series names come straight from caller maps/slices (`BarChart.Data` keys, `LineSeries.Name`, `StackedBarSeries.Values` keys, etc. — see `charts.go`). Go's `encoding/json` does **not** escape `<`, `>`, `/` in a way that is safe for an HTML `<script>` context: a label such as `</script><script>alert(1)</script>` is emitted verbatim, terminating the script element early. Result: broken report at best, executed attacker script at worst, for reports built from untrusted data (a common reporting use case).

Related, lower-severity facets of the same "raw string building instead of context-aware templating" root cause:

- `FreeText` with `IsHTML: true` (`html.go:213`) injects `Content` unescaped. This is documented and intentional, but there is no `// SECURITY:` warning at the field/code site (`elements.go:94`) and no guidance that callers must sanitize.
- `renderElement`'s `id` is escaped with `html.EscapeString` (`html.go:273`) but then placed inside a **JavaScript single-quoted string** (`getElementById('%s')`); HTML escaping is the wrong escaping for that context. Currently safe only because `slugify` already restricts the charset — i.e. safety is accidental, not enforced at the injection point.

## Severity

**Critical** — stored XSS / document corruption from ordinary user data flowing into reports.

## Proposed change / acceptance criteria

1. Make the script-context emission safe regardless of label content. Acceptable approaches:
   - After `json.Marshal`, replace `<`→`<`, `>`→`>`, `&`→`&`, and ` `/` `, **or**
   - Render the chart config into a `<script type="application/json">` element (HTML-escaped via `html/template`) and have a small bootstrap read+`JSON.parse` it.
2. Add a regression test: a chart whose label is `</script><img src=x onerror=alert(1)>` must produce output where that substring is not present unescaped and the document still has exactly one closing `</script>` per opened script.
3. Add an explicit `// SECURITY:` doc comment on `FreeText.IsHTML` (`elements.go`) and in the README stating the caller is responsible for sanitizing raw HTML.
4. Stop relying on `slugify` for JS-context safety incidentally — either keep `slugify` and document it as a security boundary, or escape the id for the JS string context explicitly.
5. Consider switching HTML emission to `html/template` / `text/template` (see [015-html-monolith-refactor.md](015-html-monolith-refactor.md)) so context-aware escaping is structural rather than manual.

## Dependencies

- Best implemented alongside the templating refactor in [015-html-monolith-refactor.md](015-html-monolith-refactor.md).
- Test additions coordinate with [016-test-quality.md](016-test-quality.md).
