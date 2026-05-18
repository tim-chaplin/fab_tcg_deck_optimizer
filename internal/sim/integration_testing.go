package sim

// Integration-test entry points that drive one or two turns in source order (no shuffle).
// Tests assert on summary.Value (chain + start-of-turn tick damage) and summary.State
// (post end-of-turn cleanup + next-hand draw).

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// EvalOneTurnForTesting drives one turn and returns its TurnSummary. State reflects the
// end-of-turn boundary (pitched recycled, next hand drawn; partial or empty hands OK).
// initial seeds carryover (hero, arsenal, auras, items, banished, graveyard, opponentMarked,
// incoming damage); nil uses gameengine's default state (20hp/4int, no abilities).
// initialHand is the dealt hand — required (callers that want a deck-drawn hand do the
// d.Draw themselves and pass the result). d is consumed in place — callers shouldn't reuse
// it across calls.
func EvalOneTurnForTesting(d *deck.Deck, initial *gameengine.GameState, initialHand []card.Card) TurnSummary {
	master, weapons := setupTurn(d, initial)
	summary, _ := playOneTurn(master, initialHand, d, weapons, ev(), nil, nil)
	return summary
}

// EvalTwoTurnsForTesting drives two turns, threading carryover between them. Each turn's
// Value folds in its own start-of-turn trigger damage. Turn 2 runs even with an empty or
// partial hand. The returned turn1.State is an independent snapshot; turn2 mutates the
// shared master.
func EvalTwoTurnsForTesting(d *deck.Deck, initial *gameengine.GameState, hand1 []card.Card) (TurnSummary, TurnSummary) {
	master, weapons := setupTurn(d, initial)

	turn1, _ := playOneTurn(master, hand1, d, weapons, ev(), nil, nil)
	stable := turn1
	stable.State = snapshotState(turn1.State, turn1.State.Deck(), turn1.State.Hand(), turn1.Value, turn1.State.CardsDrawn())

	turn2, _ := playOneTurn(turn1.State, turn1.State.Hand(), turn1.State.Deck(), weapons, ev(), nil, nil)
	return stable, turn2
}

// setupTurn assembles the per-turn fixture from d: master *GameState (`initial` when
// provided, else the builder's default state) and the weapon list.
func setupTurn(d *deck.Deck, initial *gameengine.GameState) (master *gameengine.GameState, weapons []weapon.Weapon) {
	if initial != nil {
		master = initial
	} else {
		master = gameengine.GameStateBuilder().Build()
	}
	weapons = make([]weapon.Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}
	return master, weapons
}

// ev returns the package-level Evaluator so test helpers share its cache/scratch state.
func ev() *Evaluator { return sharedEvaluator }

// snapshotState clones src's persistent carryover and overlays deck (copied), hand, value,
// and cardsDrawn. The result is independent of future mutations to src.
func snapshotState(src *gameengine.GameState, d *deck.Deck, hand []card.Card, value, cardsDrawn int) *gameengine.GameState {
	out := src.CopyPersistentState()
	if d != nil {
		out.SetDeck(d.Copy())
	}
	out.SetHand(hand)
	out.SetValue(value)
	out.SetCardsDrawn(cardsDrawn)
	return out
}
