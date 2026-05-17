package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TurnExtras carries the turn-1 results that don't sit naturally on the returned
// start-of-next-turn *gameengine.GameState: score, BestLine, and the damage / graveyard
// produced by the start-of-next-turn aura tick.
type TurnExtras struct {
	Value            int
	BestLine         []deck.CardAssignment
	TriggerDamage    int
	TriggerGraveyard []deck.Card
}

// EvalOneTurnForTesting drives one turn against the deck in source order (no shuffle) so
// tests can assert chain outcomes plus the start-of-next-turn state without the full
// multi-shuffle Evaluate loop. Returns the GameState at the start of the next turn plus
// turn-1 extras; the returned GameState's Value() mirrors extras.Value. initial seeds the
// turn's carryover (Hero / Arsenal / Auras / Items / Banished / Graveyard /
// OpponentMarked); pass nil to start clean inheriting the deck's hero.
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
	handBuf := make([]card.Card, handSize, handSize+startOfTurnRevealRoom)
	var hand []card.Card
	if initialHand != nil {
		hand = handBuf[:len(initialHand)]
		for i, c := range initialHand {
			hand[i] = c.(card.Card)
		}
	} else {
		hand = handBuf[:handSize]
		for i, c := range d.Draw(handSize) {
			hand[i] = c.(card.Card)
		}
	}
	sortHandByID(hand)
	weapons := make([]weapon.Weapon, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.(weapon.Weapon)
	}
	play := best(weapons, hand, mp, d, master)
	d = play.State.Deck()
	pitched := pitchedFromBestLine(play.BestLine)
	recycled := make([]deck.Card, len(pitched))
	for i, c := range pitched {
		recycled[i] = c
	}
	d.PutBottom(recycled)

	auraQueue := make([]*aura.Aura, 0, len(play.State.Auras()))
	for _, a := range play.State.Auras() {
		auraQueue = append(auraQueue, a.(*aura.Aura))
	}
	itemQueue := make([]*item.Item, 0, len(play.State.Items()))
	for _, it := range play.State.Items() {
		itemQueue = append(itemQueue, it.(*item.Item))
	}

	held := append([]card.Card(nil), play.State.Hand()...)
	graveyardOut := append([]card.Card(nil), play.State.Graveyard()...)
	extras := TurnExtras{
		Value:    play.Value,
		BestLine: append([]deck.CardAssignment(nil), play.BestLine...),
	}

	// No room or not enough deck to refill — skip the start-of-next-turn aura tick and
	// return the state as it stood at the end of turn 1.
	if len(held) >= handSize || d.Size() < handSize-len(held) {
		gs := gameengine.GameStateBuilder().
			SetHero(h).
			SetArsenal(play.State.Arsenal()).
			SetDeck(d.Copy()).
			SetGraveyard(graveyardOut).
			SetOpponentMarked(play.State.OpponentMarked()).
			SetValue(play.Value).
			Build()
		gs.SetCardsDrawn(play.State.CardsDrawn())
		for _, a := range auraQueue {
			gs.CreateAura(a)
		}
		for _, it := range itemQueue {
			gs.CreateItem(it)
		}
		return gs, extras
	}

	turn2Hand := append([]card.Card(nil), held...)
	for _, c := range d.Draw(handSize - len(held)) {
		turn2Hand = append(turn2Hand, c.(card.Card))
	}
	// Drive the start-of-next-turn aura tick on a fresh carryover state so reveals / damage
	// / destroyed-with-graveyard side effects are observable on the returned extras.
	tick := gameengine.GameStateBuilder().Build()
	for _, a := range auraQueue {
		tick.CreateAura(a)
	}
	preGrav := len(tick.Graveyard())
	trigDamage := processAurasAtStartOfTurn(tick, d, &turn2Hand)
	extras.TriggerDamage = trigDamage
	if newGrav := tick.Graveyard(); len(newGrav) > preGrav {
		extras.TriggerGraveyard = make([]deck.Card, 0, len(newGrav)-preGrav)
		for _, c := range newGrav[preGrav:] {
			extras.TriggerGraveyard = append(extras.TriggerGraveyard, c)
		}
	}

	gs := gameengine.GameStateBuilder().
		SetHero(h).
		SetArsenal(play.State.Arsenal()).
		SetDeck(d.Copy()).
		SetGraveyard(graveyardOut).
		SetOpponentMarked(play.State.OpponentMarked()).
		SetValue(play.Value).
		Build()
	gs.SetHand(turn2Hand)
	gs.SetCardsDrawn(play.State.CardsDrawn())
	for _, a := range tick.Auras() {
		gs.CreateAura(a.(*aura.Aura))
	}
	for _, it := range itemQueue {
		gs.CreateItem(it)
	}
	return gs, extras
}
