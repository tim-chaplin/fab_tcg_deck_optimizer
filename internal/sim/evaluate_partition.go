package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// evaluatePartition is the shared "given a fixed role assignment, score it" body used
// by both findBest's recurse leaf (one of many partitions explored during the search)
// and replayBest (the cached partition replayed without searching). It groups hand
// cards into pitched/attackers/defenders/held, folds the arsenal-in card into the right
// bucket per rolesBuf[n], computes the arsenal indices, and forwards everything to
// bestAttackWithWeapons. The winning *GameState for this partition's chain is returned
// for the caller to keep — discarding losers requires no explicit work since the only
// reference to them was inside this leaf.
//
// rolesBuf must be in sync with hand (rolesBuf[i] is the role assigned to hand[i]) and
// must use only Pitch/Attack/Defend/Held in hand-card slots — hand cards never have
// Arsenal role during the chain run. The cache-replay caller, whose stored BestLine may
// contain a post-hoc-promoted hand entry tagged Arsenal, flips that entry back to Held
// before calling and restores the Arsenal tag on the returned BestLine afterward.
//
// Mutates the bufs scratch slices (pitchedBuf, attackersBuf, defendersBuf, heldBuf) in
// place; both callers feed pooled scratch through bufs and tolerate the rewrite.
func (e *Evaluator) evaluatePartition(
	masterState *gameengine.GameState,
	weapons []Weapon, hand []card.Card,
	d *deck.Deck,
	rolesBuf []deck.Role, n int, bufs *attackBufs,
	mp Matchup, defenseSum int,
	skipLog bool,
) (
	attackDealt, defenseDealt int,
	swung []string, winner *gameengine.GameState,
	ok, cacheable bool,
	arsenalAtChainStart card.Card,
) {
	arsenalCardIn := masterState.Arsenal()
	// Group hand cards into played / pitched / defending buckets, then fold in the
	// arsenal-in card based on its slot's role.
	p, a, defs := groupByRoleInto(
		hand, rolesBuf[:n],
		bufs.pitchedBuf[:0], bufs.attackersBuf[:0], bufs.defendersBuf[:0],
	)
	if arsenalCardIn != nil {
		switch rolesBuf[n] {
		case deck.Attack:
			a = append(a, arsenalCardIn)
		case deck.Defend:
			defs = append(defs, arsenalCardIn)
		}
	}
	arsenalInIdx := -1
	if arsenalCardIn != nil && rolesBuf[n] == deck.Attack {
		arsenalInIdx = len(a) - 1
	}
	arsenalDefenderIdx := -1
	if arsenalCardIn != nil && rolesBuf[n] == deck.Defend {
		arsenalDefenderIdx = len(defs) - 1
	}
	h := gatherHeldCards(hand, rolesBuf[:n], bufs.heldBuf[:0])
	arsenalAtChainStart = findArsenalCard(rolesBuf, arsenalCardIn, n)

	attackDealt, defenseDealt, _, swung, winner, ok, cacheable = bestAttackWithWeapons(
		masterState, weapons, a, defs, p, h, d, bufs,
		mp, defenseSum,
		arsenalInIdx, arsenalDefenderIdx, arsenalAtChainStart,
		skipLog,
	)
	return
}
