package runeblade

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
		pc := &card.PlayedCard{Card: tc.c}
		s := card.TurnState{Self: pc}
		if got := tc.c.Play(&s); got != tc.wantDmg {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.wantDmg)
		}
		if pc.GrantedGoAgain {
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
// be) dealt this turn, so the conditional go again fires via Self.GrantedGoAgain.
func TestMeatAndGreet_ArcaneDamageDealtGrantsGoAgain(t *testing.T) {
	cases := []card.Card{
		MeatAndGreetRed{},
		MeatAndGreetYellow{},
		MeatAndGreetBlue{},
	}
	for _, c := range cases {
		pc := &card.PlayedCard{Card: c}
		s := card.TurnState{Self: pc, ArcaneDamageDealt: true}
		_ = c.Play(&s)
		if !pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain = false, want true (ArcaneDamageDealt → go again)", c.Name())
		}
	}
}
