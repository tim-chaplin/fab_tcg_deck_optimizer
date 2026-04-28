package sim

// Cache-hit replay: rebuild a TurnSummary from a cached winning partition by running the
// chain dispatcher against just that one role assignment. Skips the partition search
// (the dominant cost — exponential in hand size) but still runs bestAttackWithWeapons
// once and the post-hoc arsenal promotion + Hand carryover bookkeeping that findBest
// does after the search loop, so the resulting summary is byte-identical to a full
// from-scratch Best call.

import "fmt"

// replayBest is the cache-hit body. Parallels findBest's miss path but runs only one
// partition: the one whose roles entry caches. Caller has already verified the key is
// valid (priorAuraTriggers empty, hand size in bounds, key matches an existing entry).
//
// Steps:
//  1. Walk the cached BestLine and project its (cardID, role, fromArsenal) tuples back
//     onto the new call's hand + arsenal-in. Produces the same a / d / p / h slices the
//     original search's winning leaf passed to bestAttackWithWeapons.
//  2. Run bestAttackWithWeapons once with that one partition. Returns the same Value /
//     CarryState the original search produced (Best is deterministic given inputs).
//  3. Build the TurnSummary: copy BestLine roles, attach SwungWeapons, adopt CarryState,
//     re-do the post-hoc arsenal promotion when needed.
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

	// Project the cached role tuples back onto this call's slice positions. The cached
	// BestLine and the new hand share the same multiset (we keyed on it), so we walk the
	// new hand in order and assign each card the next-available cached role for its ID.
	// arsenalCardIn (if present) maps to whichever cached entry has FromArsenal=true.
	rolesBuf := bufs.rolesBuf[:totalN]
	postPromotedFromHeld := -1 // index in `hand` of the post-hoc-promoted Held card, -1 if no promotion happened
	if !mapCachedRolesToHand(entry.line, hand, arsenalCardIn, rolesBuf, &postPromotedFromHeld) {
		// Multiset mismatch can't happen by construction — the cache key sorts hand IDs
		// and the entry was stored under that exact key, so a hit means the multisets are
		// identical. Reaching here indicates a bug somewhere (key collision, cache
		// corruption, mid-call mutation of cachedLine, etc.) that's already compromised
		// correctness; panic loudly so the operator notices instead of falling back to a
		// silent re-search that hides the cache bug.
		panic(fmt.Sprintf("replayBest: mapCachedRolesToHand failed despite cache hit — cache invariant violated (hand=%d, cachedLine=%d, arsenal=%v)",
			len(hand), len(entry.line), arsenalCardIn != nil))
	}

	// The cached partition tells us which cards Pitched / Attacked / Defended / Held /
	// went to Arsenal. The post-hoc promotion already happened in the cached entry: a
	// Held card may show role=Arsenal; we treat it as Held during the chain run (so it
	// threads through state.Hand / the held buf), then re-do the promotion below to make
	// sure best.State.Arsenal lands the right card and BestLine.Role flips back.
	pitched := bufs.pitchedBuf[:0]
	attackers := bufs.attackersBuf[:0]
	defenders := bufs.defendersBuf[:0]
	held := bufs.heldBuf[:0]
	for i, c := range hand {
		role := rolesBuf[i]
		if i == postPromotedFromHeld {
			role = Held // chain-run treats it as Held; promotion re-runs after
		}
		switch role {
		case Pitch:
			pitched = append(pitched, c)
		case Attack:
			attackers = append(attackers, c)
		case Defend:
			defenders = append(defenders, c)
		case Held:
			held = append(held, c)
		}
	}
	arsenalInIdx := -1
	arsenalDefenderIdx := -1
	var arsenalAtChainStart Card
	if arsenalCardIn != nil {
		switch rolesBuf[n] {
		case Attack:
			attackers = append(attackers, arsenalCardIn)
			arsenalInIdx = len(attackers) - 1
		case Defend:
			defenders = append(defenders, arsenalCardIn)
			arsenalDefenderIdx = len(defenders) - 1
		case Arsenal:
			arsenalAtChainStart = arsenalCardIn
		}
	}

	// defenseSum has to match what the original search computed — sum of Defense() across
	// every Defend-role card (DR or plain), per fillPartitionPerCardBufs. It feeds
	// state.BlockTotal so DR Plays that read "did we block all incoming?" see the right
	// shape. defendersDamage's per-DR seed reads it from state.IncomingDamage, not
	// BlockTotal, so this only matters for cards that consult state.BlockTotal directly.
	defenseSum := 0
	for _, d := range defenders {
		defenseSum += d.Defense()
		// Arsenal-defender DRs that opt into +N{d} from arsenal use ArsenalDefenseBonus;
		// match fillPartitionPerCardBufs's behavior so BlockTotal includes the rider.
	}
	if arsenalDefenderIdx >= 0 {
		if ab, ok := defenders[arsenalDefenderIdx].(ArsenalDefenseBonus); ok {
			defenseSum += ab.ArsenalDefenseBonus()
		}
	}

	attackDealt, defenseDealt, _, _, swung, carry, ok, _ := bestAttackWithWeapons(
		hero, weapons, attackers, defenders, pitched, held, deck, bufs,
		runechantCarryover, incomingDamage, defenseSum, arsenalInIdx, arsenalDefenderIdx,
		arsenalAtChainStart, nil, skipLog,
	)
	if !ok {
		// Infeasible-partition replay can't happen by construction — the cached entry was
		// only stored after best.Cacheable=true, which means a feasible partition was
		// found AND no leaf read deck/graveyard. Reaching here means either the cache
		// stored an infeasible result (bug), or the inputs that should be deterministic
		// (resourceBudget, costs) somehow shifted. Panic so the operator notices rather
		// than silently re-searching and hiding a real correctness bug.
		panic(fmt.Sprintf("replayBest: cached partition is infeasible — cache invariant violated (hand=%d, runechantCarryover=%d, incomingDamage=%d)",
			len(hand), runechantCarryover, incomingDamage))
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

