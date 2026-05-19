package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Plunder Run from hand registers the on-hit-draw trigger and skips the +N{p} grant.
func TestPlunderRun_FromHandQueuesTriggerNoBonus(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(0, 4)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	pc := &card.CardState{Card: cards.PlunderRunRed{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if got := triggerHitCount(ge); got != 1 {
		t.Errorf("queued triggers = %d, want 1", got)
	}
	if target.BonusAttack != 0 {
		t.Errorf("target.BonusAttack = %d, want 0 (hand-played skips +N{p})", target.BonusAttack)
	}
}

// From arsenal: registers the trigger and adds +N{p} to the next attack action in
// CardsRemaining. Each printing carries its own N.
func TestPlunderRun_FromArsenalAddsBonusAttack(t *testing.T) {
	cases := []struct {
		c        card.Card
		wantBoon int
	}{
		{cards.PlunderRunRed{}, 3},
		{cards.PlunderRunYellow{}, 2},
		{cards.PlunderRunBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 4)}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		pc := &card.CardState{Card: tc.c, FromArsenal: true}
		ge.ResolveChainStep(ge.Logger(), pc)
		if got := triggerHitCount(ge); got != 1 {
			t.Errorf("%s: queued triggers = %d, want 1", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.wantBoon {
			t.Errorf("%s: target.BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.wantBoon)
		}
	}
}

// Multiple Plunder Runs queue independent triggers — they all fire on the same hit.
func TestPlunderRun_TriggersStack(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.PlunderRunRed{}})
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.PlunderRunBlue{}})
	if got := triggerHitCount(ge); got != 2 {
		t.Errorf("queued triggers = %d, want 2 (two independent listeners)", got)
	}
}
