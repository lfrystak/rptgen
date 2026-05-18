# 005 — `Element` interface is illusory: rendering is a hardcoded type switch

**Problem:** The public `Element` interface and `Renderer` interface suggest extensibility, but neither a new element type nor a new renderer can actually be added from outside the package — the renderer dispatches via a closed type switch on concrete unexported-logic functions.

## Background / context

- `Element` (`rptgen.go:8`) only requires `ElementType() string`.
- `renderElement` (`html.go:151`) is a `switch e := elem.(type)` over the package's own concrete types. Any `Element` implemented by an external package hits `default:` and renders as `<!-- unknown element -->` (`html.go:179`).
- The `Renderer` interface (`rptgen.go:104`) is documented in the README (`README.md:329`) as the extension point ("Custom renderers can be implemented by satisfying the `Renderer` interface"), but all the machinery a renderer needs — element dispatch, chart config structs, CSS generation, ID generation — is unexported and HTML-specific (`html.go`). A second renderer (PDF, Markdown) would have to re-implement the entire type switch and could not reuse anything.

Net effect: the two interfaces advertise an extensibility story the architecture does not support. Adding a new report element *or* a new output format requires editing core files, defeating the purpose of the interfaces.

## Severity

**Important** — extensibility is a stated design goal (interfaces + README) that is not met; it shapes every future feature.

## Proposed change / acceptance criteria

Two models were considered:

- **Option A (closed, honest):** Accept that the element set is closed. Make `Element` a sealed interface (e.g. an unexported marker method) so external types can't pretend to be elements, and document that new elements require a library change. Remove the misleading README claim about custom renderers unless renderer reuse is actually provided.
- **Option B (open, visitor):** Add a render-dispatch interface elements implement (e.g. `Render(ctx, w) ` or a visitor: `Accept(ElementVisitor)`), so a custom element supplies its own rendering and a custom renderer implements the visitor. Export the helpers a renderer realistically needs.

**Decision: implement Option B (or, at minimum, an explicit chart render seam).** Adding new chart types is a known near-term goal; with Option A every new chart type remains a surgical edit to the central `renderElement` type switch, which does not satisfy that goal. The open model makes a new element/chart type self-contained — it supplies its own rendering rather than being dispatched by a closed switch.

Acceptance criteria:

1. A documented, tested answer to "how do I add a new element type?" and "how do I add a new output format?" that matches the actual code.
2. The element set is open: `renderElement`'s closed type switch is replaced by a dispatch interface/visitor. The `default:` / unknown-element path still returns an error (coordinate with [004-error-handling.md](004-error-handling.md)) rather than silently emitting an HTML comment, and the README extensibility section is corrected ([013-documentation-drift.md](013-documentation-drift.md)).
3. At least one example custom `Element` rendered correctly by `HtmlRenderer` in a test.
4. **Chart-type extension acceptance test:** a new chart type (e.g. `ScatterChart`) can be added — struct, constructor, and its own render logic — and renders correctly **without modifying `renderElement` / the central dispatch**. This test is the concrete bar that "I might want to add more chart types" must clear. Its Chart.js config needs depend on the config model being open for extension — see [015-html-monolith-refactor.md](015-html-monolith-refactor.md).

## Dependencies

- Drives the package split in [015-html-monolith-refactor.md](015-html-monolith-refactor.md).
- Unknown-element behaviour shared with [004-error-handling.md](004-error-handling.md).
- README correction in [013-documentation-drift.md](013-documentation-drift.md).
