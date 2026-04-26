// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
//
// Entry points are Best / BestWithTriggers (evaluator.go): they partition a hand across five
// roles (Pitch, Attack, Defend, Held, Arsenal) and return the TurnSummary (types.go) with the
// highest Value.
//
// Internally the search runs in two layers:
//
//   - Partition enumeration (partition.go) walks every role assignment and hands each leaf to
//     bestAttackWithWeapons.
//   - Attack-chain search (sequence.go) enumerates phase / weapon masks and permutes the
//     resulting attackers via playSequenceWithMeta, which replays one ordering through a
//     pooled TurnState while firing hero triggers and AuraTrigger / EphemeralAttackTrigger
//     handlers. Per-card damage / block / pitch attribution is read off the chain's LogEntry
//     stream.
//
// The Evaluator type owns per-goroutine scratch buffers (attackbufs.go) so concurrent callers
// each get their own alloc-free state. Per-card metadata (cardmeta.go) is cached lazily into a
// uint16-keyed table so the chain inner loop avoids interface dispatch on Types / GoAgain /
// Cost.
//
// Format-layer helpers (format.go) render a TurnSummary for display: FormatBestLine is the
// compact one-liner, FormatBestTurn is the sectioned play-order printout.
package hand
