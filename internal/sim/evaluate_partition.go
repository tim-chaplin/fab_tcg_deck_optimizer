package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// evaluatePartition scores a fixed role assignment: it groups the hand-card slots of
// pcards into pitched/attackers/defenders/held, folds the arsenal-in slot (pcards[n], when
// present) into its bucket per its role, computes the arsenal indices, and forwards to
// bestAttackWithWeapons. The winning *GameState for this partition's chain is returned for
// the caller to keep.
//
// pcards[:n] are the hand-card slots and carry only Pitch/Attack/Defend/Held roles — hand
// cards never have Arsenal role during the chain run.
//
// Mutates the bufs scratch slices (pitchedBuf, attackersBuf, defendersBuf, heldBuf) in
// place.
func (e *Evaluator) evaluatePartition(
	masterState *gameengine.GameState,
	weapons []weapon.Weapon,
	d *deck.Deck,
	pcards []partitionCard, n int, bufs *attackBufs,
	defenseSum int,
) (
	attackDealt, defenseDealt int,
	swung []string, winner *gameengine.GameState,
	ok, cacheable bool,
	arsenalAtChainStart card.Card,
) {
	p, a, defs, h, arsenalInIdx, arsenalDefenderIdx, arsenalAtChainStart := groupPartition(pcards, n, bufs)
	attackDealt, defenseDealt, _, swung, winner, ok, cacheable = bestAttackWithWeapons(
		masterState, weapons, a, defs, p, h, d, bufs,
		defenseSum,
		arsenalInIdx, arsenalDefenderIdx, arsenalAtChainStart,
	)
	return
}

// groupPartition splits a role-assigned partition into its played / defending / pitched /
// held card buckets and resolves the arsenal-in indices, reusing the bufs scratch slices.
// The buckets fold in the arsenal-in card (pcards[n], when present) per its slot's role.
func groupPartition(pcards []partitionCard, n int, bufs *attackBufs) (p, a, defs, h []card.Card, arsenalInIdx, arsenalDefenderIdx int, arsenalAtChainStart card.Card) {
	p, a, defs = groupByRole(
		pcards[:n],
		bufs.pitchedBuf[:0], bufs.attackersBuf[:0], bufs.defendersBuf[:0],
	)
	arsenalInIdx = -1
	arsenalDefenderIdx = -1
	if len(pcards) > n {
		switch arsenalPC := pcards[n]; arsenalPC.role {
		case card.Attack:
			a = append(a, arsenalPC.card)
			arsenalInIdx = len(a) - 1
		case card.Defend:
			defs = append(defs, arsenalPC.card)
			arsenalDefenderIdx = len(defs) - 1
		}
	}
	h = gatherHeldCards(pcards[:n], bufs.heldBuf[:0])
	arsenalAtChainStart = findArsenalCard(pcards, n)
	return
}
