package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestRunicReaping_NoNextAttackReturnsZero(t *testing.T) {
	// No attack action following → no bonus at all, and AuraCreated must remain false.
	s := card.TurnState{Pitched: []card.Card{stubAttackWithPower{power: 4}}}
	if got := (RunicReapingRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Fatalf("want 0 when no next attack, got %d", got)
	}
	if s.AuraCreated {
		t.Fatalf("AuraCreated should stay false when no bonus fires")
	}
}

func TestRunicReaping_WeaponNextDoesNotQualify(t *testing.T) {
	// A Runeblade weapon swing later in the turn is not an attack action card, so the rider doesn't
	// trigger.
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubRunebladeWeapon{}}}}
	if got := (RunicReapingRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Fatalf("want 0 with weapon-only next, got %d", got)
	}
}

func TestRunicReaping_LikelyHitTargetNoPitchedAttack(t *testing.T) {
	// Target's printed power (4) is in the likely-to-hit set and nothing attack-typed was
	// pitched → N Runechant tokens created. Play returns N (each token credited +1 at creation);
	// the pitched-attack +1 rider doesn't fire.
	cases := []struct {
		c card.Card
		n int
	}{
		{RunicReapingRed{}, 3},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []*card.CardState{{Card: stubAttackWithPower{power: 4}}},
			Pitched:        []card.Card{stubNonAttack{}},
		}
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set when bonus fires", tc.c.Name())
		}
	}
}

func TestRunicReaping_LikelyHitTargetWithPitchedAttack(t *testing.T) {
	// Target is likely-to-hit AND an attack card was pitched → Play returns N (token credits) plus
	// 1 (the pitched-attack rider). state.Runechants holds only the N tokens — the rider damage is
	// direct, not a runechant.
	cases := []struct {
		c card.Card
		n int
	}{
		{RunicReapingRed{}, 3},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []*card.CardState{{Card: stubAttackWithPower{power: 4}}},
			Pitched:        []card.Card{stubRunebladeAttack{}},
		}
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.n+1 {
			t.Errorf("%s: Play() = %d, want %d (N tokens + 1 pitched-attack bonus)", tc.c.Name(), got, tc.n+1)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
	}
}

// TestRunicReaping_BlockableTargetDropsRunechants pins the LikelyToHit gate: when the target's
// printed power (3) falls in the blockable range, the "if this hits" Runechant clause fizzles.
// The pitched-attack +1{p} rider still fires because it isn't gated on hitting.
func TestRunicReaping_BlockableTargetDropsRunechants(t *testing.T) {
	s := card.TurnState{
		CardsRemaining: []*card.CardState{{Card: stubAttackWithPower{power: 3}}},
		Pitched:        []card.Card{stubRunebladeAttack{}},
	}
	if got := (RunicReapingRed{}).Play(&s, &card.CardState{}); got != 1 {
		t.Errorf("Play() = %d, want 1 (blockable target drops Runechants, pitched-attack +1 still fires)", got)
	}
	if s.Runechants != 0 {
		t.Errorf("Runechants = %d, want 0 (no tokens when target is blockable)", s.Runechants)
	}
	if s.AuraCreated {
		t.Error("AuraCreated should stay false when no Runechant is created")
	}
}
