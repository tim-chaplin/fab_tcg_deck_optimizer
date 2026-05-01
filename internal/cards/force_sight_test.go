package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestForceSight_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestForceSight_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{ForceSightRed{}, ForceSightYellow{}, ForceSightBlue{}} {
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestForceSight_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestForceSight_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAction()}}}
	(ForceSightRed{}).Play(&s, &sim.CardState{Card: ForceSightRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestForceSight_NextAttackReturnsBonus: first attack-action triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestForceSight_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{ForceSightRed{}, 3},
		{ForceSightYellow{}, 2},
		{ForceSightBlue{}, 1},
	}
	for _, tc := range cases {
		target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
		s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}

// Tests that Force Sight played from hand skips the arsenal-gated Opt.
func TestForceSight_HandPlaySkipsOpt(t *testing.T) {
	prev := sim.CurrentHero
	sim.CurrentHero = testutils.Hero{}
	defer func() { sim.CurrentHero = prev }()

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	for _, c := range []sim.Card{ForceSightRed{}, ForceSightYellow{}, ForceSightBlue{}} {
		s := sim.NewTurnState([]sim.Card{a, b}, nil)
		c.Play(s, &sim.CardState{Card: c})
		if s.Value != 0 {
			t.Errorf("%s: Play() from hand Value = %d, want 0", c.Name(), s.Value)
		}
		// Just the LogPlay chain step, no Opt sub-entry.
		if len(s.Log) != 1 {
			t.Errorf("%s: Log len = %d, want 1 (LogPlay only — Opt arsenal-gated)",
				c.Name(), len(s.Log))
		}
	}
}

// Tests that Force Sight played from arsenal emits an Opt 2 log entry after LogPlay.
func TestForceSight_ArsenalPlayCallsOpt2(t *testing.T) {
	prev := sim.CurrentHero
	sim.CurrentHero = testutils.Hero{}
	defer func() { sim.CurrentHero = prev }()

	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	for _, c := range []sim.Card{ForceSightRed{}, ForceSightYellow{}, ForceSightBlue{}} {
		s := sim.NewTurnState([]sim.Card{a, b}, nil)
		c.Play(s, &sim.CardState{Card: c, FromArsenal: true})
		if s.Value != 0 {
			t.Errorf("%s: Play() from arsenal Value = %d, want 0", c.Name(), s.Value)
		}
		if len(s.Log) != 2 {
			t.Errorf("%s: Log len = %d, want 2 (LogPlay + Opted ...)", c.Name(), len(s.Log))
			continue
		}
		want := "Opted [a, b], put [a, b] on top, put [] on bottom"
		if got := s.Log[1].Text; got != want {
			t.Errorf("%s: Opt log entry = %q, want %q", c.Name(), got, want)
		}
	}
}
