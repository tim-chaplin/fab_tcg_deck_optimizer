package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// TestArcanicSpike_ArcaneDamageNotDealtReturnsBaseAttack covers the unsatisfied branch: when
// TurnState.ArcaneDamageDealt is false the +2{p} rider doesn't fire and Play returns only the
// printed attack.
func TestArcanicSpike_ArcaneDamageNotDealtReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ArcanicSpikeRed{}, 5},
		{cards.ArcanicSpikeYellow{}, 4},
		{cards.ArcanicSpikeBlue{}, 3},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, ArcaneDamageDealt=false)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestArcanicSpike_ArcaneDamageDealtTriggersBonus exercises the satisfied path: when
// ArcaneDamageDealt is set (an earlier attack fired a Runechant, or a direct-arcane card flipped
// the flag) the +2{p} rider activates and Play returns attack + 2.
func TestArcanicSpike_ArcaneDamageDealtTriggersBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ArcanicSpikeRed{}, 5 + 2},
		{cards.ArcanicSpikeYellow{}, 4 + 2},
		{cards.ArcanicSpikeBlue{}, 3 + 2},
	}
	for _, tc := range cases {
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{ArcaneDamageDealt: true})
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (attack + arcane bonus)", tc.c.Name(), got, tc.want)
		}
	}
}
