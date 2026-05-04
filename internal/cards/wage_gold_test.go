package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Wage Gold registers a wager handler when the player holds a Gold at
// Play time. With no Gold to wager, no handler should register — there's nothing
// to give over if the attack is fully blocked.
func TestWageGold_RegistersOnlyWhenGoldHeld(t *testing.T) {
	t.Run("no gold: no handler", func(t *testing.T) {
		s := sim.NewTurnState(nil, nil)
		pc := &sim.CardState{Card: WageGoldRed{}}
		(WageGoldRed{}).Play(s, pc)
		if got := len(pc.OnHit); got != 0 {
			t.Fatalf("OnHit entries = %d, want 0 (no Gold to wager)", got)
		}
	})
	t.Run("gold held: handler registered", func(t *testing.T) {
		s := sim.NewTurnState(nil, nil)
		s.CreateGold(1)
		pc := &sim.CardState{Card: WageGoldRed{}}
		(WageGoldRed{}).Play(s, pc)
		if got := len(pc.OnHit); got != 1 {
			t.Fatalf("OnHit entries = %d, want 1 (wager declared)", got)
		}
		if pc.OnHit[0].Fire != nil {
			t.Errorf("Fire (on-hit) = non-nil, want nil (no modelled effect on hit)")
		}
		if pc.OnHit[0].FireBlocked == nil {
			t.Errorf("FireBlocked (on-blocked) = nil, want non-nil (wager-loss path)")
		}
	})
}

// Tests the FireBlocked branch directly: a fully-blocked Wage Gold should hand over
// one Gold token and credit -1 Value to capture the opponent's eventual draw.
func TestWageGold_FireBlockedConsumesGoldAndCreditsMinusOne(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	s.CreateGold(2)
	pc := &sim.CardState{Card: WageGoldRed{}}
	(WageGoldRed{}).Play(s, pc)
	if len(pc.OnHit) != 1 || pc.OnHit[0].FireBlocked == nil {
		t.Fatalf("Wage Gold didn't register the FireBlocked handler with Gold in play")
	}
	startValue := s.Value
	pc.OnHit[0].FireBlocked(s, pc, &pc.OnHit[0])
	if s.Gold() != 1 {
		t.Errorf("Gold after wager loss = %d, want 1 (started at 2, gave one over)", s.Gold())
	}
	if delta := s.Value - startValue; delta != -1 {
		t.Errorf("Value delta = %d, want -1 (modelled opponent's draw)", delta)
	}
}

// Tests that FireBlocked is a no-op if the wagered Gold was consumed mid-chain (e.g.
// by a Gold-token activated ability between Play and finalize-active-attack). We can't
// hand over a Gold we no longer hold, so don't credit -1 either.
func TestWageGold_FireBlockedNoOpWhenGoldAlreadyGone(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	s.CreateGold(1)
	pc := &sim.CardState{Card: WageGoldRed{}}
	(WageGoldRed{}).Play(s, pc)
	s.ConsumeItem(sim.TokenTypeGold, 1)
	startValue := s.Value
	pc.OnHit[0].FireBlocked(s, pc, &pc.OnHit[0])
	if s.Gold() != 0 {
		t.Errorf("Gold = %d, want 0", s.Gold())
	}
	if delta := s.Value - startValue; delta != 0 {
		t.Errorf("Value delta = %d, want 0 (no Gold left to give over)", delta)
	}
}
