package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// TestRunebloodIncantation_PlayRegistersStartOfTurnTriggerWithCountN: Play flips AuraCreated
// and registers a TriggerStartOfTurn Aura with Count=N (Red 3, Yellow 2, Blue 1). No
// same-turn damage credit — every Runechant lands on a real future-turn fire.
func TestRunebloodIncantation_PlayRegistersStartOfTurnTriggerWithCountN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.RunebloodIncantationRed{}, 3},
		{cards.RunebloodIncantationYellow{}, 2},
		{cards.RunebloodIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (every rune fires on a future turn)", tc.c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if ge.RunechantCount() != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (no same-turn rune)", tc.c.Name(), ge.RunechantCount())
		}
		if len(ge.Auras()) != 1 {
			t.Fatalf("%s: Auras len = %d, want 1", tc.c.Name(), len(ge.Auras()))
		}
		tr := ge.Auras()[0]
		if tr.TriggerType() != triggertype.StartOfTurn {
			t.Errorf("%s: trigger Type = %d, want TriggerStartOfTurn", tc.c.Name(), tr.TriggerType())
		}
		if tr.Count() != tc.n {
			t.Errorf("%s: Count = %d, want %d (one per verse counter)", tc.c.Name(), tr.Count(), tc.n)
		}
	}
}

// TestRunebloodIncantation_HandlerCreatesOneRunechantPerFire: each invocation of the handler
// creates exactly one live Runechant — the multi-fire behaviour comes from the sim ticking
// Count, not from the handler doing more work each call.
func TestRunebloodIncantation_HandlerCreatesOneRunechantPerFire(t *testing.T) {
	for _, c := range []card.Card{cards.RunebloodIncantationRed{}, cards.RunebloodIncantationYellow{}, cards.RunebloodIncantationBlue{}} {
		play := gameengine.New()
		play.ResolveChainStep(play.Logger(), &card.CardState{Card: c})
		fire := gameengine.New()
		fire.CreateAura(play.Auras()[0])
		fire.FireStartOfTurn()
		if fire.Value() != 1 {
			t.Errorf("%s: handler Value = %d, want 1", c.Name(), fire.Value())
		}
		if fire.RunechantCount() != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (one rune per fire)", c.Name(), fire.RunechantCount())
		}
	}
}
