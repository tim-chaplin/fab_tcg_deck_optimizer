package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func TestDeathlyDuet_BaseDamage(t *testing.T) {
	// Nothing attributed → just printed power.
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.DeathlyDuetRed{}, 4},
		{cards.DeathlyDuetYellow{}, 3},
		{cards.DeathlyDuetBlue{}, 2},
	}
	for _, tc := range cases {
		var s sim.TurnState
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestDeathlyDuet_AttackAttributedAddsPower(t *testing.T) {
	// Attack attributed → +2{p}.
	var s sim.TurnState
	self := &card.CardState{
		Card:          cards.DeathlyDuetRed{},
		PitchedToPlay: []card.Card{testutils.RunebladeAttack{}},
	}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if got := s.Value(); got != 6 {
		t.Errorf("Deathly Duet Red with attack attributed: Play() = %d, want 6", got)
	}
}

func TestDeathlyDuet_NonAttackActionAttributedCreatesRunechants(t *testing.T) {
	// Non-attack action attributed → 2 Runechant tokens enter play, credited +1 each at creation.
	// Play returns base + 2 (Deathly Duet Red base 4 + 2 token credits = 6). state.Runechants=2
	// for downstream consume bookkeeping.
	var s sim.TurnState
	self := &card.CardState{
		Card:          cards.DeathlyDuetRed{},
		PitchedToPlay: []card.Card{testutils.NonAttack{}},
	}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if got := s.Value(); got != 6 {
		t.Errorf("Deathly Duet Red with non-attack attributed: Play() = %d, want 6 (base 4 + 2 token credits)", got)
	}
	if s.RunechantCount() != 2 {
		t.Errorf("Runechants = %d, want 2", s.RunechantCount())
	}
	if !s.AuraCreated() {
		t.Errorf("AuraCreated should be set when Runechants are created")
	}
}

func TestDeathlyDuet_BothBranchesFire(t *testing.T) {
	// Both an attack AND a non-attack action attributed → both riders fire: +2 power bonus,
	// plus 2 Runechants credited +1 each at creation. Play returns base 4 + 2 power + 2 = 8.
	var s sim.TurnState
	self := &card.CardState{
		Card:          cards.DeathlyDuetRed{},
		PitchedToPlay: []card.Card{testutils.RunebladeAttack{}, testutils.NonAttack{}},
	}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if got := s.Value(); got != 8 {
		t.Errorf("Deathly Duet Red with both attributed: Play() = %d, want 8 (base 4 + 2 power + 2 token credits)", got)
	}
	if s.RunechantCount() != 2 {
		t.Errorf("Runechants = %d, want 2", s.RunechantCount())
	}
}
