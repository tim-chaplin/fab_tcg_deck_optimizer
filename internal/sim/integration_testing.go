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
// post end-of-turn cleanup boundary (pitched recycled, arsenal refilled, next hand drawn —
// partial or empty hands are fine). initial seeds carryover (Hero, Arsenal, Auras, Items,
// Banished, Graveyard, OpponentMarked); nil starts clean inheriting the deck's hero.
func EvalOneTurnForTesting(masterDeck *deck.Deck, mp Matchup, initial *gameengine.GameState, initialHand []deck.Card) TurnSummary {
	d, master, hand, handSize, weapons, ok := setupTurn(masterDeck, mp, initial, initialHand)
	if !ok {
		return TurnSummary{State: gameengine.GameStateBuilder().SetHero(d.Hero.(hero.Hero)).Build()}
	}
	summary, _ := playOneTurn(master, hand, d, mp, weapons, ev())
	turnDraws := summary.State.CardsDrawn()
	d, master, nextHand := advanceAndDraw(summary, d, handSize)
	summary.State = snapshotState(master, d, nextHand, summary.Value, turnDraws)
	return summary
}

// EvalTwoTurnsForTesting drives two turns, threading carryover between them. Each turn's
// Value folds in its own start-of-turn trigger damage (turn2 sees auras created by turn1).
// Turn 2 runs even with an empty or partial hand (matching the real-game rule that turns
// continue past deck exhaustion).
func EvalTwoTurnsForTesting(masterDeck *deck.Deck, mp Matchup, initial *gameengine.GameState, hand1 []deck.Card) (TurnSummary, TurnSummary) {
	d, master, hand, handSize, weapons, ok := setupTurn(masterDeck, mp, initial, hand1)
	if !ok {
		return TurnSummary{State: gameengine.GameStateBuilder().SetHero(d.Hero.(hero.Hero)).Build()}, TurnSummary{}
	}

	turn1, _ := playOneTurn(master, hand, d, mp, weapons, ev())
	turn1Draws := turn1.State.CardsDrawn()
	d, master, turn2Hand := advanceAndDraw(turn1, d, handSize)
	turn1.State = snapshotState(master, d, turn2Hand, turn1.Value, turn1Draws)

	turn2, _ := playOneTurn(master, turn2Hand, d, mp, weapons, ev())
	turn2Draws := turn2.State.CardsDrawn()
	d, master, turn3Hand := advanceAndDraw(turn2, d, handSize)
	turn2.State = snapshotState(master, d, turn3Hand, turn2.Value, turn2Draws)

	return turn1, turn2
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

// advanceAndDraw runs end-of-turn cleanup (recycle pitched, reset ephemeral, capture held)
// and draws the next hand (partial draws OK). Returns the post-recycle deck, mutated
// master, and next hand (nil when held is empty and no draws were possible).
func advanceAndDraw(summary TurnSummary, d *deck.Deck, handSize int) (*deck.Deck, *gameengine.GameState, []card.Card) {
	carry, _ := advanceToNextTurn(summary, nil, nil)
	d = carry.deck
	master := carry.master
	held := carry.held

	toDraw := handSize - len(held)
	if toDraw > d.Size() {
		toDraw = d.Size()
	}
	if toDraw <= 0 {
		if len(held) == 0 {
			return d, master, nil
		}
		return d, master, held
	}
	nextHand := make([]card.Card, 0, handSize)
	nextHand = append(nextHand, held...)
	for _, c := range d.Draw(toDraw) {
		nextHand = append(nextHand, c.(card.Card))
	}
	return d, master, nextHand
}

// ev returns the package-level Evaluator so test helpers share its cache/scratch state.
func ev() *Evaluator { return sharedEvaluator }

// snapshotState clones src's persistent carryover and overlays deck (copied), hand, value,
// and cardsDrawn.
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
