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

// TestArcanicSpike_RunechantRescuesBlockableBuff: a lone Runechant firing alongside is likely to
// slip through, which counts as the attack connecting and credits the buff.
func TestArcanicSpike_RunechantRescuesBlockableBuff(t *testing.T) {
	s := card.TurnState{ArcaneDamageDealt: true, Runechants: 1}
	if got := (ArcanicSpikeYellow{}).Play(&s); got != 4+2 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 6 (runechant slips → buff credited)", got)
	}
}
