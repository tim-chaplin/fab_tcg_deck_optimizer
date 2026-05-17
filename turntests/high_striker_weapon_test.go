package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that High Striker's "next attack hits" rider fires on a follow-up weapon swing.
func TestHighStriker_WeaponHitCreatesCopper(t *testing.T) {
	d := deck.New(heroes.Viserai{}, []deck.Weapon{testutils.ClubWeapon{}}, fillerDeck())
	hand := []deck.Card{
		cards.HighStrikerRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if got := state.CopperCount(); got != 6 {
		t.Fatalf("Copper at start of next turn = %d, want 6 (HSR rider on Club swing)\nBestLine: %s",
			got, formatBestLine(state.BestLine))
	}
}
