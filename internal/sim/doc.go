// Package sim is the hand-and-deck evaluator at the heart of the optimizer.
//
// Given a deck it shuffles, walks hands, and for each hand brute-forces the optimal turn:
// the partition of the hand into roles (Pitch, Attack, Defend, Held, Arsenal) and the
// attack-chain ordering that maximises turn value (damage dealt plus damage prevented).
//
// The search runs in two layers:
//
//   - Partition enumeration (partition.go) walks every role assignment and hands each leaf
//     to the attack-chain search.
//   - Attack-chain search (sequence.go) enumerates phase / weapon masks and permutes the
//     resulting attackers via playSequenceWithMeta, replaying one ordering through a pooled
//     GameState while firing triggers, Aura handlers, and OnHit closures.
//
// Best (hand_eval.go) is the single-turn entry point; Evaluate (deck_eval.go) folds per-hand
// results into deck.Stats. EvalOneTurnForTesting / EvalTwoTurnsForTesting are the public
// test entry points. The Evaluator type owns per-goroutine scratch buffers (attackbufs.go)
// so concurrent callers each get their own alloc-free state.
//
// See README.md in this directory for the full walkthrough.
package sim
