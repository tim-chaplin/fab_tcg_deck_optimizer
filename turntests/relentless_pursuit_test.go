package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Tests that Relentless Pursuit recycles to the bottom of the deck (rather than the
// graveyard) when an attack has already resolved this turn.
func TestRelentlessPursuit_RecyclesToDeckBottomAfterAttack(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		testutils.RedAttack{},         // attack first to satisfy the recycle gate
		cards.RelentlessPursuitBlue{}, // resolves second; recycles to deck bottom
		testutils.RedPitch{},          // funds RedAttack's cost-1
	}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.Prior{}, hand)
	rp := cards.RelentlessPursuitBlue{}
	rpName := rp.DisplayName()
	for _, c := range state.Graveyard {
		if c.ID() == rp.ID() {
			t.Fatalf("Relentless Pursuit unexpectedly in graveyard after recycle; got %v", state.Graveyard)
		}
	}
	if state.StartOfNextTurnDeck.NameCounts()[rpName] == 0 {
		t.Fatalf("Relentless Pursuit missing from next turn's deck; graveyard=%v", state.Graveyard)
	}
}

// Tests that Relentless Pursuit goes to the graveyard normally when it resolves before any
// attack this turn (no prior attack to recycle).
func TestRelentlessPursuit_GoesToGraveyardWithoutPriorAttack(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.RelentlessPursuitBlue{}, cards.OutedRed{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.Prior{}, hand)
	if state.Value != 4 {
		t.Fatalf("Value = %d, want 4 (RP marks, Outed reads mark for 3+1=4)", state.Value)
	}
	rpID := cards.RelentlessPursuitBlue{}.ID()
	for _, c := range state.Graveyard {
		if c.ID() == rpID {
			return
		}
	}
	t.Fatalf("Relentless Pursuit missing from graveyard (RP→Outed order, no prior attack); graveyard=%v", state.Graveyard)
}
