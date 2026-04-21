package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestZealousBelting_NoQualifyingPitchNoGoAgain covers the miss branch: no pitched card this turn
// has base power greater than the played variant's own base power, so the rider fizzles.
// Red's base power is 5 — a pitched power-5 card fails the strict ">" check.
func TestZealousBelting_NoQualifyingPitchNoGoAgain(t *testing.T) {
	c := ZealousBeltingRed{}
	s := card.TurnState{
		Pitched: []card.Card{stubGenericAttack(0, 5)},
	}
	if got := c.Play(&s); got != c.Attack() {
		t.Errorf("Play() = %d, want %d (no qualifying pitch)", got, c.Attack())
	}
	if s.SelfGoAgain {
		t.Errorf("SelfGoAgain = true, want false (no pitched card with power > %d)", c.Attack())
	}
}

// TestZealousBelting_HigherPowerPitchGrantsGoAgain exercises the hit branch: a pitched card whose
// base power is strictly greater than the variant's own base power triggers the go-again grant.
// Printed base powers are Red 5, Yellow 4, Blue 3.
func TestZealousBelting_HigherPowerPitchGrantsGoAgain(t *testing.T) {
	cases := []struct {
		c        card.Card
		pitchPow int
	}{
		{ZealousBeltingRed{}, 6},    // base 5, pitched power 6
		{ZealousBeltingYellow{}, 5}, // base 4, pitched power 5
		{ZealousBeltingBlue{}, 4},   // base 3, pitched power 4
	}
	for _, tc := range cases {
		s := card.TurnState{
			Pitched: []card.Card{stubGenericAttack(0, tc.pitchPow)},
		}
		_ = tc.c.Play(&s)
		if !s.SelfGoAgain {
			t.Errorf("%s: SelfGoAgain = false, want true (pitched power %d > base)", tc.c.Name(), tc.pitchPow)
		}
	}
}
