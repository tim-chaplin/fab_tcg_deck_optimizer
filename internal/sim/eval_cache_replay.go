package sim

// Cache-hit replay: rebuild a TurnSummary from a cached winning partition by running the
// chain dispatcher against just that one role assignment. Skips the partition search
// (the dominant cost — exponential in hand size) but still runs bestAttackWithWeapons
// once and the post-hoc arsenal promotion + Hand carryover bookkeeping that findBest
// does after the search loop, so the resulting summary is byte-identical to a full
// from-scratch Best call.

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// replayBest is the cache-hit body. Thin wrapper around evaluatePartition: project the
// cached BestLine onto the new call's hand to fill rolesBuf, hand off to evaluatePartition
// for the actual chain run, then assemble the TurnSummary from its outputs. The cache key
// already locked the inputs (hand multiset, runechantCarryover, arsenalCardIn, auras) so
// the chain output here is byte-identical to what the original cached call produced.
//
// One quirk in projecting the BestLine: the cached entry may tag a hand card with
// Role=Arsenal (the post-hoc promotion target). Hand cards never have that role during
// the chain run — the search treats them as Held and the post-hoc step re-flips at the
// end — so we flip the entry back to Held before evaluatePartition and re-stamp Arsenal
// on the BestLine afterward.
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
	rolesBuf := bufs.rolesBuf[:totalN]
	postPromotedFromHeld := -1
	if !mapCachedRolesToHand(entry.line, hand, arsenalCardIn, rolesBuf, &postPromotedFromHeld) {
		panic(fmt.Sprintf("replayBest: mapCachedRolesToHand failed despite cache hit — cache invariant violated (hand=%d, cachedLine=%d, arsenal=%v)",
			len(hand), len(entry.line), arsenalCardIn != nil))
	}

	if postPromotedFromHeld >= 0 {
		rolesBuf[postPromotedFromHeld] = deck.Held
	}

	defenseSum := defenseSumFromRoles(hand, arsenalCardIn, rolesBuf, n)

	attackDealt, defenseDealt, swung, winner, ok, _, arsenalAtChainStart := e.evaluatePartition(
		masterState, weapons, hand, d,
		rolesBuf, n, bufs,
		defenseSum,
	)
	if !ok {
		panic(fmt.Sprintf("replayBest: cached partition is infeasible — cache invariant violated (hand=%d, incoming=%d)",
			len(hand), masterState.IncomingDamage()))
	}

	if postPromotedFromHeld >= 0 {
		rolesBuf[postPromotedFromHeld] = deck.Arsenal
	}

	winner.SetArsenal(arsenalAtChainStart)
	best := TurnSummary{
		BestLine:       make([]deck.CardAssignment, totalN),
		Value:          attackDealt + defenseDealt,
		SwungWeapons:   append([]string(nil), swung...),
		IncomingDamage: masterState.IncomingDamage(),
		Cacheable:      true,
		State:          winner,
	}
	for i := 0; i < n; i++ {
		best.BestLine[i] = deck.CardAssignment{Card: hand[i], Role: rolesBuf[i]}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = deck.CardAssignment{Card: arsenalCardIn, Role: rolesBuf[n], FromArsenal: true}
	}
	if best.State.Arsenal() == nil {
		promoteRandomHandCardToArsenal(&best, hand, arsenalCardIn)
	}
	return best
}

// defenseSumFromRoles totals Defense() across every Defend-role card per the rolesBuf
// assignment.
func defenseSumFromRoles(hand []card.Card, arsenalCardIn card.Card, rolesBuf []deck.Role, n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		if rolesBuf[i] == deck.Defend {
			sum += hand[i].Defense()
		}
	}
	if arsenalCardIn != nil && rolesBuf[n] == deck.Defend {
		sum += arsenalCardIn.Defense() + card.ArsenalDefenseBonusOf(arsenalCardIn)
	}
	return sum
}

// mapCachedRolesToHand walks entry.line and the new call's hand, assigning each hand /
// arsenal-in card a role from the cached entry by ID. Returns false on multiset mismatch.
func mapCachedRolesToHand(cachedLine []deck.CardAssignment, hand []card.Card, arsenalCardIn card.Card, rolesBuf []deck.Role, postPromotedFromHeld *int) bool {
	*postPromotedFromHeld = -1
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
	for hi, c := range hand {
		matched := false
		for i, a := range cachedLine {
			if used[i] || a.Card.ID() != c.ID() || a.FromArsenal {
				continue
			}
			rolesBuf[hi] = a.Role
			used[i] = true
			matched = true
			if a.Role == deck.Arsenal {
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
