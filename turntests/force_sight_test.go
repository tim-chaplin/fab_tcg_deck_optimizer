package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestForceSight_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestForceSight_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{cards.ForceSightRed{}, cards.ForceSightYellow{}, cards.ForceSightBlue{}} {
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestForceSight_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestForceSight_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.FakeRedAction()}}).Build()}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.ForceSightRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestForceSight_NextAttackReturnsBonus: first attack-action triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestForceSight_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ForceSightRed{}, 3},
		{cards.ForceSightYellow{}, 2},
		{cards.ForceSightBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.FakeRedAttack()}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}

// Tests that Force Sight played from hand returns Value 0 — the arsenal-gated Opt rider
// doesn't fire in this path, and Force Sight isn't an attack so no damage is credited.
func TestForceSight_HandPlayValueZero(t *testing.T) {
	a, b := testutils.FakeRedAction().WithName("a"), testutils.FakeRedAction().WithName("b")
	for _, c := range []card.Card{cards.ForceSightRed{}, cards.ForceSightYellow{}, cards.ForceSightBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b}).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if ge.Value() != 0 {
			t.Errorf("%s: Play() from hand Value = %d, want 0", c.Name(), ge.Value())
		}
	}
}

// Tests that Force Sight played from arsenal returns Value 0 (the rider's effect is the
// deck reshape, not a value credit).
func TestForceSight_ArsenalPlayValueZero(t *testing.T) {
	a, b := testutils.FakeRedAction().WithName("a"), testutils.FakeRedAction().WithName("b")
	for _, c := range []card.Card{cards.ForceSightRed{}, cards.ForceSightYellow{}, cards.ForceSightBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b}).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c, FromArsenal: true})
		if ge.Value() != 0 {
			t.Errorf("%s: Play() from arsenal Value = %d, want 0", c.Name(), ge.Value())
		}
	}
}
