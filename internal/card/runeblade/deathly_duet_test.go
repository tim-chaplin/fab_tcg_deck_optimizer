package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestDeathlyDuet_BaseDamage: nothing pitched → neither rider fires, printed power stands.
func TestDeathlyDuet_BaseDamage(t *testing.T) {
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
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

// TestDeathlyDuet_AttackPitchedBuffGated: the +2{p} buff is credited only when the buffed
// total lands in the likely-to-hit set. Blue (2+2=4) qualifies; Red (4+2=6) and Yellow (3+2=5)
// don't — the opponent comfortably blocks their buffed totals.
func TestDeathlyDuet_AttackPitchedBuffGated(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
		note string
	}{
		{DeathlyDuetRed{}, 4, "buffed total 6 blockable"},
		{DeathlyDuetYellow{}, 3, "buffed total 5 blockable"},
		{DeathlyDuetBlue{}, 2 + 2, "buffed total 4 likely to hit"},
	}
	for _, tc := range cases {
		s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s with attack pitched: Play() = %d, want %d (%s)", tc.c.Name(), got, tc.want, tc.note)
		}
	}
}

// TestDeathlyDuet_AttackPitchedRunechantsDontRescueBuff: the +2{p} is a physical buff;
// runechants firing alongside are separate arcane damage, so they can't rescue a blockable
// buffed total.
func TestDeathlyDuet_AttackPitchedRunechantsDontRescueBuff(t *testing.T) {
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}, Runechants: 1}
	if got := (DeathlyDuetRed{}).Play(&s); got != 4 {
		t.Errorf("Red with attack pitched + 1 Runechant: Play() = %d, want 4 (buff physical, arcane separate)", got)
	}
}

// TestDeathlyDuet_NonAttackActionPitchedCreatesRunechants: Runechants are arcane damage, not
// gated by physical hit likelihood. The 2 tokens credit regardless of blockability.
func TestDeathlyDuet_NonAttackActionPitchedCreatesRunechants(t *testing.T) {
	s := card.TurnState{Pitched: []card.Card{stubNonAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s); got != 4+2 {
		t.Errorf("Red with non-attack pitched: Play() = %d, want 6 (base 4 + 2 runechant credits)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
	if !s.AuraCreated {
		t.Errorf("AuraCreated should be set when Runechants are created")
	}
}

// TestDeathlyDuet_BothBranchesFire: both riders' conditions met. Red base 4 with attack pitched
// has buffed total 6 (blockable → don't credit +2); runechant creation credits +2 independently.
// Total = 4 + 2 = 6.
func TestDeathlyDuet_BothBranchesFire(t *testing.T) {
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}, stubNonAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s); got != 4+2 {
		t.Errorf("Red with both pitched: Play() = %d, want 6 (base 4 + 2 runechant credits; buff blocked)", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
}
