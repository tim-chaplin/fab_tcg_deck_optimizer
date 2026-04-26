package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfSuffering_PlayCreditsArcaneAndRider exercises the printed conditional: the Sigil's
// own arcane fires first and flips ArcaneDamageDealt, so the +1{d} rider always lands when the
// Sigil plays. Total Play credit is 1 (arcane) + 1 (rider) = 2 across every variant; the printed
// block stays in Defense() and is consumed separately by the chain.
func TestSigilOfSuffering_PlayCreditsArcaneAndRider(t *testing.T) {
	cases := []card.Card{SigilOfSufferingRed{}, SigilOfSufferingYellow{}, SigilOfSufferingBlue{}}
	for _, c := range cases {
		var s card.TurnState
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 2 {
			t.Errorf("%s: Play() = %d, want 2 (1 arcane + 1 rider)", card.DisplayName(c), got)
		}
		if !s.ArcaneDamageDealt {
			t.Errorf("%s: ArcaneDamageDealt = false, want true (Sigil's own arcane should flip the flag)", card.DisplayName(c))
		}
	}
}

// TestSigilOfSuffering_DefenseIsPrinted pins each variant's Defense() to its printed block value
// — the +1{d} bonus is credited as a Play-time rider, not baked into Defense.
func TestSigilOfSuffering_DefenseIsPrinted(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{SigilOfSufferingRed{}, 3},
		{SigilOfSufferingYellow{}, 2},
		{SigilOfSufferingBlue{}, 1},
	}
	for _, tc := range cases {
		if got := tc.c.Defense(); got != tc.want {
			t.Errorf("%s: Defense() = %d, want %d (printed)", card.DisplayName(tc.c), got, tc.want)
		}
	}
}
