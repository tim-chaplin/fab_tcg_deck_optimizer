package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestDeathlyDuet_BaseDamage(t *testing.T) {
	// Nothing attributed → just printed power.
	cases := []struct {
		c    sim.Card
		want int
	}{
		{DeathlyDuetRed{}, 4},
		{DeathlyDuetYellow{}, 3},
		{DeathlyDuetBlue{}, 2},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestDeathlyDuet_AttackAttributedAddsPower(t *testing.T) {
	// Attack attributed → +2{p}.
	var s sim.TurnState
	self := &sim.CardState{
		Card:          DeathlyDuetRed{},
		PitchedToPlay: []sim.Card{testutils.RunebladeAttack{}},
	}
	(DeathlyDuetRed{}).Play(&s, self)
	if got := s.Value; got != 6 {
		t.Errorf("Deathly Duet Red with attack attributed: Play() = %d, want 6", got)
	}
}

func TestDeathlyDuet_NonAttackActionAttributedCreatesRunechants(t *testing.T) {
	// Non-attack action attributed → 2 Runechant tokens enter play, credited +1 each at creation.
	// Play returns base + 2 (Deathly Duet Red base 4 + 2 token credits = 6). state.Runechants=2
	// for downstream consume bookkeeping.
	var s sim.TurnState
	self := &sim.CardState{
		Card:          DeathlyDuetRed{},
		PitchedToPlay: []sim.Card{testutils.NonAttack{}},
	}
	(DeathlyDuetRed{}).Play(&s, self)
	if got := s.Value; got != 6 {
		t.Errorf("Deathly Duet Red with non-attack attributed: Play() = %d, want 6 (base 4 + 2 token credits)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
	if !s.AuraCreated {
		t.Errorf("AuraCreated should be set when Runechants are created")
	}
}

func TestDeathlyDuet_BothBranchesFire(t *testing.T) {
	// Both an attack AND a non-attack action attributed → both riders fire: +2 power bonus,
	// plus 2 Runechants credited +1 each at creation. Play returns base 4 + 2 power + 2 = 8.
	var s sim.TurnState
	self := &sim.CardState{
		Card:          DeathlyDuetRed{},
		PitchedToPlay: []sim.Card{testutils.RunebladeAttack{}, testutils.NonAttack{}},
	}
	(DeathlyDuetRed{}).Play(&s, self)
	if got := s.Value; got != 8 {
		t.Errorf("Deathly Duet Red with both attributed: Play() = %d, want 8 (base 4 + 2 power + 2 token credits)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
}
