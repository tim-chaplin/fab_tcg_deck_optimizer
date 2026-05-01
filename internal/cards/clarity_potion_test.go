package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Clarity Potion's Play emits a LogPlay step and an Opt 2 log entry.
func TestClarityPotion_PlayCallsOpt2(t *testing.T) {
	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	prev := sim.CurrentHero
	sim.CurrentHero = testutils.Hero{}
	defer func() { sim.CurrentHero = prev }()

	c := ClarityPotionBlue{}
	s := sim.NewTurnState([]sim.Card{a, b}, nil)
	c.Play(s, &sim.CardState{Card: c})
	if s.Value != 0 {
		t.Errorf("Play() Value = %d, want 0", s.Value)
	}
	if len(s.Log) != 2 {
		t.Fatalf("Log len = %d, want 2 (LogPlay + Opted ...)", len(s.Log))
	}
	want := "Opted [a, b], put [a, b] on top, put [] on bottom"
	if got := s.Log[1].Text; got != want {
		t.Errorf("Opt log entry = %q, want %q", got, want)
	}
}
