// Package sim evaluates the value of a hand of Flesh and Blood cards played in isolation,
// then folds per-hand results into per-deck statistics.
//
// Entry points are Best / BestWithTriggers (evaluator.go): they partition a hand across five
// roles (Pitch, Attack, Defend, Held, Arsenal) and return the TurnSummary with the highest
// Value.
//
// The search runs in two layers:
//
//   - Partition enumeration (partition.go) walks every role assignment and hands each leaf
//     to bestAttackWithWeapons.
//   - Attack-chain search (sequence.go) enumerates phase / weapon masks and permutes the
//     resulting attackers via playSequenceWithMeta, replaying one ordering through a pooled
//     TurnState while firing hero triggers, Aura handlers, and OnHit closures. Per-card
//     damage / block / pitch attribution is read off the chain's LogEntry stream.
//
// The Evaluator type owns per-goroutine scratch buffers (attackbufs.go) so concurrent
// callers each get their own alloc-free state. Per-card metadata (cardmeta.go) is cached
// lazily into a uint16-keyed table so the chain inner loop avoids interface dispatch.
//
// Format-layer helpers render the winning BestLine for display: FormatBestLine is the
// compact one-liner. PrintBestTurn re-runs Best against an eval-time snapshot and streams
// the sectioned play-order printout to an io.Writer.
package sim
