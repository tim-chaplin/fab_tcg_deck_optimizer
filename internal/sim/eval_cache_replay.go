package sim

// Cache-hit replay: rebuild a TurnSummary from a cached winning partition by running the
// chain dispatcher against just that one role assignment. Skips the partition search
// (the dominant cost — exponential in hand size) but still runs bestAttackWithWeapons
// once and the post-hoc arsenal promotion + Hand carryover bookkeeping that findBest
// does after the search loop, so the resulting summary is byte-identical to a full
// from-scratch Best call.

import "fmt"

// replayBest is the cache-hit body. Thin wrapper around evaluatePartition: project the
// cached BestLine onto the new call's hand to fill rolesBuf, hand off to evaluatePartition
// for the actual chain run, then assemble the TurnSummary from its outputs. The cache
// search / store gate in findBest already guarantees priorAuraTriggers is empty and the
// hand multiset matches the entry's, so the chain output here is byte-identical to what
// the original cached call produced.
//
// One quirk in projecting the BestLine: the cached entry may tag a hand card with
// Role=Arsenal (the post-hoc promotion target). Hand cards never have that role during
// the chain run — the search treats them as Held and the post-hoc step re-flips at the
// end — so we flip the entry back to Held before evaluatePartition and re-stamp Arsenal
// on the BestLine afterward.
func (e *Evaluator) replayBest(
	entry evalCacheEntry,
	hero Hero, weapons []Weapon, hand []Card,
	incomingDamage int, deck []Card, runechantCarryover int,
	arsenalCardIn Card, skipLog bool,
) TurnSummary {
	n := len(hand)
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	bufs := e.getAttackBufs(n, weapons)
	rolesBuf := bufs.rolesBuf[:totalN]
	postPromotedFromHeld := -1
	if !mapCachedRolesToHand(entry.line, hand, arsenalCardIn, rolesBuf, &postPromotedFromHeld) {
		// Multiset mismatch can't happen by construction — the cache key sorts hand IDs
		// and the entry was stored under that exact key, so a hit means the multisets are
		// identical. Reaching here indicates a bug (key collision, cache corruption,
		// mid-call mutation of cachedLine, etc.) that's already compromised correctness;
		// panic loudly so the operator notices instead of falling back to a silent
		// re-search that hides the cache bug.
		panic(fmt.Sprintf("replayBest: mapCachedRolesToHand failed despite cache hit — cache invariant violated (hand=%d, cachedLine=%d, arsenal=%v)",
			len(hand), len(entry.line), arsenalCardIn != nil))
	}

	// Flip post-hoc-promoted hand entry from Arsenal back to Held for the chain run; the
	// promotion re-runs below and re-stamps Arsenal on the BestLine.
	if postPromotedFromHeld >= 0 {
		rolesBuf[postPromotedFromHeld] = Held
	}

	// defenseSum has to match what the original search computed — sum of Defense() across
	// every Defend-role card (DR or plain), per fillPartitionPerCardBufs. It feeds
	// state.BlockTotal so DR Plays that read "did we block all incoming?" see the right
	// shape. We compute it here rather than via fillPartitionPerCardBufs because the
	// recurse path's accumulator-arg threading isn't available — replay knows the role
	// assignment directly.
	defenseSum := defenseSumFromRoles(hand, arsenalCardIn, rolesBuf, n)

	attackDealt, defenseDealt, _, swung, carry, ok, _, arsenalAtChainStart := e.evaluatePartition(
		hero, weapons, hand, deck, arsenalCardIn,
		rolesBuf, n, bufs,
		runechantCarryover, incomingDamage, defenseSum,
		nil, skipLog,
	)
	if !ok {
		// Infeasible-partition replay can't happen by construction — the cached entry
		// was only stored after best.Cacheable=true with a feasible winning partition.
		// Reaching here means either the cache stored an infeasible result (bug) or some
		// "should be deterministic" input drifted. Panic so the operator notices rather
		// than silently re-searching and hiding a real correctness bug.
		panic(fmt.Sprintf("replayBest: cached partition is infeasible — cache invariant violated (hand=%d, runechantCarryover=%d, incomingDamage=%d)",
			len(hand), runechantCarryover, incomingDamage))
	}

	// Re-stamp the post-hoc-promoted entry's Arsenal role so the BestLine matches the
	// cached layout (and the post-promotion step below sees the same shape findBest did).
	if postPromotedFromHeld >= 0 {
		rolesBuf[postPromotedFromHeld] = Arsenal
	}

	// Build the TurnSummary. BestLine cards come from the new call's hand (so the printout
	// names the right Card values) but roles come from the cached entry. Mirror findBest's
	// final wiring: adopt CarryState, override Arsenal from arsenalAtChainStart so an
	// arsenal-in card that stayed is preserved, then re-do the post-hoc promotion.
	best := TurnSummary{
		BestLine:       make([]CardAssignment, totalN),
		Value:          attackDealt + defenseDealt,
		SwungWeapons:   append([]string(nil), swung...),
		IncomingDamage: incomingDamage,
		Cacheable:      true,
		State:          carry,
	}
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: rolesBuf[i]}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: rolesBuf[n], FromArsenal: true}
	}
	best.State.Arsenal = arsenalAtChainStart
	if best.State.Arsenal == nil {
		promoteRandomHandCardToArsenal(&best, hand, arsenalCardIn)
	}
	return best
}

