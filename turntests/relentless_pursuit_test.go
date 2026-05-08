package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Relentless Pursuit recycles to the bottom of the deck (rather than the
// graveyard) when an attack has already resolved this turn. Driven through 2 turns so
// the post-recycle deck shows RP after applyTurnResult routes pitches through.
func TestRelentlessPursuit_RecyclesToDeckBottomAfterAttack(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		testutils.RedAttack{},         // attack first to satisfy the recycle gate
		cards.RelentlessPursuitBlue{}, // resolves second; recycles to deck bottom
		testutils.RedPitch{},          // funds RedAttack's cost-1
	}
	got := d.EvalTwoTurnsForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	rpID := cards.RelentlessPursuitBlue{}.ID()
	for _, c := range got.Turn1.State.Graveyard {
		if c.ID() == rpID {
			t.Fatalf("Relentless Pursuit unexpectedly in graveyard after recycle; got %v", got.Turn1.State.Graveyard)
		}
	}
	// After turn 2's deal off the post-recycle deck, RP either remains in turn 2's deck or
	// got dealt into turn 2's hand — either confirms RP was on the deck post turn 1.
	for _, c := range got.Turn2.DealtHand {
		if c.ID() == rpID {
			return
		}
	}
	for _, c := range got.Turn2.State.Deck {
		if c.ID() == rpID {
			return
		}
	}
	t.Fatalf("Relentless Pursuit missing from turn 2's deal/deck; turn2 hand=%v deck=%v turn1 graveyard=%v",
		got.Turn2.DealtHand, got.Turn2.State.Deck, got.Turn1.State.Graveyard)
}

// Tests that Relentless Pursuit goes to the graveyard normally when it resolves before any
// attack this turn (no prior attack to recycle).
func TestRelentlessPursuit_GoesToGraveyardWithoutPriorAttack(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.RelentlessPursuitBlue{}, cards.OutedRed{}}
	summary := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	state := summary.State
	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (RP marks, Outed reads mark for 3+1=4)", summary.Value)
	}
	rpID := cards.RelentlessPursuitBlue{}.ID()
	for _, c := range state.Graveyard {
		if c.ID() == rpID {
			return
		}
	}
	t.Fatalf("Relentless Pursuit missing from graveyard (RP→Outed order, no prior attack); graveyard=%v", state.Graveyard)
}
