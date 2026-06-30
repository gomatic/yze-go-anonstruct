# yze-go-anonstruct

A [`yze`](https://github.com/gomatic/yze) analyzer (categories `types`/`structure`) enforcing the gomatic Go standard that struct types are named rather than anonymous. It flags anonymous struct types that carry fields — as variables, parameters, results, composite literals, struct fields, and map/slice element types; empty anonymous structs (`struct{}`, idiomatic for sets and signaling channels) and the right-hand side of a named type declaration are allowed.

- **Rule:** `yze/anonstruct`
- **Library:** exports `Analyzer` and `Registration` for the [`yze`](https://github.com/gomatic/yze) aggregator and [`stickler`](https://github.com/gomatic/stickler) runner.
- **Binary:** `cmd/yze-go-anonstruct` runs it standalone (`text`/`-json`, and as a `go vet -vettool`).

Built on the [`go-yze`](https://github.com/gomatic/go-yze) framework.
