package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that with enough IncomingDamage Value sums printed Defense + +1{d} boost + 1 arcane.
func TestSigilOfSuffering_FullCreditWhenIncomingAbsorbsBoost(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.SigilOfSufferingRed{}, 5},    // 3 block + 1 boost + 1 arcane
		{cards.SigilOfSufferingYellow{}, 4}, // 2 block + 1 boost + 1 arcane
		{cards.SigilOfSufferingBlue{}, 3},   // 1 block + 1 boost + 1 arcane
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(10).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play(IncomingDamage=10) Value = %d, want %d (block + boost + arcane)",
				tc.c.DisplayName(), got, tc.want)
		}
	}
}

// Tests that when IncomingDamage equals printed Defense the +1{d} boost is clamped away;
// Value collapses to printed Defense + 1 arcane.
func TestSigilOfSuffering_BoostWastedWhenIncomingMatchesDefense(t *testing.T) {
	cases := []struct {
		c        card.Card
		incoming int
		want     int
	}{
		{cards.SigilOfSufferingRed{}, 3, 4},    // 3 block + 1 arcane (boost wasted)
		{cards.SigilOfSufferingYellow{}, 2, 3}, // 2 block + 1 arcane
		{cards.SigilOfSufferingBlue{}, 1, 2},   // 1 block + 1 arcane
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(tc.incoming).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play(IncomingDamage=%d) Value = %d, want %d (block at cap + arcane only)",
				tc.c.DisplayName(), tc.incoming, got, tc.want)
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
		{cards.SigilOfSufferingRed{}, 3},
		{cards.SigilOfSufferingYellow{}, 2},
		{cards.SigilOfSufferingBlue{}, 1},
	}
	for _, tc := range cases {
		if got := tc.c.Defense(); got != tc.want {
			t.Errorf("%s: Defense() = %d, want %d (printed)", tc.c.DisplayName(), got, tc.want)
		}
	}
}
