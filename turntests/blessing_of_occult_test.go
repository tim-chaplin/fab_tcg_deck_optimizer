package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
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
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (rune creation deferred to trigger)", tc.c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if ge.RunechantCount() != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (tokens are next-turn)", tc.c.Name(), ge.RunechantCount())
		}
		if len(ge.Auras()) != 1 {
			t.Fatalf("%s: Auras len = %d, want 1", tc.c.Name(), len(ge.Auras()))
		}
		if ge.Auras()[0].TriggerType() != triggertype.StartOfTurn {
			t.Errorf("%s: trigger Type = %d, want TriggerStartOfTurn", tc.c.Name(), ge.Auras()[0].TriggerType())
		}
		if ge.Auras()[0].Count() != tc.wantCount {
			t.Errorf("%s: Count = %d, want %d", tc.c.Name(), ge.Auras()[0].Count(), tc.wantCount)
		}
	}
}

// Tests that a carried Blessing of Occult fires at the start of the next turn, creating N
// Runechants (R=3, Y=2, B=1) and destroying itself.
func TestBlessingOfOccult_StartOfTurnCreatesNRunes(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.BlessingOfOccultRed{}, 3},
		{cards.BlessingOfOccultYellow{}, 2},
		{cards.BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		prior := gameengine.GameStateBuilder().CreateAuraFromCard(tc.c).Build()
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{testutils.FakeBlueResource()}

		summary := sim.EvalOneTurnForTesting(d, prior, hand)

		if got := summary.State.RunechantCount(); got != tc.n {
			t.Errorf("%s: Runechants = %d, want %d (start-of-turn handler created them)", tc.c.Name(), got, tc.n)
		}
		if got := len(summary.State.Auras()); got != 1 {
			t.Errorf("%s: Auras = %d, want 1 (Blessing destroyed itself, leaving only the consolidated Runechant entry)",
				tc.c.Name(), got)
		}
	}
}
