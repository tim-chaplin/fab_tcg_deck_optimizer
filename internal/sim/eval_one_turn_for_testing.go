package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TurnExtras carries a turn's chain results (Value, BestLine) plus the damage and
// graveyard adds produced by the following turn's start-of-turn aura tick.
type TurnExtras struct {
	Value            int
	BestLine         []deck.CardAssignment
	TriggerDamage    int
	TriggerGraveyard []deck.Card
}

// EvalOneTurnForTesting drives one turn against the deck in source order (no shuffle) and
// returns the start-of-next-turn GameState plus the turn's extras. The returned state's
// Value() mirrors extras.Value. initial seeds the turn's carryover (Hero / Arsenal /
// Auras / Items / Banished / Graveyard / OpponentMarked); pass nil for a clean start
// inheriting the deck's hero.
func EvalOneTurnForTesting(masterDeck *deck.Deck, mp Matchup, initial *gameengine.GameState, initialHand []deck.Card) (*gameengine.GameState, TurnExtras) {
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
		return gameengine.GameStateBuilder().SetHero(h).Build(), TurnExtras{}
	}
	if initialHand == nil && d.Size() < handSize {
		return gameengine.GameStateBuilder().SetHero(h).Build(), TurnExtras{}
	}
	if initialHand != nil && (len(initialHand) == 0 || len(initialHand) > handSize) {
		return gameengine.GameStateBuilder().SetHero(h).Build(), TurnExtras{}
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
	sortHandByID(hand)
	weapons := make([]weapon.Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}
	play := best(weapons, hand, mp, d, master)
	extras := TurnExtras{
		Value:    play.Value,
		BestLine: append([]deck.CardAssignment(nil), play.BestLine...),
	}
	turnDraws := play.State.CardsDrawn()

	carry, _ := advanceToNextTurn(play, nil, nil)
	d = carry.deck
	master = carry.master
	held := carry.held

	if len(held) >= handSize || d.Size() < handSize-len(held) {
		// Can't refill — skip the start-of-next-turn tick.
		return snapshotState(master, d, held, play.Value, turnDraws), extras
	}

	// Refill before the tick so reveal-handling auras see the post-draw deck top.
	turn2Hand := make([]card.Card, 0, handSize+startOfTurnRevealRoom)
	turn2Hand = append(turn2Hand, held...)
	for _, c := range d.Draw(handSize - len(held)) {
		turn2Hand = append(turn2Hand, c.(card.Card))
	}
	preGrav := len(master.Graveyard())
	damage := processAurasAtStartOfTurn(master, d, &turn2Hand)
	extras.TriggerDamage = damage
	extras.TriggerGraveyard = cardsToDeckCards(master.Graveyard()[preGrav:])

	return snapshotState(master, d, turn2Hand, play.Value, turnDraws), extras
}
