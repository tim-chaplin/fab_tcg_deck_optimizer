package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestShrillOfSkullform_BaseDamage: without an aura played or created this turn the +3 buff
// can't fire.
func TestShrillOfSkullform_BaseDamage(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ShrillOfSkullformRed{}, 4},
		{ShrillOfSkullformYellow{}, 3},
		{ShrillOfSkullformBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		got := tc.c.Play(&s)
		if got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

// TestShrillOfSkullform_LikelyBuffedTotalCreditsBonus: Red (4+3=7) is the only variant whose
// buffed total lands in the likely set — the opponent can't comfortably block 7, so the buff
// delivers.
func TestShrillOfSkullform_LikelyBuffedTotalCreditsBonus(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	if got := (ShrillOfSkullformRed{}).Play(&s); got != 4+3 {
		t.Errorf("Red with aura played: Play() = %d, want 7 (buffed total 7 likely to hit)", got)
	}
}

// TestShrillOfSkullform_BlockableBuffedTotalSuppressesBonus: Yellow (3+3=6) and Blue (2+3=5)
// produce blockable buffed totals; the buff delivers nothing and isn't credited.
func TestShrillOfSkullform_BlockableBuffedTotalSuppressesBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{ShrillOfSkullformYellow{}, 3},
		{ShrillOfSkullformBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s with aura played: Play() = %d, want %d (buffed total blockable)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestShrillOfSkullform_RunechantsDontRescueBuff: the +3{p} is a physical buff; runechants
// firing alongside are a separate arcane stream and don't change whether the physical attack
// hits.
func TestShrillOfSkullform_RunechantsDontRescueBuff(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}, Runechants: 1}
	if got := (ShrillOfSkullformYellow{}).Play(&s); got != 3 {
		t.Errorf("Yellow with aura + 1 Runechant: Play() = %d, want 3 (buff physical, arcane separate)", got)
	}
}
