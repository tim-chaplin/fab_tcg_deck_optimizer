package sim

// Cache-hit replay: rebuild a TurnSummary by replaying the cached winning solution
// verbatim — captured attacker order, modal modes, pitch ordering, blocker modes run as
// a single attack turn, no partition search, no permutation enumeration. Post-hoc arsenal
// promotion + Hand carryover still run, so the summary is byte-identical to a from-scratch
// Best call.

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// replayBest is the cache-hit body. Builds the partitionCard set, maps the cached BestLine's
// roles onto it, runs replaySolution for the verbatim attack turn, then assembles the TurnSummary.
//
// Quirk: a cached entry may tag a hand card with Role=Arsenal (post-hoc promotion target),
// but hand cards never carry that role during the attack turn run. Flip it to Held before
// replaySolution; re-stamp Arsenal on the BestLine afterward.
func (e *Evaluator) replayBest(
	entry evalCacheEntry,
	weapons []weapon.Weapon, hand []card.Card,
	d *deck.Deck,
	masterState *gameengine.GameState,
) TurnSummary {
	arsenalCardIn := masterState.Arsenal()
	n := len(hand)
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	bufs := e.getAttackBufs(n, weapons)
	// maxResourceBonus is unset: it feeds only the attack-budget prune, which replay skips.
	pcards := bufs.partitionCards[:totalN]
	fillPartitionCards(hand, n, totalN, arsenalCardIn, pcards)
	postPromotedFromHeld := -1
	if !mapCachedRolesToHand(entry.line, pcards, n, &postPromotedFromHeld) {
		panic(fmt.Sprintf("replayBest: mapCachedRolesToHand failed despite cache hit — cache invariant violated (hand=%d, cachedLine=%d, arsenal=%v)",
			len(hand), len(entry.line), arsenalCardIn != nil))
	}

	if postPromotedFromHeld >= 0 {
		pcards[postPromotedFromHeld].role = card.Held
	}

	attackDealt, defenseDealt, winner, arsenalAtAttackTurnStart := e.replaySolution(masterState, weapons, d, entry, pcards, n, bufs)
	if winner == nil {
		panic(fmt.Sprintf("replayBest: cached solution is infeasible — cache invariant violated (hand=%d, incoming=%d)",
			len(hand), masterState.IncomingPhysicalDamage()))
	}

	if postPromotedFromHeld >= 0 {
		pcards[postPromotedFromHeld].role = card.Arsenal
	}

	// Mirror partition.go's mid-turn-fill-aware restore: only re-seed the arsenal when
	// the per-perm didn't fill it mid-turn (Promise of Plenty's on-hit). A nil arsenal at
	// end of replay means no fill landed; restore via arsenalAtAttackTurnStart.
	if winner.Arsenal() == nil {
		winner.SetArsenal(arsenalAtAttackTurnStart)
	}
	best := TurnSummary{
		BestLine:       make([]card.CardAssignment, totalN),
		Value:          attackDealt + defenseDealt,
		SwungWeapons:   append([]string(nil), entry.swungWeapons...),
		IncomingPhysicalDamage: masterState.IncomingPhysicalDamage(),
		Cacheable:      true,
		State:          winner,
	}
	for i := 0; i < n; i++ {
		best.BestLine[i] = card.CardAssignment{Card: hand[i], Role: pcards[i].role}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = card.CardAssignment{Card: arsenalCardIn, Role: pcards[n].role, FromArsenal: true}
	}
	if best.State.Arsenal() == nil {
		promoteHeldToArsenal(best.State, hand, arsenalCardIn)
	}
	return best
}

// replaySolution replays a cached winning solution verbatim — no search, no prune. Groups
// the role-assigned partition, runs defense with the cached blocker modes, then runs the
// single winning attack sequence. Output is byte-identical to the originating Best call.
func (e *Evaluator) replaySolution(
	masterState *gameengine.GameState, weapons []weapon.Weapon, d *deck.Deck,
	entry evalCacheEntry, pcards []partitionCard, n int, bufs *attackBufs,
) (attackDealt, defenseDealt int, winner *gameengine.GameState, arsenalAtAttackTurnStart card.Card) {
	p, a, defs, h, arsenalInIdx, arsenalDefenderIdx, arsenalAtAttackTurnStart := groupPartition(pcards, n, bufs)
	ctx := newSequenceContext(masterState, weapons, a, defs, p, h, d, bufs, defenseSumFromRoles(pcards), arsenalInIdx, arsenalAtAttackTurnStart)
	defer ctx.releaseLeafState()

	ctx.attackPitchPerm = wrapPitchedCards(entry.pitchOrder)
	ctx.attackPitchVals = bufs.pitchPermValsBuf[:0]
	for _, c := range entry.pitchOrder {
		ctx.attackPitchVals = append(ctx.attackPitchVals, c.Pitch())
	}
	ctx.resourceBudget = 0

	incoming := masterState.IncomingPhysicalDamage()
	arcaneIncoming := masterState.IncomingArcaneDamage()
	if len(defs) > 0 {
		installLeafDeck(ctx, bufs, d)
		defenseDealt, _, ctx.handStart = ctx.runDefense(defs, p, h, ctx.deck, incoming, noBlockBudgetCap, arsenalDefenderIdx, entry.defenders)
	} else if incoming > 0 || arcaneIncoming > 0 {
		defenseDealt = ctx.fireUndefendedDamageTriggers()
	}
	attackDealt, _, _, _ = ctx.playSequenceModal(entry.attackOrder)
	return attackDealt, defenseDealt, ctx.permState, arsenalAtAttackTurnStart
}

// defenseSumFromRoles totals the defenseVal of every Defend-role card.
func defenseSumFromRoles(pcards []partitionCard) int {
	sum := 0
	for _, pc := range pcards {
		if pc.role == card.Defend {
			sum += pc.defenseVal
		}
	}
	return sum
}

// mapCachedRolesToHand walks entry.line and assigns each partitionCard a role from the
// cached entry by matching Card.ID(). Returns false on multiset mismatch.
func mapCachedRolesToHand(cachedLine []card.CardAssignment, pcards []partitionCard, n int, postPromotedFromHeld *int) bool {
	*postPromotedFromHeld = -1
	used := make([]bool, len(cachedLine))
	if len(pcards) > n {
		matched := false
		for i, a := range cachedLine {
			if a.FromArsenal && a.Card.ID() == pcards[n].card.ID() {
				pcards[n].role = a.Role
				used[i] = true
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	for hi := 0; hi < n; hi++ {
		matched := false
		for i, a := range cachedLine {
			if used[i] || a.Card.ID() != pcards[hi].card.ID() || a.FromArsenal {
				continue
			}
			pcards[hi].role = a.Role
			used[i] = true
			matched = true
			if a.Role == card.Arsenal {
				*postPromotedFromHeld = hi
			}
			break
		}
		if !matched {
			return false
		}
	}
	for _, u := range used {
		if !u {
			return false
		}
	}
	return true
}
