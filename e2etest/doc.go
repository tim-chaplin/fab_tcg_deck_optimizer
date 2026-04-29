// Package e2etest holds end-to-end tests that exercise the simulator through its public
// entry points only — Best, EvaluateWith, and the deck / hero / card constructors —
// without reaching into sim internals via exports_test.go-style hacks. Tests living here
// are larger than unit tests (they construct a full hand + hero + chain runner state and
// assert on the optimizer's output) but smaller than integration smoke tests (no I/O, no
// CSV parsing, no shuffling driver). New cross-cutting behaviour tests should land here
// rather than in `internal/sim/*_test.go` so the sim package's test surface stays focused
// on internals, not end-to-end fixtures.
package e2etest
