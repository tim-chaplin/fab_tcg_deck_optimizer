package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Plunder Run from hand registers the on-hit-draw trigger and skips the +N{p} grant.
func TestPlunderRun_FromHandQueuesTriggerNoBonus(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(0, 4)}
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{target}})
	self := &card.CardState{Card: cards.PlunderRunRed{}}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if got := triggerHitCount(&s); got != 1 {
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
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{target}})
		self := &card.CardState{Card: tc.c, FromArsenal: true}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if got := triggerHitCount(&s); got != 1 {
			t.Errorf("%s: queued triggers = %d, want 1", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.wantBoon {
			t.Errorf("%s: target.BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.wantBoon)
		}
	}
}

// Multiple Plunder Runs queue independent triggers — they all fire on the same hit.
func TestPlunderRun_TriggersStack(t *testing.T) {
	s := sim.TurnState{}
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: cards.PlunderRunRed{}})
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: cards.PlunderRunBlue{}})
	if got := triggerHitCount(&s); got != 2 {
		t.Errorf("queued triggers = %d, want 2 (two independent listeners)", got)
	}
}
