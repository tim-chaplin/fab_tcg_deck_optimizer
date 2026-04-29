// Package e2etest holds end-to-end tests that exercise the simulator through its public
// deck-level entry points — (*Deck).EvalOneTurnForTesting for chain evaluation,
// (*Deck).EvaluateWith for full multi-turn runs — without reaching into sim internals via
// exports_test.go-style hacks. Tests living here are larger than unit tests (they construct
// a real Deck + hero + initial hand and assert on the optimizer's output) but smaller than
// integration smoke tests (no I/O, no CSV parsing, no shuffling driver). New cross-cutting
// behaviour tests should land here rather than in `internal/sim/*_test.go` so the sim
// package's test surface stays focused on internals, not end-to-end fixtures. See
// docs/dev-standards.md "Test layout" for the full convention.
package e2etest
