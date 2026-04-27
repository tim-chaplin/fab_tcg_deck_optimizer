package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestFlyingHigh_NoAttackReturnsZero covers the miss branch: with nothing attack-typed in
// CardsRemaining the grant fizzles and Play returns 0.
func TestFlyingHigh_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	(FlyingHighRed{}).Play(&s, &card.CardState{Card: FlyingHighRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestFlyingHigh_NonAttackInRemainingFizzles confirms a non-attack action in CardsRemaining is
// skipped by the attack-action predicate.
func TestFlyingHigh_NonAttackInRemainingFizzles(t *testing.T) {
	skipped := &card.CardState{Card: testutils.GenericAction()}
	s := card.TurnState{CardsRemaining: []*card.CardState{skipped}}
	(FlyingHighRed{}).Play(&s, &card.CardState{Card: FlyingHighRed{}})
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
		c          card.Card
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
			pc := &card.CardState{Card: testutils.GenericAttackPitch(0, 0, target.pitch)}
			s := card.TurnState{CardsRemaining: []*card.CardState{pc}}
			tc.c.Play(&s, &card.CardState{Card: tc.c})
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
