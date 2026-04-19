package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestArcanicSpike_ArcaneDamageNotDealtReturnsBaseAttack: the +2{p} rider is gated on the
// arcane-damage clause; without it the printed attack stands alone.
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

// TestArcanicSpike_LikelyBuffedTotalCreditsBonus: Red (5+2=7) is the only variant whose buffed
// total lands in the likely-to-hit set ({1,4,7}). The opponent can't comfortably block 7, so the
// +2 buff delivers and we credit it.
func TestArcanicSpike_LikelyBuffedTotalCreditsBonus(t *testing.T) {
	s := card.TurnState{ArcaneDamageDealt: true}
	if got := (ArcanicSpikeRed{}).Play(&s); got != 5+2 {
		t.Errorf("Red with ArcaneDamageDealt: Play() = %d, want 7 (5+2 likely to hit)", got)
	}
}

// TestArcanicSpike_BlockableBuffedTotalSuppressesBonus: Yellow (4+2=6) and Blue (3+2=5) produce
// buffed totals the opponent comfortably blocks. The buff delivers nothing, so we don't credit
// it — only the base attack stands.
func TestArcanicSpike_BlockableBuffedTotalSuppressesBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ArcanicSpikeYellow{}, 4},
		{ArcanicSpikeBlue{}, 3},
	}
	for _, tc := range cases {
		s := card.TurnState{ArcaneDamageDealt: true}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s with ArcaneDamageDealt: Play() = %d, want %d (buffed total blockable)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestArcanicSpike_RunechantsDontRescueBuff: the +2{p} is a physical buff; the opponent blocks
// the physical attack based on physical damage. Runechants firing alongside are a separate
// arcane stream and don't change the block decision, so they can't rescue a blockable buff.
func TestArcanicSpike_RunechantsDontRescueBuff(t *testing.T) {
	s := card.TurnState{ArcaneDamageDealt: true, Runechants: 1}
	if got := (ArcanicSpikeYellow{}).Play(&s); got != 4 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 4 (buff is physical, arcane separate)", got)
	}
}
