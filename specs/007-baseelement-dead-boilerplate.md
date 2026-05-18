# 007 — `BaseElement` / `newBaseElement()` is dead boilerplate (C#-base-class carryover)

**Problem:** Every element type embeds an empty `BaseElement` struct and every struct literal is expected to set `BaseElement: newBaseElement()`, but `BaseElement` carries no fields, no methods, and contributes nothing to satisfying the `Element` interface.

## Background / context

`rptgen.go:13`:

```go
type BaseElement struct{}
func newBaseElement() BaseElement { return BaseElement{} }
```

It is embedded in `NumberTile`, `DateTile`, `FreeText`, `Table`, `Canvas`, and (via `ChartBase`) every chart. The `Element` interface is satisfied solely by each concrete type's own `ElementType()` method — `BaseElement` has no methods, so it does not help satisfy the interface and provides no shared behaviour or state.

Costs:

- Boilerplate everywhere: `&NumberTile{BaseElement: newBaseElement(), …}` throughout `*_test.go`, and `newBaseElement()` calls in `elements.go` constructors. `BaseElement{}` is the zero value anyway, so `newBaseElement()` is pure ceremony.
- It is an inheritance idiom imported from the C# original; in Go the empty embed is noise. New element types must remember to embed it for no reason, and forgetting it changes nothing — proving it is inert.

## Severity

**Important** — pervasive, idiom-violating boilerplate that taxes every element definition, test, and example; cheap to remove.

## Proposed change / acceptance criteria

1. Remove `BaseElement` and `newBaseElement()` entirely **or**, if a genuine shared concern emerges (e.g. an unexported sealing method per [005-extensibility-element-interface.md](005-extensibility-element-interface.md)), repurpose it to actually carry that — not an empty struct.
2. Drop the `BaseElement` embed from all element/chart structs and `ChartBase`.
3. Remove all `BaseElement: newBaseElement()` from constructors (`elements.go`, `charts.go`), tests (`*_test.go`), README, and example.
4. Tests still pass; `Element` is still satisfied by `ElementType()`.

## Dependencies

- If [005-extensibility-element-interface.md](005-extensibility-element-interface.md) chooses a sealed-interface marker, implement that here instead of a no-op embed (coordinate to avoid double churn).
