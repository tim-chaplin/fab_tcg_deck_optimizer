package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// EvalTwoTurnsForTesting drives two turns in source order (no shuffle), threading
// cross-turn carryover so tests can assert that turn-1 state shapes turn 2. Returns the
// start-of-turn-3 GameState plus per-turn extras.
//
// initial seeds turn-1 carryover (nil for a clean start). hand1 optionally fixes turn 1's
// hand; pass nil to draw from the deck top. Turn 2's hand always comes from real
// carryover (held + fresh draws + start-of-turn reveals).
//
// Extras: Value / BestLine describe the turn's chain. TriggerDamage / TriggerGraveyard
// report the FOLLOWING turn's start-of-turn aura tick (extras1 → start-of-turn-2, extras2
// → start-of-turn-3). extras2.Value already folds in the start-of-turn-2 tick damage.
//
// Turn 1 runs no start-of-turn aura tick. To exercise start-of-turn behavior for turn
// N ≥ 2, install the aura via a turn-1 play and read extras2 / the returned gs.
//
// On failure the snapshot stops at the furthest boundary reached and missing turns get
// zero extras: validation failure → hero-only GameState; turn 2 can't field → end-of-1;
// turn 3 can't refill → end-of-2 with extras2.TriggerDamage zero.
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

	// Refill before the tick so reveal-handling auras see the post-draw deck top.
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
