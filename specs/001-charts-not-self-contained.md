# 001 — Reports are not self-contained: Chart.js is loaded from a CDN

**Problem:** Any report containing a chart silently depends on a live internet connection and a third-party CDN, directly contradicting the library's core "self-contained, no external dependencies" promise.

## Background / context

`html.go:79` emits:

```go
b.WriteString("  <script src=\"https://cdn.jsdelivr.net/npm/chart.js\"></script>\n")
```

The script tag is emitted unconditionally whenever `Render` runs, and the chart init scripts call `new Chart(...)` against this remotely-loaded global.

This breaks the headline value proposition stated in:

- `README.md:12`: "render everything to a single HTML file with **no external dependencies**"
- `README.md:196`: "Charts are rendered using Chart.js, **embedded inline — no CDN calls required**"
- `README.md:323`: "It produces a fully self-contained HTML document — all CSS and Chart.js JavaScript are embedded inline."

Concrete consequences:

- Reports do not render charts offline, in air-gapped environments, or when archived/emailed.
- `cdn.jsdelivr.net` resolving an unpinned `npm/chart.js` (latest, no version) means a future Chart.js major release can break every previously generated report.
- It is a privacy/security concern: opening any report phones home to a third party.

This is also why the smoke test `TestHtmlSmokeTest` asserts `cdn.jsdelivr.net` is present — the test currently encodes the wrong behaviour.

## Severity

**Critical** — the central product claim is false and reports silently fail offline.

## Proposed change / acceptance criteria

1. Vendor a pinned Chart.js (UMD build, specific version, e.g. `chart.umd.min.js`) into the package, embedded via `//go:embed`, and inline it inside a `<script>` block in the rendered output.
2. The rendered HTML must contain **zero** external `src`/`href`/`@import`/network references for the default renderer. Verify with a test that greps the output for `http://`/`https://` and fails on any (excluding user-supplied `LogoURL`, which should be documented as the one user-controlled exception).
3. Pin and document the bundled Chart.js version (e.g. expose `rptgen.ChartJSVersion`); document how to upgrade it.
4. Update `TestHtmlSmokeTest` (`html_test.go:69`) to assert the absence of the CDN URL and the presence of inline Chart.js source.
5. If a CDN mode is still desired, make it an explicit opt-in renderer option (see [005-extensibility-element-interface.md](005-extensibility-element-interface.md) and [006-render-to-io-writer.md](006-render-to-io-writer.md) for the renderer-configuration mechanism) — never the default.

## Dependencies

- Documentation correction tracked in [013-documentation-drift.md](013-documentation-drift.md) (must be reconciled once behaviour is fixed).
- Renderer option plumbing relates to [008-api-consistency.md](008-api-consistency.md).
