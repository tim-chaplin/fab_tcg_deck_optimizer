package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestArcaneCussing_BlockCoversIncomingReturnsN confirms the aura's value is N when the
// partition's block total meets or exceeds incoming damage — we don't take damage, the aura
// survives to pay out later.
func TestArcaneCussing_BlockCoversIncomingReturnsN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{ArcaneCussingRed{}, 3},
		{ArcaneCussingYellow{}, 2},
		{ArcaneCussingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 3}
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d (block == incoming)", tc.c.Name(), got, tc.n)
		}
	}
}

// TestArcaneCussing_OverBlockReturnsN pins that the BlockTotal is uncapped — over-blocking still
// counts as covering incoming and the aura survives.
func TestArcaneCussing_OverBlockReturnsN(t *testing.T) {
	s := card.TurnState{IncomingDamage: 3, BlockTotal: 7}
	if got := (ArcaneCussingRed{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3 (over-block still covers)", got)
	}
}

// TestArcaneCussing_BlockShortReturnsZero confirms the aura collapses to 0 when any incoming
// damage gets through — we take damage, aura dies without pay-out.
func TestArcaneCussing_BlockShortReturnsZero(t *testing.T) {
	cases := []card.Card{
		ArcaneCussingRed{},
		ArcaneCussingYellow{},
		ArcaneCussingBlue{},
	}
	for _, c := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 2}
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (block < incoming)", c.Name(), got)
		}
	}
}

// TestArcaneCussing_NoIncomingReturnsN covers the low-incoming edge: with 0 incoming damage,
// any BlockTotal (including 0) satisfies the guard and the aura pays its full N.
func TestArcaneCussing_NoIncomingReturnsN(t *testing.T) {
	s := card.TurnState{IncomingDamage: 0, BlockTotal: 0}
	if got := (ArcaneCussingRed{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3 (no incoming to kill the aura)", got)
	}
}
