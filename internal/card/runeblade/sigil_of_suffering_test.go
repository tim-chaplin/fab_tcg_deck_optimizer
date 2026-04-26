package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfSuffering_FullCreditWhenIncomingAbsorbsBoost: with enough IncomingDamage to consume
// the printed Defense plus the +1{d} bonus, total Value reflects every component: printed
// Defense + 1 (the +1{d} rider folded into BonusDefense) + 1 (arcane sub-line). Each variant
// scales by its printed Defense (Red 3, Yellow 2, Blue 1).
func TestSigilOfSuffering_FullCreditWhenIncomingAbsorbsBoost(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{SigilOfSufferingRed{}, 5},    // 3 block + 1 boost + 1 arcane
		{SigilOfSufferingYellow{}, 4}, // 2 block + 1 boost + 1 arcane
		{SigilOfSufferingBlue{}, 3},   // 1 block + 1 boost + 1 arcane
	}
	for _, tc := range cases {
		s := card.TurnState{IncomingDamage: 10}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play(IncomingDamage=10) Value = %d, want %d (block + boost + arcane)",
				card.DisplayName(tc.c), got, tc.want)
		}
	}
}

// TestSigilOfSuffering_BoostWastedWhenIncomingMatchesDefense: with IncomingDamage exactly equal
// to the printed Defense, the +1{d} bonus is over-block — ApplyAndLogEffectiveDefense's clamp
// drops the extra point. Total Value collapses to printed Defense + 1 (arcane) — the boost
// adds nothing because there's no more incoming for it to consume.
func TestSigilOfSuffering_BoostWastedWhenIncomingMatchesDefense(t *testing.T) {
	cases := []struct {
		c        card.Card
		incoming int
		want     int
	}{
		{SigilOfSufferingRed{}, 3, 4},    // 3 block + 1 arcane (boost wasted)
		{SigilOfSufferingYellow{}, 2, 3}, // 2 block + 1 arcane
		{SigilOfSufferingBlue{}, 1, 2},   // 1 block + 1 arcane
	}
	for _, tc := range cases {
		s := card.TurnState{IncomingDamage: tc.incoming}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play(IncomingDamage=%d) Value = %d, want %d (block at cap + arcane only)",
				card.DisplayName(tc.c), tc.incoming, got, tc.want)
		}
	}
}

// TestSigilOfSuffering_DefenseIsPrinted pins each variant's Defense() to its printed block value
// — the +1{d} bonus is credited via BonusDefense at Play time, not baked into Defense.
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
