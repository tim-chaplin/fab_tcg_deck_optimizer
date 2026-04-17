package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestArcanicSpike_NoRunechantsReturnsBaseAttack covers the fix: if no Runechant exists at Play
// time there's no source of arcane damage this turn, so the +2{p} rider doesn't fire and Play
// returns only the printed attack — no flat bonus.
func TestArcanicSpike_NoRunechantsReturnsBaseAttack(t *testing.T) {
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
			t.Errorf("%s: Play() = %d, want %d (base attack, no runechants)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestArcanicSpike_RunechantsTriggerBonus exercises the satisfied path: any live Runechant means
// arcane damage will fire on this attack, so the +2{p} rider activates and Play returns attack + 2.
func TestArcanicSpike_RunechantsTriggerBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ArcanicSpikeRed{}, 5 + 2},
		{ArcanicSpikeYellow{}, 4 + 2},
		{ArcanicSpikeBlue{}, 3 + 2},
	}
	for _, tc := range cases {
		s := card.TurnState{Runechants: 1}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (attack + arcane bonus)", tc.c.Name(), got, tc.want)
		}
	}
}
