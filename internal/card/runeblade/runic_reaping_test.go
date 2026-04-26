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

// TestRunicReaping_PitchedAttackBuffPushesIntoHitWindow: a printed-3 target buffed to 4 by
// the pitched-attack +1{p} rider lands in the LikelyToHit window — both riders fire. The +1
// rides on the target's BonusAttack so Play returns just the Runechant credits, and
// EffectiveAttack already reflects the buff when the LikelyToHit check runs.
func TestRunicReaping_PitchedAttackBuffPushesIntoHitWindow(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{RunicReapingRed{}, 3},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubAttackWithPower{power: 3}}
		s := card.TurnState{
			CardsRemaining: []*card.CardState{target},
			Pitched:        []card.Card{stubRunebladeAttack{}},
		}
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d (Runechant tokens; +1 rides on target's BonusAttack)", tc.c.Name(), got, tc.n)
		}
		if target.BonusAttack != 1 {
			t.Errorf("%s: target BonusAttack = %d, want 1 (pitched-attack +1{p} rider)", tc.c.Name(), target.BonusAttack)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
	}
}

// TestRunicReaping_PitchedAttackBuffPushesPastHitWindow: a printed-4 target buffed to 5 by
// the pitched-attack +1{p} rider falls OUT of the LikelyToHit window (5 isn't in {1,4,7} and
// isn't dominate 5+). The Runechant rider doesn't fire even though it would have on the
// printed-4 attack alone — buffs flow through the LikelyToHit check, not around it.
func TestRunicReaping_PitchedAttackBuffPushesPastHitWindow(t *testing.T) {
	target := &card.CardState{Card: stubAttackWithPower{power: 4}}
	s := card.TurnState{
		CardsRemaining: []*card.CardState{target},
		Pitched:        []card.Card{stubRunebladeAttack{}},
	}
	if got := (RunicReapingRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (buffed to 5 = blockable; runechant rider drops)", got)
	}
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1 (pitched-attack +1{p} still grants)", target.BonusAttack)
	}
	if s.Runechants != 0 {
		t.Errorf("Runechants = %d, want 0 (buffed-out-of-window target drops the rider)", s.Runechants)
	}
}

// TestRunicReaping_BlockableTargetEvenAfterBuffDropsRunechants: a printed-2 target buffed to
// 3 by the pitched-attack +1{p} rider still doesn't hit (3 isn't in {1,4,7}). The +1 lands
// on the target's BonusAttack but the Runechant rider stays off.
func TestRunicReaping_BlockableTargetEvenAfterBuffDropsRunechants(t *testing.T) {
	target := &card.CardState{Card: stubAttackWithPower{power: 2}}
	s := card.TurnState{
		CardsRemaining: []*card.CardState{target},
		Pitched:        []card.Card{stubRunebladeAttack{}},
	}
	if got := (RunicReapingRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (buffed to 3 still blockable)", got)
	}
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1", target.BonusAttack)
	}
	if s.Runechants != 0 {
		t.Errorf("Runechants = %d, want 0", s.Runechants)
	}
	if s.AuraCreated {
		t.Error("AuraCreated should stay false when no Runechant is created")
	}
}