// defenseSumFromRoles totals Defense() across every Defend-role card per the rolesBuf
// assignment. The arsenal-in slot's bonus (ArsenalDefenseBonus) is added when it took the
// Defend role — matching fillPartitionPerCardBufs's per-card dvals layout. Hand cards
// that opt into ArsenalDefenseBonus don't get the bonus here because they aren't in the
// arsenal slot; the bonus only applies to cards actually played from arsenal.
func defenseSumFromRoles(hand []Card, arsenalCardIn Card, rolesBuf []Role, n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		if rolesBuf[i] == Defend {
			sum += hand[i].Defense()
		}
	}
	if arsenalCardIn != nil && rolesBuf[n] == Defend {
		sum += arsenalCardIn.Defense()
		if ab, ok := arsenalCardIn.(ArsenalDefenseBonus); ok {
			sum += ab.ArsenalDefenseBonus()
		}
	}
	return sum
}

// mapCachedRolesToHand walks entry.line and the new call's hand, assigning each hand /
// arsenal-in card a role from the cached entry by ID. Returns false on multiset mismatch
// — a should-never-happen condition because the cache key locks the multiset down;
// replayBest panics on a false return so any cache invariant violation is loud.
//
// The arsenal-in card (if present) maps to the cached entry whose FromArsenal is true.
// Hand cards consume the remaining ID-matched roles in order. postPromotedFromHeld is set
// to the hand index of a card whose cached role is Arsenal but FromArsenal=false (the
// post-hoc promotion target); -1 if no such card exists.
func mapCachedRolesToHand(cachedLine []CardAssignment, hand []Card, arsenalCardIn Card, rolesBuf []Role, postPromotedFromHeld *int) bool {
	*postPromotedFromHeld = -1
	// First pass: pick out the FromArsenal=true entry (if any) and reserve it for the
	// arsenal-in card. The rest stay available for hand-card matching.
	used := make([]bool, len(cachedLine))
	if arsenalCardIn != nil {
		matched := false
		for i, a := range cachedLine {
			if a.FromArsenal && a.Card.ID() == arsenalCardIn.ID() {
				rolesBuf[len(hand)] = a.Role
				used[i] = true
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	// Second pass: walk hand in order, assigning each card the first unused cached entry
	// matching its ID. Track post-hoc promotion: if any unused entry has Role=Arsenal and
	// FromArsenal=false, the matched hand card is the promoted-from-Held card; record its
	// hand index so the chain-run treats it as Held (and the post-promotion step re-flips
	// it). Multiple Held-then-Arsenal candidates aren't possible (post-hoc promotes one
	// card at most).
	for hi, c := range hand {
		matched := false
		for i, a := range cachedLine {
			if used[i] || a.Card.ID() != c.ID() || a.FromArsenal {
				continue
			}
			rolesBuf[hi] = a.Role
			used[i] = true
			matched = true
			if a.Role == Arsenal {
				*postPromotedFromHeld = hi
			}
			break
		}
		if !matched {
			return false
		}
	}
	// Sanity: every cached entry should be claimed.
	for _, u := range used {
		if !u {
			return false
		}
	}
	return true
}

