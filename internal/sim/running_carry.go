package sim

// runningCarry encapsulates findBest's "track the best CarryState across an iteration,
// finalize once at exit" pattern. It owns:
//
//   - the per-iteration scratch buffer (so updates are allocation-free via
//     CarryState.CopyFrom);
//   - the tiebreaker scalars (Value, leftoverRunechants, futureValuePlayed) so callers
//     don't materialise the running TurnSummary on every leaf;
//   - the "any candidate promoted yet" flag so Finalize knows whether to clone the
//     scratch into the caller's output or leave the seed in place.
//
// Use:
//
//	r := runningCarry{
//	    scratch:            &bufs.findBestCarryScratch,
//	    arsenal:            arsenalSeed,
//	    hasHeld:            n > 0,
//	    leftoverRunechants: runechantCarryover,
//	}
//	for each candidate {
//	    if !r.Beats(value, leftoverRunechants, futureValuePlayed, willOccupy) { continue }
//	    r.Promote(value, leftoverRunechants, futureValuePlayed, hasHeld, arsenal, &carry)
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
	value, leftoverRunechants, futureValuePlayed int
	// arsenal is the running winner's end-of-chain arsenal slot — read by Beats's
	// willOccupy tiebreaker, written by Promote with the leaf's arsenalAtChainStart.
	arsenal Card
	// itemCarry sums end-of-chain Item.Count across all token types — used as the
	// final tiebreaker to prefer chains that spent items (drew cards now) over chains
	// that hoarded them. Symmetric to "use it or lose it" without forbidding carry.
	itemCarry int

	// seen flips true on the first Promote. Finalize uses it to decide between
	// "clone the scratch into out" and "leave out's seed value alone."
	seen bool
}

// Beats reports whether (value, leftoverRunechants, futureValuePlayed, willOccupy)
// describes a candidate that should displace the running winner. Tiebreaker order:
// higher Value wins; equal Value, higher leftoverRunechants wins; equal both, more
// futureValuePlayed wins; equal all three, only displace if the candidate ends with
// arsenal occupied AND the running winner doesn't.
func (r *runningCarry) Beats(
	value, leftoverRunechants, futureValuePlayed int, willOccupy bool, itemCarry int,
) bool {
	if !r.seen {
		// No candidate yet — any feasible leaf wins, regardless of stats.
		return true
	}
	if value != r.value {
		return value > r.value
	}
	if leftoverRunechants != r.leftoverRunechants {
		return leftoverRunechants > r.leftoverRunechants
	}
	if futureValuePlayed != r.futureValuePlayed {
		return futureValuePlayed > r.futureValuePlayed
	}
	// Both willOccupy bits read end-of-chain carry (arsenal slot or Hand cards a
	// post-hoc promotion can pull from); reading r.scratch keeps the comparison
	// symmetric across the running winner and the candidate.
	bestWillOccupy := r.arsenal != nil || len(r.scratch.Hand) > 0
	if willOccupy != bestWillOccupy {
		return willOccupy
	}
	// Final tiebreaker: chain that spent more items (left fewer in the carry) wins.
	// Symmetric to "use it or lose it" — both chains have equal damage and equal
	// next-turn surface area, so the one that cashed in a Gold token now (drew a
	// card into arsenal) outranks the one that hoarded it for next turn.
	return itemCarry < r.itemCarry
}

// Promote records the candidate as the new running winner. carry's slice contents are
// copied into the scratch (allocation-free after the first sizing); arsenal overrides
// the snapshot's Arsenal so an arsenal-in card that stayed is preserved.
func (r *runningCarry) Promote(
	value, leftoverRunechants, futureValuePlayed int,
	arsenal Card, carry *CarryState, itemCarry int,
) {
	r.scratch.CopyFrom(carry)
	r.scratch.Arsenal = arsenal
	r.value = value
	r.leftoverRunechants = leftoverRunechants
	r.futureValuePlayed = futureValuePlayed
	r.arsenal = arsenal
	r.itemCarry = itemCarry
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
