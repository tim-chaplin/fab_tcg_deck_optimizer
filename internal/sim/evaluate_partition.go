package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// evaluatePartition scores a fixed role assignment: it groups hand cards into
// pitched/attackers/defenders/held, folds the arsenal-in card into the bucket per
// rolesBuf[n], computes the arsenal indices, and forwards to bestAttackWithWeapons. The
// winning *GameState for this partition's chain is returned for the caller to keep.
//
// rolesBuf must be in sync with hand (rolesBuf[i] is the role of hand[i]) and use only
// Pitch/Attack/Defend/Held in hand-card slots — hand cards never have Arsenal role during
// the chain run.
//
// Mutates the bufs scratch slices (pitchedBuf, attackersBuf, defendersBuf, heldBuf) in
// place.
func (e *Evaluator) evaluatePartition(
	masterState *gameengine.GameState,
	weapons []weapon.Weapon, hand []card.Card,
	d *deck.Deck,
	rolesBuf []deck.Role, n int, bufs *attackBufs,
	defenseSum int,
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
		defenseSum,
		arsenalInIdx, arsenalDefenderIdx, arsenalAtChainStart,
	)
	return
}
