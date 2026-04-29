package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestFlyingHigh_NoAttackReturnsZero covers the miss branch: with nothing attack-typed in
// CardsRemaining the grant fizzles and Play returns 0.
func TestFlyingHigh_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	(FlyingHighRed{}).Play(&s, &sim.CardState{Card: FlyingHighRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestFlyingHigh_NonAttackInRemainingFizzles confirms a non-attack action in CardsRemaining is
// skipped by the attack-action predicate.
func TestFlyingHigh_NonAttackInRemainingFizzles(t *testing.T) {
	skipped := &sim.CardState{Card: testutils.GenericAction()}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{skipped}}
	(FlyingHighRed{}).Play(&s, &sim.CardState{Card: FlyingHighRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
	if skipped.GrantedGoAgain {
		t.Error("non-attack target should not be granted go again")
	}
}

// TestFlyingHigh_ColorMatchGrantsBonus: each variant's '+1{p} if matching color' rider only
// fires when the granted target's pitch matches this card's own pitch. Every variant grants
// go again to any attack target regardless. Granter returns 0; the +1 (when applicable)
// rides on the target's BonusAttack.
func TestFlyingHigh_ColorMatchGrantsBonus(t *testing.T) {
	cases := []struct {
		name       string
		c          sim.Card
		wantRed    int
		wantYellow int
		wantBlue   int
	}{
		{"FlyingHighRed", FlyingHighRed{}, 1, 0, 0},
		{"FlyingHighYellow", FlyingHighYellow{}, 0, 1, 0},
		{"FlyingHighBlue", FlyingHighBlue{}, 0, 0, 1},
	}
	for _, tc := range cases {
		for _, target := range []struct {
			pitch int
			want  int
		}{{1, tc.wantRed}, {2, tc.wantYellow}, {3, tc.wantBlue}} {
			pc := &sim.CardState{Card: testutils.GenericAttackPitch(0, 0, target.pitch)}
			s := sim.TurnState{CardsRemaining: []*sim.CardState{pc}}
			tc.c.Play(&s, &sim.CardState{Card: tc.c})
			if got := s.Value; got != 0 {
				t.Errorf("%s vs pitch-%d target: Play() = %d, want 0 (granter returns 0; +1 rides on target's BonusAttack when colour matches)",
					tc.name, target.pitch, got)
			}
			if pc.BonusAttack != target.want {
				t.Errorf("%s vs pitch-%d target: BonusAttack = %d, want %d",
					tc.name, target.pitch, pc.BonusAttack, target.want)
			}
			if !pc.GrantedGoAgain {
				t.Errorf("%s vs pitch-%d target: GrantedGoAgain = false, want true (go again is unconditional)",
					tc.name, target.pitch)
			}
		}
	}
}

// TestFlyingHigh_GrantsGoAgainToWeaponSwing pins the "your next attack" wording: the next
// scheduled attack can be a weapon swing (TypeWeapon, no TypeAction), and Flying High
// must grant go again to it. Weapons have no printed pitch so the "+1{p} if matching
// colour" rider never fires; only the go-again grant lands.
func TestFlyingHigh_GrantsGoAgainToWeaponSwing(t *testing.T) {
	pc := &sim.CardState{Card: testutils.RunebladeWeapon{}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{pc}}
	(FlyingHighRed{}).Play(&s, &sim.CardState{Card: FlyingHighRed{}})
	if !pc.GrantedGoAgain {
		t.Error("weapon swing should get go again ('your next attack' has no 'action card' qualifier)")
	}
	if pc.BonusAttack != 0 {
		t.Errorf("weapon BonusAttack = %d, want 0 (weapons have no pitch)", pc.BonusAttack)
	}
}
