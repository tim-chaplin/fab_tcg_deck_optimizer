package sim

// Cache-hit replay: rebuild a TurnSummary by replaying the cached winning solution
// verbatim — the captured attacker order, modal modes, pitch ordering, and blocker modes
// run as a single chain, with no partition search and no permutation enumeration. The
// post-hoc arsenal promotion + Hand carryover bookkeeping still run, so the resulting
// summary is byte-identical to a from-scratch Best call.

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// replayBest is the cache-hit body. Builds the partitionCard set for the new call's hand,
// maps the cached BestLine's roles onto it, hands off to replaySolution for the verbatim
// chain run, then assembles the TurnSummary.
//
// Quirk: the cached entry may tag a hand card with Role=Arsenal (post-hoc promotion
// target), but hand cards never have that role during the chain run. Flip the slot back
// to Held before replaySolution; re-stamp Arsenal on the BestLine afterward.
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
	// maxResourceBonus is left unset: it feeds only the attack-budget prune, which the
	// verbatim replay path never runs.
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

	attackDealt, defenseDealt, winner, arsenalAtChainStart := e.replaySolution(masterState, weapons, d, entry, pcards, n, bufs)
	if winner == nil {
		panic(fmt.Sprintf("replayBest: cached solution is infeasible — cache invariant violated (hand=%d, incoming=%d)",
			len(hand), masterState.IncomingDamage()))
	}

	if postPromotedFromHeld >= 0 {
		pcards[postPromotedFromHeld].role = card.Arsenal
	}

	winner.SetArsenal(arsenalAtChainStart)
	best := TurnSummary{
		BestLine:       make([]card.CardAssignment, totalN),
		Value:          attackDealt + defenseDealt,
		SwungWeapons:   append([]string(nil), entry.swungWeapons...),
		IncomingDamage: masterState.IncomingDamage(),
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

// replaySolution replays a cached winning solution verbatim — no search, no prune. It
// groups the role-assigned partition, runs the defense phase with the cached blocker
// modes, and runs the single winning attack sequence with the cached order + modes. The
// chain output is byte-identical to the from-scratch Best call that produced the entry.
func (e *Evaluator) replaySolution(
	masterState *gameengine.GameState, weapons []weapon.Weapon, d *deck.Deck,
	entry evalCacheEntry, pcards []partitionCard, n int, bufs *attackBufs,
) (attackDealt, defenseDealt int, winner *gameengine.GameState, arsenalAtChainStart card.Card) {
	p, a, defs, h, arsenalInIdx, arsenalDefenderIdx, arsenalAtChainStart := groupPartition(pcards, n, bufs)
	ctx := newSequenceContext(masterState, weapons, a, defs, p, h, d, bufs, defenseSumFromRoles(pcards), arsenalInIdx, arsenalAtChainStart)

	ctx.attackPitchPerm = entry.pitchOrder
	ctx.attackPitchVals = bufs.pitchPermValsBuf[:0]
	for _, c := range entry.pitchOrder {
		ctx.attackPitchVals = append(ctx.attackPitchVals, c.Pitch())
	}
	ctx.resourceBudget = 0

	incoming := masterState.IncomingDamage()
	if len(defs) > 0 {
		defenseDealt, _, ctx.handStart = ctx.runDefense(defs, p, h, d, incoming, noBlockBudgetCap, arsenalDefenderIdx, entry.defenders)
	} else if incoming > 0 {
		ctx.leafState.SetIsMyTurn(false)
		ctx.leafState.Engine().FireTriggers(triggertype.DamageTaken, nil)
	}
	ctx.leafState.SetDeck(nil)
	ctx.seedPoolGravBuf(len(entry.attackOrder), len(ctx.attackPitchPerm))

	attackDealt, _, _, _ = ctx.playSequenceModal(entry.attackOrder)
	return attackDealt, defenseDealt, ctx.permState, arsenalAtChainStart
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
