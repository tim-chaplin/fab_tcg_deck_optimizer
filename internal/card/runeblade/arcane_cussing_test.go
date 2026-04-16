package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestArcaneCussing_NoFollowingAttackFlatValue covers the path where no same-turn attack exists
// to trigger the aura's destruction — Play returns a flat N as a future-turn payout without
// touching state.Runechants.
func TestArcaneCussing_NoFollowingAttackFlatValue(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{ArcaneCussingRed{}, 3},
		{ArcaneCussingYellow{}, 2},
		{ArcaneCussingBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState // empty CardsRemaining
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (no tokens tracked when no follow-up)", tc.c.Name(), s.Runechants)
		}
	}
}

// TestArcaneCussing_FollowingAttackCreatesRunechants covers the path where a same-turn attack
// will trigger Arcane Cussing's destruction, so the N Runechants enter state this turn and the
// attack downstream consumes them.
func TestArcaneCussing_FollowingAttackCreatesRunechants(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{ArcaneCussingRed{}, 3},
		{ArcaneCussingYellow{}, 2},
		{ArcaneCussingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}}}
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
	}
}

// TestArcaneCussing_FollowingWeaponCountsAsAttack — weapon swings in CardsRemaining also deal
// damage and trigger Cussing's destruction, so they count the same as an attack-action follow-up.
func TestArcaneCussing_FollowingWeaponCountsAsAttack(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeWeapon{}}}}
	if got := (ArcaneCussingRed{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
	if s.Runechants != 3 {
		t.Errorf("Runechants = %d, want 3", s.Runechants)
	}
}
