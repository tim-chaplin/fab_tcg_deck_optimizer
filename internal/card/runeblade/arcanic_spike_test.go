package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestArcanicSpike_ArcaneDamageNotDealtReturnsBaseAttack covers the unsatisfied branch: when
// TurnState.ArcaneDamageDealt is false the +2{p} rider doesn't fire and Play returns only the
// printed attack.
func TestArcanicSpike_ArcaneDamageNotDealtReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ArcanicSpikeRed{}, 5},
		{ArcanicSpikeYellow{}, 4},
		{ArcanicSpikeBlue{}, 3},
	}
	for _, tc := range cases {
		s := card.TurnState{}
		if got := tc.c.Play(&s); got != tc.want {
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
		{ArcanicSpikeRed{}, 5 + 2},
		{ArcanicSpikeYellow{}, 4 + 2},
		{ArcanicSpikeBlue{}, 3 + 2},
	}
	for _, tc := range cases {
		s := card.TurnState{ArcaneDamageDealt: true}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (attack + arcane bonus)", tc.c.Name(), got, tc.want)
		}
	}
}
