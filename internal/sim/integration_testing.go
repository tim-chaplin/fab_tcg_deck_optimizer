package sim

// Integration-test entry points that drive one or two turns in source order (no shuffle).
// Tests assert on summary.Value (chain + start-of-turn tick damage) and summary.State
// (post end-of-turn cleanup + next-hand draw).

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// EvalOneTurnForTesting drives one turn and returns its TurnSummary. State reflects the
// end-of-turn boundary (pitched recycled, next hand drawn; partial or empty hands OK).
// initial seeds carryover (hero, arsenal, auras, items, banished, graveyard, opponentMarked,
// incoming damage); nil uses gameengine's default state (heroes.Default — 20hp/4int, no
// abilities). d is consumed in place — callers shouldn't reuse it across calls.
func EvalOneTurnForTesting(d *deck.Deck, initial *gameengine.GameState, initialHand []card.Card) TurnSummary {
	master, hand, handSize, weapons, ok := setupTurn(d, initial, initialHand)
	if !ok {
		return TurnSummary{State: gameengine.GameStateBuilder().Build()}
	}
	summary, _ := playOneTurn(master, hand, d, weapons, ev(), handSize, nil, nil)
	return summary
}

// EvalTwoTurnsForTesting drives two turns, threading carryover between them. Each turn's
// Value folds in its own start-of-turn trigger damage. Turn 2 runs even with an empty or
// partial hand. The returned turn1.State is an independent snapshot; turn2 mutates the
// shared master.
func EvalTwoTurnsForTesting(d *deck.Deck, initial *gameengine.GameState, hand1 []card.Card) (TurnSummary, TurnSummary) {
	master, hand, handSize, weapons, ok := setupTurn(d, initial, hand1)
	if !ok {
		return TurnSummary{State: gameengine.GameStateBuilder().Build()}, TurnSummary{}
	}

	turn1, _ := playOneTurn(master, hand, d, weapons, ev(), handSize, nil, nil)
	stable := turn1
	stable.State = snapshotState(turn1.State, turn1.State.Deck(), turn1.State.Hand(), turn1.Value, turn1.State.CardsDrawn())

	turn2, _ := playOneTurn(turn1.State, turn1.State.Hand(), turn1.State.Deck(), weapons, ev(), handSize, nil, nil)
	return stable, turn2
}

// setupTurn assembles the per-turn fixture from d: master *GameState (`initial` when
// provided, else the builder's default state), dealt hand, handSize, weapon list. d is
// consumed in place (Draw mutates it when initialHand is nil). ok=false on invalid inputs
// (handSize <= 0; initialHand empty or oversized; no initialHand and deck can't deal
// handSize).
func setupTurn(d *deck.Deck, initial *gameengine.GameState, initialHand []card.Card) (master *gameengine.GameState, hand []card.Card, handSize int, weapons []weapon.Weapon, ok bool) {
	// Master GameState.
	if initial != nil {
		master = initial
	} else {
		master = gameengine.GameStateBuilder().Build()
	}

	// Hand size + input validation.
	handSize = master.Hero().(hero.Hero).Intelligence()
	if handSize <= 0 {
		return nil, nil, 0, nil, false
	}
	if initialHand == nil && d.Size() < handSize {
		return nil, nil, 0, nil, false
	}
	if initialHand != nil && (len(initialHand) == 0 || len(initialHand) > handSize) {
		return nil, nil, 0, nil, false
	}

	// Hand: caller's initialHand verbatim, or draw from deck top.
	hand = make([]card.Card, 0, handSize+startOfTurnRevealRoom)
	if initialHand != nil {
		hand = append(hand, initialHand...)
	} else {
		for _, c := range d.Draw(handSize) {
			hand = append(hand, c.(card.Card))
		}
	}

	// Weapons: widen []deck.Weapon to []weapon.Weapon.
	weapons = make([]weapon.Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}

	return master, hand, handSize, weapons, true
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
