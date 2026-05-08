package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that a same-turn banish (Jack Be Quick consuming a graveyard Nimblism) flips
// CardBanished so a later Tremor plays at base + 2{p}. Optimal line: pitch Bauble for
// {3} → Jack Be Quick (banishes Nimblism, +1{p} and go again, deals 4) → Tremor (cost 1,
// deals 4 + 2 bonus = 6). Total 10.
func TestTremorOfIArathael_SameTurnBanishActivatesBonus(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.TremorOfIArathaelRed{},
		cards.JackBeQuickRed{},
		cards.TitaniumBaubleBlue{},
	}
	initial := sim.NewTurnStateFromSpec(sim.TurnStateSpec{
		Graveyard: []sim.Card{cards.NimblismRed{}},
	})
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, initial, hand).Value
	if got != 10 {
		t.Fatalf("Value = %d, want 10 (Bauble pitch → JBQ banishes Nimblism for 4 → Tremor 4+2=6)", got)
	}
}

// Tests that a card already in the banished zone (carryover from a prior turn) does NOT
// flip CardBanished — Tremor plays at its base power. Hand: Tremor, Bauble. Initial
// banish has one card seeded; CardBanished stays false. Total 4.
func TestTremorOfIArathael_PriorTurnBanishedZoneDoesNotActivate(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.TremorOfIArathaelRed{}, cards.TitaniumBaubleBlue{}}
	initial := sim.NewTurnStateFromSpec(sim.TurnStateSpec{
		Banished: []sim.Card{cards.NimblismRed{}},
	})
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, initial, hand).Value
	if got != 4 {
		t.Fatalf("Value = %d, want 4 (Tremor base 4; prior-turn banish doesn't trigger bonus)", got)
	}
}
