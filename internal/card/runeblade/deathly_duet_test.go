package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestDeathlyDuet_BaseDamage(t *testing.T) {
	// Nothing pitched → just printed power.
	cases := []struct {
		c    card.Card
		want int
	}{
		{DeathlyDuetRed{}, 4},
		{DeathlyDuetYellow{}, 3},
		{DeathlyDuetBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestDeathlyDuet_AttackPitchedAddsPower(t *testing.T) {
	// Attack pitched → +2{p}.
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s, &card.CardState{}); got != 6 {
		t.Errorf("Deathly Duet Red with attack pitched: Play() = %d, want 6", got)
	}
}

func TestDeathlyDuet_NonAttackActionPitchedCreatesRunechants(t *testing.T) {
	// Non-attack action pitched → 2 Runechant tokens enter play, credited +1 each at creation.
	// Play returns base + 2 (Deathly Duet Red base 4 + 2 token credits = 6). state.Runechants=2
	// for downstream consume bookkeeping.
	s := card.TurnState{Pitched: []card.Card{stubNonAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s, &card.CardState{}); got != 6 {
		t.Errorf("Deathly Duet Red with non-attack pitched: Play() = %d, want 6 (base 4 + 2 token credits)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
	if !s.AuraCreated {
		t.Errorf("AuraCreated should be set when Runechants are created")
	}
}

func TestDeathlyDuet_BothBranchesFire(t *testing.T) {
	// Both an attack AND a non-attack action in Pitched → both riders fire: +2 power bonus, plus
	// 2 Runechants credited +1 each at creation. Play returns base 4 + 2 power + 2 tokens = 8.
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}, stubNonAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s, &card.CardState{}); got != 8 {
		t.Errorf("Deathly Duet Red with both pitched: Play() = %d, want 8 (base 4 + 2 power + 2 token credits)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
}
