package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMeatAndGreet_NoRunechantsNoGoAgain covers the fix: without a live Runechant at Play time,
// the conditional "this gets go again" rider doesn't trigger. GoAgain() returns false and
// Self.GrantedGoAgain stays unset. The attack damage still includes the on-hit Runechant.
func TestMeatAndGreet_NoRunechantsNoGoAgain(t *testing.T) {
	cases := []struct {
		c        card.Card
		wantDmg  int
	}{
		{MeatAndGreetRed{}, 4 + 1},
		{MeatAndGreetYellow{}, 3 + 1},
		{MeatAndGreetBlue{}, 2 + 1},
	}
	for _, tc := range cases {
		pc := &card.PlayedCard{Card: tc.c}
		s := card.TurnState{Self: pc}
		if got := tc.c.Play(&s); got != tc.wantDmg {
			t.Errorf("%s: Play() = %d, want %d (attack + created Runechant)", tc.c.Name(), got, tc.wantDmg)
		}
		if pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = true, want false (no prior runechant → no go again)", tc.c.Name())
		}
		// Card's printed GoAgain must also be false — the rider is the only source.
		if tc.c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (rider is conditional, not printed)", tc.c.Name())
		}
	}
}

// TestMeatAndGreet_RunechantsGrantGoAgain exercises the satisfied branch: a live Runechant at the
// start of Play means arcane damage has been / will be dealt this turn, so the conditional go
// again fires via Self.GrantedGoAgain.
func TestMeatAndGreet_RunechantsGrantGoAgain(t *testing.T) {
	cases := []card.Card{
		MeatAndGreetRed{},
		MeatAndGreetYellow{},
		MeatAndGreetBlue{},
	}
	for _, c := range cases {
		pc := &card.PlayedCard{Card: c}
		s := card.TurnState{Self: pc, Runechants: 1}
		_ = c.Play(&s)
		if !pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (runechant live at Play → go again)", c.Name())
		}
	}
}
