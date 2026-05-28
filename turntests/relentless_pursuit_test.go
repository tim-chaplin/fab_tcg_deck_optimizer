package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Relentless Pursuit recycles to the bottom of the deck (rather than the
// graveyard) when an attack has already resolved this turn.
func TestRelentlessPursuit_RecyclesToDeckBottomAfterAttack(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, fillerDeck())
	hand := []card.Card{
		testutils.FakeRedAttack().
			WithCost(1).
			WithPower(3).
			WithGoAgain(), // attack first to satisfy the recycle gate
		cards.RelentlessPursuitBlue{}, // resolves second; recycles to deck bottom
		testutils.FakeRedAction(),     // funds RedAttack's cost-1
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)
	rp := cards.RelentlessPursuitBlue{}
	rpName := rp.DisplayName()
	for _, c := range summary.State.Graveyard() {
		if c.ID() == rp.ID() {
			t.Fatalf("Relentless Pursuit unexpectedly in graveyard after recycle; got %v", summary.State.Graveyard())
		}
	}
	if summary.State.Deck().NameCounts()[rpName] == 0 {
		t.Fatalf("Relentless Pursuit missing from next turn's deck; graveyard=%v", summary.State.Graveyard())
	}
}

// Tests that Relentless Pursuit goes to the graveyard normally when it resolves before any
// attack this turn (no prior attack to recycle).
func TestRelentlessPursuit_GoesToGraveyardWithoutPriorAttack(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, fillerDeck())
	hand := []card.Card{cards.RelentlessPursuitBlue{}, cards.OutedRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)
	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (RP marks, Outed reads mark for 3+1=4)", summary.Value)
	}
	rpID := cards.RelentlessPursuitBlue{}.ID()
	for _, c := range summary.State.Graveyard() {
		if c.ID() == rpID {
			return
		}
	}
	t.Fatalf("Relentless Pursuit missing from graveyard (RP→Outed order, no prior attack); graveyard=%v", summary.State.Graveyard())
}
