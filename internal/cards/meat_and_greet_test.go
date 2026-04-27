package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMeatAndGreet_OnHitRunechantGatedByLikelyToHit: the Runechant rider fires only when the
// variant's printed power satisfies card.LikelyToHit. Red (4) qualifies and gets +1 for the
// token; Yellow (3) and Blue (2) are blockable and drop the rider.
func TestMeatAndGreet_OnHitRunechantGatedByLikelyToHit(t *testing.T) {
	cases := []struct {
		c       card.Card
		wantDmg int
	}{
		{MeatAndGreetRed{}, 4 + 1},
		{MeatAndGreetYellow{}, 3},
		{MeatAndGreetBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{}
		self := &card.CardState{Card: tc.c}
		tc.c.Play(&s, self)
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

// TestMeatAndGreet_ArcaneDamageDealtGrantsGoAgain exercises the satisfied branch: when
// ArcaneDamageDealt is set at the start of Play, arcane damage has already been (or is about to
// be) dealt this turn, so the conditional go again fires via self.GrantedGoAgain.
func TestMeatAndGreet_ArcaneDamageDealtGrantsGoAgain(t *testing.T) {
	cases := []card.Card{
		MeatAndGreetRed{},
		MeatAndGreetYellow{},
		MeatAndGreetBlue{},
	}
	for _, c := range cases {
		s := card.TurnState{ArcaneDamageDealt: true}
		self := &card.CardState{Card: c}
		c.Play(&s, self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (ArcaneDamageDealt → go again)", c.Name())
		}
	}
}
