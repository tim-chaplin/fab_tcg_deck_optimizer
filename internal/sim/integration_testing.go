package sim

// Integration-test entry points that drive one or two turns in source order (no shuffle),
// letting tests assert chain outcomes and cross-turn state without the full Evaluate loop.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TurnExtras carries a turn's chain results plus the following turn's start-of-turn aura
// tick (damage credit and any auras that exhausted into graveyard). TriggerDamage /
// TriggerGraveyard are zero unless populated by a multi-turn entry point.
type TurnExtras struct {
	Value            int
	BestLine         []deck.CardAssignment
	TriggerDamage    int
	TriggerGraveyard []deck.Card
}

// EvalOneTurnForTesting drives one turn and returns the post end-of-turn-cleanup GameState
// (pitched cards recycled, arsenal refilled, next hand drawn) plus the turn's chain extras.
// initial seeds carryover (Hero / Arsenal / Auras / Items / Banished / Graveyard /
// OpponentMarked); nil means a clean start inheriting the deck's hero.
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
		// Not enough deck left to refill — return post-recycle state with held cards intact.
		return snapshotState(master, d, held, play.Value, turnDraws), extras
	}

	nextHand := make([]card.Card, 0, handSize)
	nextHand = append(nextHand, held...)
	for _, c := range d.Draw(handSize - len(held)) {
		nextHand = append(nextHand, c.(card.Card))
	}
	return snapshotState(master, d, nextHand, play.Value, turnDraws), extras
}

// EvalTwoTurnsForTesting drives two turns, threading carryover between them, and returns
// the start-of-turn-3 GameState plus per-turn extras. extras1's TriggerDamage/Graveyard
// report the start-of-turn-2 aura tick; extras2's report start-of-turn-3. Snapshot
// truncates at the furthest turn that ran when refills or validation fall short.
func EvalTwoTurnsForTesting(masterDeck *deck.Deck, mp Matchup, initial *gameengine.GameState, hand1 []deck.Card) (*gameengine.GameState, TurnExtras, TurnExtras) {
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
		return gameengine.GameStateBuilder().SetHero(h).Build(), TurnExtras{}, TurnExtras{}
	}
	if hand1 == nil && d.Size() < handSize {
		return gameengine.GameStateBuilder().SetHero(h).Build(), TurnExtras{}, TurnExtras{}
	}
	if hand1 != nil && (len(hand1) == 0 || len(hand1) > handSize) {
		return gameengine.GameStateBuilder().SetHero(h).Build(), TurnExtras{}, TurnExtras{}
	}
	weapons := make([]weapon.Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}

	hand := make([]card.Card, 0, handSize+startOfTurnRevealRoom)
	if hand1 != nil {
		for _, c := range hand1 {
			hand = append(hand, c.(card.Card))
		}
	} else {
		for _, c := range d.Draw(handSize) {
			hand = append(hand, c.(card.Card))
		}
	}
	sortHandByID(hand)
	play1 := best(weapons, hand, mp, d, master)
	extras1 := TurnExtras{
		Value:    play1.Value,
		BestLine: append([]deck.CardAssignment(nil), play1.BestLine...),
	}
	turn1Draws := play1.State.CardsDrawn()

	carry1, _ := advanceToNextTurn(play1, nil, nil)
	d = carry1.deck
	master = carry1.master
	held := carry1.held

	if len(held) >= handSize || d.Size() < handSize-len(held) {
		return snapshotState(master, d, held, play1.Value, turn1Draws), extras1, TurnExtras{}
	}

	turn2Hand := make([]card.Card, 0, handSize+startOfTurnRevealRoom)
	turn2Hand = append(turn2Hand, held...)
	for _, c := range d.Draw(handSize - len(held)) {
		turn2Hand = append(turn2Hand, c.(card.Card))
	}
	preGrav2 := len(master.Graveyard())
	damage2 := processAurasAtStartOfTurn(master, d, &turn2Hand)
	extras1.TriggerDamage = damage2
	extras1.TriggerGraveyard = cardsToDeckCards(master.Graveyard()[preGrav2:])
	sortHandByID(turn2Hand)

	play2 := best(weapons, turn2Hand, mp, d, master)
	play2.Value += damage2
	extras2 := TurnExtras{
		Value:    play2.Value,
		BestLine: append([]deck.CardAssignment(nil), play2.BestLine...),
	}
	turn2Draws := play2.State.CardsDrawn()

	carry2, _ := advanceToNextTurn(play2, nil, nil)
	d = carry2.deck
	master = carry2.master
	held3 := carry2.held

	if len(held3) >= handSize || d.Size() < handSize-len(held3) {
		return snapshotState(master, d, held3, play2.Value, turn2Draws), extras1, extras2
	}

	turn3Hand := make([]card.Card, 0, handSize+startOfTurnRevealRoom)
	turn3Hand = append(turn3Hand, held3...)
	for _, c := range d.Draw(handSize - len(held3)) {
		turn3Hand = append(turn3Hand, c.(card.Card))
	}
	preGrav3 := len(master.Graveyard())
	damage3 := processAurasAtStartOfTurn(master, d, &turn3Hand)
	extras2.TriggerDamage = damage3
	extras2.TriggerGraveyard = cardsToDeckCards(master.Graveyard()[preGrav3:])

	return snapshotState(master, d, turn3Hand, play2.Value, turn2Draws), extras1, extras2
}

// snapshotState clones src's persistent carryover and overlays deck (copied), hand,
// value, and cardsDrawn.
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

// cardsToDeckCards widens []card.Card to []deck.Card; returns nil for empty input.
func cardsToDeckCards(in []card.Card) []deck.Card {
	if len(in) == 0 {
		return nil
	}
	out := make([]deck.Card, len(in))
	for i, c := range in {
		out[i] = c
	}
	return out
}
