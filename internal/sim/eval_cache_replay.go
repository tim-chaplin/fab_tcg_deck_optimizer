package sim

// Cache-hit replay: rebuild a TurnSummary from a cached winning partition by running the
// chain dispatcher against just that one role assignment. Skips the partition search (the
// dominant cost — exponential in hand size) but still runs bestAttackWithWeapons once and
// the post-hoc arsenal promotion + Hand carryover bookkeeping, so the resulting summary
// is byte-identical to a from-scratch Best call.

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// replayBest is the cache-hit body. Builds the partitionCard set for the new call's hand,
// maps the cached BestLine's roles onto it, hands off to evaluatePartition for the chain
// run, then assembles the TurnSummary. The cache key locks all inputs so the chain output
// here is byte-identical to the original call.
//
// Quirk: the cached entry may tag a hand card with Role=Arsenal (post-hoc promotion
// target), but hand cards never have that role during the chain run. Flip the slot back
// to Held before evaluatePartition; re-stamp Arsenal on the BestLine afterward.
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

	defenseSum := defenseSumFromRoles(pcards)

	attackDealt, defenseDealt, swung, winner, ok, _, arsenalAtChainStart := e.evaluatePartition(
		masterState, weapons, d,
		pcards, n, bufs,
		defenseSum,
	)
	if !ok {
		panic(fmt.Sprintf("replayBest: cached partition is infeasible — cache invariant violated (hand=%d, incoming=%d)",
			len(hand), masterState.IncomingDamage()))
	}

	if postPromotedFromHeld >= 0 {
		pcards[postPromotedFromHeld].role = card.Arsenal
	}

	winner.SetArsenal(arsenalAtChainStart)
	best := TurnSummary{
		BestLine:       make([]card.CardAssignment, totalN),
		Value:          attackDealt + defenseDealt,
		SwungWeapons:   append([]string(nil), swung...),
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
