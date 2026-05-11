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
	sim.ResolveChainStep(&s, s.Logger(), &sim.CardState{Card: FlyingHighRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestFlyingHigh_NonAttackInRemainingFizzles confirms a non-attack action in CardsRemaining is
// skipped by the attack-action predicate.
func TestFlyingHigh_NonAttackInRemainingFizzles(t *testing.T) {
	skipped := &sim.CardState{Card: testutils.GenericAction()}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{skipped}}
	sim.ResolveChainStep(&s, s.Logger(), &sim.CardState{Card: FlyingHighRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
	if skipped.GrantedGoAgain {
		t.Error("non-attack target should not be granted go again")
	}
}

// Tests that the +1{p} rider only fires on a colour-matched target; go again is unconditional.
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
			sim.ResolveChainStep(&s, s.Logger(), &sim.CardState{Card: tc.c})
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

// Tests that "your next attack" grants go again to a weapon swing target, with no +1{p}
// rider (weapons have no pitch).
func TestFlyingHigh_GrantsGoAgainToWeaponSwing(t *testing.T) {
	pc := &sim.CardState{Card: testutils.RunebladeWeapon{}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{pc}}
	sim.ResolveChainStep(&s, s.Logger(), &sim.CardState{Card: FlyingHighRed{}})
	if !pc.GrantedGoAgain {
		t.Error("weapon swing should get go again ('your next attack' has no 'action card' qualifier)")
	}
	if pc.BonusAttack != 0 {
		t.Errorf("weapon BonusAttack = %d, want 0 (weapons have no pitch)", pc.BonusAttack)
	}
}
