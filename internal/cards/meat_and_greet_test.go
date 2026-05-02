package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the on-hit Runechant rider fires only on likely-hit variants.
func TestMeatAndGreet_OnHitRunechantGatedByLikelyToHit(t *testing.T) {
	cases := []struct {
		c       sim.Card
		wantDmg int
	}{
		{MeatAndGreetRed{}, 4 + 1},
		{MeatAndGreetYellow{}, 3},
		{MeatAndGreetBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		self := &sim.CardState{Card: tc.c}
		tc.c.Play(&s, self)
		testutils.FireOnHitIfLikely(&s, self)
		if got := s.Value; got != tc.wantDmg {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.wantDmg)
		}
		if self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true, want false (no prior arcane damage → no go again)", tc.c.Name())
		}
		// Card's printed GoAgain must also be false — the rider is the only source.
		if tc.c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (rider is conditional, not printed)", tc.c.Name())
		}
	}
}

// Tests that ArcaneDamageDealt at Play time grants conditional go again.
func TestMeatAndGreet_ArcaneDamageDealtGrantsGoAgain(t *testing.T) {
	cases := []sim.Card{
		MeatAndGreetRed{},
		MeatAndGreetYellow{},
		MeatAndGreetBlue{},
	}
	for _, c := range cases {
		s := sim.TurnState{ArcaneDamageDealt: true}
		self := &sim.CardState{Card: c}
		c.Play(&s, self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (ArcaneDamageDealt → go again)", c.Name())
		}
	}
}
