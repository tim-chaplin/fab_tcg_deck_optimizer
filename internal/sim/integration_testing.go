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
// initial seeds carryover; nil starts clean inheriting the deck's hero.
func EvalOneTurnForTesting(masterDeck *deck.Deck, mp Matchup, initial *gameengine.GameState, initialHand []deck.Card) TurnSummary {
	d, master, hand, handSize, weapons, ok := setupTurn(masterDeck, mp, initial, initialHand)
	if !ok {
		return TurnSummary{State: gameengine.GameStateBuilder().SetHero(d.Hero.(hero.Hero)).Build()}
	}
	summary, _ := playOneTurn(master, hand, d, mp, weapons, ev(), handSize, nil, nil, nil)
	return summary
}

// EvalTwoTurnsForTesting drives two turns, threading carryover between them. Each turn's
// Value folds in its own start-of-turn trigger damage. Turn 2 runs even with an empty or
// partial hand. The returned turn1.State is an independent snapshot; turn2 mutates the
// shared master.
func EvalTwoTurnsForTesting(masterDeck *deck.Deck, mp Matchup, initial *gameengine.GameState, hand1 []deck.Card) (TurnSummary, TurnSummary) {
	d, master, hand, handSize, weapons, ok := setupTurn(masterDeck, mp, initial, hand1)
	if !ok {
		return TurnSummary{State: gameengine.GameStateBuilder().SetHero(d.Hero.(hero.Hero)).Build()}, TurnSummary{}
	}

	turn1, _ := playOneTurn(master, hand, d, mp, weapons, ev(), handSize, nil, nil, nil)
	stable := turn1
	stable.State = snapshotState(turn1.State, turn1.State.Deck(), turn1.State.Hand(), turn1.Value, turn1.State.CardsDrawn())

	turn2, _ := playOneTurn(turn1.State, turn1.State.Hand(), turn1.State.Deck(), mp, weapons, ev(), handSize, nil, nil, nil)
	return stable, turn2
}

// setupTurn assembles the per-turn fixture: deck copy, master *GameState seeded with
// carryover, dealt hand, handSize, weapon list. ok=false on invalid inputs (handSize <= 0;
// initialHand empty or oversized; no initialHand and deck can't deal handSize).
func setupTurn(masterDeck *deck.Deck, mp Matchup, initial *gameengine.GameState, initialHand []deck.Card) (*deck.Deck, *gameengine.GameState, []card.Card, int, []weapon.Weapon, bool) {
	d := masterDeck.Copy()
	h := d.Hero.(hero.Hero)
	var master *gameengine.GameState
	if initial != nil {
		master = initial
		if hv := master.Hero(); hv != nil {
			h = hv.(hero.Hero)
		} else {
			master.SetHero(h)
		}
	} else {
		master = gameengine.GameStateBuilder().SetHero(h).Build()
	}
	master.SetIncomingDamage(mp.IncomingDamage)
	master.SetArcaneIncomingDamage(mp.ArcaneIncomingDamage)
	handSize := h.Intelligence()
	if handSize <= 0 {
		return d, nil, nil, 0, nil, false
	}
	if initialHand == nil && d.Size() < handSize {
		return d, nil, nil, 0, nil, false
	}
	if initialHand != nil && (len(initialHand) == 0 || len(initialHand) > handSize) {
		return d, nil, nil, 0, nil, false
	}
	hand := make([]card.Card, 0, handSize+startOfTurnRevealRoom)
	if initialHand != nil {
		for _, c := range initialHand {
			hand = append(hand, c.(card.Card))
		}
	} else {
		for _, c := range d.Draw(handSize) {
			hand = append(hand, c.(card.Card))
		}
	}
	weapons := make([]weapon.Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}
	return d, master, hand, handSize, weapons, true
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
