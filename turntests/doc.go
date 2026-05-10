// Package turntests holds tests that drive a single turn of the simulator from a fixed
// hand and assert the resulting per-turn Value. Public deck-level entry points only (no
// exports_test.go hacks). See docs/dev-standards.md "Test layout" for the convention.
package turntests

// Populates sim's forward-declared hooks before any test runs. See docs/dev-standards.md
// "Registry / sim split".
import _ "github.com/tim-chaplin/fab-deck-optimizer/internal/simreg"
