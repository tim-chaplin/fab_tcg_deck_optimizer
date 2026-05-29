package sim

// Integration-test entry points that drive one or two turns in source order (no shuffle).
// Tests assert on summary.Value (attack turn + start-of-turn tick damage) and summary.State
// (post end-of-turn cleanup + next-hand draw).

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// EvalOneTurnForTesting drives one turn and returns its TurnSummary. State reflects the
// end-of-turn boundary (pitched recycled, next hand drawn; partial or empty hands OK).
// initialState seeds carryover (hero, arsenal, auras, items, banished, graveyard,
// opponentMarked, incoming damage); nil uses gameengine's default state (20hp/4int,
// no abilities).
// Equipped weapons come from d.Weapons. initialHand is the dealt hand — required (callers
// that want a deck-drawn hand do the d.Draw themselves and pass the result). d is consumed
// in place — callers shouldn't reuse it across calls.
func EvalOneTurnForTesting(d *deck.Deck, initialState *gameengine.GameState, initialHand []card.Card) TurnSummary {
	if initialState == nil {
		initialState = gameengine.GameStateBuilder().Build()
	}
	gameengine.EquipFromCards(initialState, weaponCards(d))
	sortHandByID(initialHand)
	initialState.SetHand(initialHand)
	summary, _ := playOneTurn(initialState, d, ev(), nil, nil)
	return summary
}

// EvalTwoTurnsForTesting drives two turns, threading carryover between them. Each turn's
// Value folds in its own start-of-turn trigger damage. Turn 2 runs even with an empty or
// partial hand. The returned turn1.State is an independent deep copy; turn2 mutates the
// shared state.
func EvalTwoTurnsForTesting(d *deck.Deck, initialState *gameengine.GameState, hand1 []card.Card) (TurnSummary, TurnSummary) {
	if initialState == nil {
		initialState = gameengine.GameStateBuilder().Build()
	}
	gameengine.EquipFromCards(initialState, weaponCards(d))
	sortHandByID(hand1)
	initialState.SetHand(hand1)

	turn1, _ := playOneTurn(initialState, d, ev(), nil, nil)
	stable := turn1
	stable.State = turn1.State.Copy()

	turn2, _ := playOneTurn(turn1.State, turn1.State.Deck(), ev(), nil, nil)
	return stable, turn2
}

// EvalFourTurnsForTesting drives four turns, threading carryover between them — the harness
// for weapons / effects whose state accumulates over several turns (Talishar's rust counters
// reach the self-destruct threshold on turn 3). Each returned summary's State is an
// independent deep copy except the last, which mutates the shared state.
func EvalFourTurnsForTesting(d *deck.Deck, initialState *gameengine.GameState, hand1 []card.Card) (TurnSummary, TurnSummary, TurnSummary, TurnSummary) {
	if initialState == nil {
		initialState = gameengine.GameStateBuilder().Build()
	}
	gameengine.EquipFromCards(initialState, weaponCards(d))
	sortHandByID(hand1)
	initialState.SetHand(hand1)

	turn1, _ := playOneTurn(initialState, d, ev(), nil, nil)
	s1 := turn1
	s1.State = turn1.State.Copy()

	turn2, _ := playOneTurn(turn1.State, turn1.State.Deck(), ev(), nil, nil)
	s2 := turn2
	s2.State = turn2.State.Copy()

	turn3, _ := playOneTurn(turn2.State, turn2.State.Deck(), ev(), nil, nil)
	s3 := turn3
	s3.State = turn3.State.Copy()

	turn4, _ := playOneTurn(turn3.State, turn3.State.Deck(), ev(), nil, nil)
	return s1, s2, s3, turn4
}

// ev returns the package-level Evaluator so test helpers share its cache/scratch state.
func ev() *Evaluator { return sharedEvaluator }
