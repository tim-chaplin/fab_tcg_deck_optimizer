package sim

// runningCarry encapsulates findBest's "track the best CarryState across an iteration,
// finalize once at exit" pattern. It owns:
//
//   - the per-iteration scratch buffer (so updates are allocation-free via
//     CarryState.CopyFrom);
//   - the tiebreaker scalars (value, futureValuePlayed, cardsDrawn) so callers don't
//     materialise the running TurnSummary on every leaf;
//   - the "any candidate promoted yet" flag so Finalize knows whether to clone the
//     scratch into the caller's output or leave the seed in place.
//
// Use:
//
//	r := runningCarry{scratch: &bufs.findBestCarryScratch, arsenal: arsenalSeed}
//	for each candidate {
//	    if !r.Beats(value, futureValuePlayed, cardsDrawn, willOccupy) { continue }
//	    r.Promote(value, futureValuePlayed, cardsDrawn, arsenal, &carry)
//	    // ... record any caller-side state (BestLine roles, swung weapons) keyed off this
//	    // promotion ...
//	}
//	r.Finalize(&out.State)  // clones scratch when something was promoted; no-op otherwise
//
// Beats / Promote split deliberately: the caller writes its OWN secondary trackers
// (BestLine roles, swung weapon names) only after the gate passes, so they stay
// consistent with the carry winner without the gate having to know about them.
type runningCarry struct {
	// scratch points at a per-Best-call attackBufs scratch (e.g. findBestCarryScratch).
	// Promote writes the candidate carry into it via CarryState.CopyFrom; Finalize
	// clones it on the way out so the returned value owns independent backing.
	scratch *CarryState

	// Tiebreaker scalars maintained per Promote. Compared against the candidate's
	// equivalents in Beats. Threading these as scalars avoids materialising a full
	// running TurnSummary inside the recurse.
	value, futureValuePlayed int
	// arsenal is the running winner's end-of-chain arsenal slot — read by Beats's
	// willOccupy tiebreaker, written by Promote with the leaf's arsenalAtChainStart.
	arsenal Card
	// cardsDrawn is the number of cards drawn during this chain. Read by Beats as a
	// tiebreaker just under value — drawing a card is a pretty good effect.
	cardsDrawn int

	// seen flips true on the first Promote. Finalize uses it to decide between
	// "clone the scratch into out" and "leave out's seed value alone."
	seen bool
}

// Beats reports whether (value, futureValuePlayed, cardsDrawn, willOccupy) describes a
// candidate that should displace the running winner. Defers the lexicographic tuple
// comparison to chainScoreCmp; on a tie there, displaces only when the candidate ends
// with the arsenal slot occupied AND the running winner doesn't.
func (r *runningCarry) Beats(
	value, futureValuePlayed, cardsDrawn int, willOccupy bool,
) bool {
	if !r.seen {
		// No candidate yet — any feasible leaf wins, regardless of stats.
		return true
	}
	if cmp := chainScoreCmp(value, cardsDrawn, futureValuePlayed, r.value, r.cardsDrawn, r.futureValuePlayed); cmp != 0 {
		return cmp > 0
	}
	// Both willOccupy bits read end-of-chain carry (arsenal slot or Hand cards a
	// post-hoc promotion can pull from); reading r.scratch keeps the comparison
	// symmetric across the running winner and the candidate.
	bestWillOccupy := r.arsenal != nil || len(r.scratch.Hand) > 0
	return willOccupy && !bestWillOccupy
}

// Promote records the candidate as the new running winner. carry's slice contents are
// copied into the scratch (allocation-free after the first sizing); arsenal overrides
// the snapshot's Arsenal so an arsenal-in card that stayed is preserved.
func (r *runningCarry) Promote(
	value, futureValuePlayed, cardsDrawn int,
	arsenal Card, carry *CarryState,
) {
	r.scratch.CopyFrom(carry)
	r.scratch.Arsenal = arsenal
	r.value = value
	r.futureValuePlayed = futureValuePlayed
	r.arsenal = arsenal
	r.cardsDrawn = cardsDrawn
	r.seen = true
}

// Finalize writes the running winner into out, cloning the scratch so out's slice
// fields own independent backing. No-op when no candidate was promoted; callers that
// want a default seeded out should write it before the iteration runs.
func (r *runningCarry) Finalize(out *CarryState) {
	if r.seen {
		*out = r.scratch.Clone()
	}
}
