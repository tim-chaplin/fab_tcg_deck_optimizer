package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestForceSight_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestForceSight_NoAttackReturnsZero(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.NewState()}
	for _, c := range []card.Card{cards.ForceSightRed{}, cards.ForceSightYellow{}, cards.ForceSightBlue{}} {
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestForceSight_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestForceSight_NonAttackInRemainingFizzles(t *testing.T) {
	s := gameengine.NewFromSpec(gameengine.Spec{CardsRemaining: []*card.CardState{{Card: testutils.GenericAction()}}})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.ForceSightRed{}})
	if got := s.Value(); got != 0 {
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
		target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
		s := gameengine.NewFromSpec(gameengine.Spec{CardsRemaining: []*card.CardState{target}})
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}

// Tests that Force Sight played from hand skips the arsenal-gated Opt.
func TestForceSight_HandPlaySkipsOpt(t *testing.T) {

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	for _, c := range []card.Card{cards.ForceSightRed{}, cards.ForceSightYellow{}, cards.ForceSightBlue{}} {
		s := gameengine.NewFromCards([]card.Card{a, b}, nil)
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if s.Value() != 0 {
			t.Errorf("%s: Play() from hand Value = %d, want 0", c.Name(), s.Value())
		}
		// Just the LogPlay chain step, no Opt sub-entry.
		if len(s.LogEntries()) != 1 {
			t.Errorf("%s: Log len = %d, want 1 (LogPlay only — Opt arsenal-gated)",
				c.Name(), len(s.LogEntries()))
		}
	}
}

// Tests that Force Sight played from arsenal emits an Opt 2 log entry after LogPlay.
func TestForceSight_ArsenalPlayCallsOpt2(t *testing.T) {

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	for _, c := range []card.Card{cards.ForceSightRed{}, cards.ForceSightYellow{}, cards.ForceSightBlue{}} {
		s := gameengine.NewFromCards([]card.Card{a, b}, nil)
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c, FromArsenal: true})
		if s.Value() != 0 {
			t.Errorf("%s: Play() from arsenal Value = %d, want 0", c.Name(), s.Value())
		}
		if len(s.LogEntries()) != 2 {
			t.Errorf("%s: Log len = %d, want 2 (Opted... + chain step)", c.Name(), len(s.LogEntries()))
			continue
		}
		// Play emits the Opted line during arsenal-gated resolution; the chain step is
		// auto-appended after Play returns, so the Opted entry lands first.
		want := "Opted [a, b], put [a, b] on top, put [] on bottom"
		if got := s.LogEntries()[0].Text; got != want {
			t.Errorf("%s: Opt log entry = %q, want %q", c.Name(), got, want)
		}
	}
}
