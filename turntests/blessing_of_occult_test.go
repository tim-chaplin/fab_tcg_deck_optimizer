package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Tests that Play flips AuraCreated, makes no runes this turn, and registers an aura with
// the per-variant Count for the deferred trigger.
func TestBlessingOfOccult_PlayCreatesAuraNoThisTurnRunes(t *testing.T) {
	cases := []struct {
		c         card.Card
		wantCount int
	}{
		{cards.BlessingOfOccultRed{}, 3},
		{cards.BlessingOfOccultYellow{}, 2},
		{cards.BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (rune creation deferred to trigger)", tc.c.Name(), got)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.RunechantCount() != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (tokens are next-turn)", tc.c.Name(), s.RunechantCount())
		}
		if len(s.Auras()) != 1 {
			t.Fatalf("%s: Auras len = %d, want 1", tc.c.Name(), len(s.Auras()))
		}
		if s.Auras()[0].TriggerType() != triggertype.StartOfTurn {
			t.Errorf("%s: trigger Type = %d, want TriggerStartOfTurn", tc.c.Name(), s.Auras()[0].TriggerType())
		}
		if s.Auras()[0].Count() != tc.wantCount {
			t.Errorf("%s: Count = %d, want %d", tc.c.Name(), s.Auras()[0].Count(), tc.wantCount)
		}
	}
}

// TestBlessingOfOccult_TriggerHandlerCreatesNRunes: invoking the trigger's handler on a
// fresh TurnState creates N live Runechants and credits matching damage.
func TestBlessingOfOccult_TriggerHandlerCreatesNRunes(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.BlessingOfOccultRed{}, 3},
		{cards.BlessingOfOccultYellow{}, 2},
		{cards.BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		play := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		play.ResolveChainStep(play.Logger(), &card.CardState{Card: tc.c})
		next := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		next.CreateAura(play.Auras()[0])
		next.FireStartOfTurn(nil)
		if next.Value() != tc.n {
			t.Errorf("%s: handler Value = %d, want %d", tc.c.Name(), next.Value(), tc.n)
		}
		if next.RunechantCount() != tc.n {
			t.Errorf("%s: Runechants = %d, want %d (live tokens on next turn)",
				tc.c.Name(), next.RunechantCount(), tc.n)
		}
	}
}
