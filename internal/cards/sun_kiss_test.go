package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestSunKiss_SoloIsHealOnly: with no Moon Wish in CardsPlayed the printed health-gain
// returns alone — no draw, no go-again grant. Pins the "synergy is opt-in, not unconditional"
// shape.
func TestSunKiss_SoloIsHealOnly(t *testing.T) {
	cases := []struct {
		c    sim.Card
		heal int
	}{
		{SunKissRed{}, 3},
		{SunKissYellow{}, 2},
		{SunKissBlue{}, 1},
	}
	for _, tc := range cases {
		s := sim.NewTurnState([]sim.Card{testutils.GenericAttack(0, 0)}, nil)
		self := &sim.CardState{Card: tc.c}
		tc.c.Play(s, self)
		got := s.Value
		if got != tc.heal {
			t.Errorf("%s: solo Play() = %d, want %d", tc.c.Name(), got, tc.heal)
		}
		if self.GrantedGoAgain {
			t.Errorf("%s: solo grant set GrantedGoAgain", tc.c.Name())
		}
		if len(s.Hand) != 0 {
			t.Errorf("%s: solo grant drew a card (got %d in Hand)", tc.c.Name(), len(s.Hand))
		}
	}
}

// TestSunKiss_SynergyFiresOnPriorMoonWish: with any Moon Wish printing earlier in CardsPlayed,
// Play returns the printed heal AND draws a card AND grants self go again. Covers all three
// (Sun Kiss variant × Moon Wish variant) cross-products with one representative each so a
// future Moon Wish printing renaming would surface as a clear failure.
func TestSunKiss_SynergyFiresOnPriorMoonWish(t *testing.T) {
	moonWishVariants := []sim.Card{MoonWishRed{}, MoonWishYellow{}, MoonWishBlue{}}
	for _, mw := range moonWishVariants {
		for _, sk := range []struct {
			c    sim.Card
			heal int
		}{
			{SunKissRed{}, 3},
			{SunKissYellow{}, 2},
			{SunKissBlue{}, 1},
		} {
			s := sim.NewTurnState([]sim.Card{testutils.GenericAttack(0, 0)}, nil)
			s.CardsPlayed = []sim.Card{mw}
			self := &sim.CardState{Card: sk.c}
			sk.c.Play(s, self)
			got := s.Value
			if got != sk.heal {
				t.Errorf("%s after %s: Play() = %d, want %d (synergy still credits printed heal)",
					sk.c.Name(), mw.Name(), got, sk.heal)
			}
			if !self.GrantedGoAgain {
				t.Errorf("%s after %s: GrantedGoAgain = false, want true", sk.c.Name(), mw.Name())
			}
			if len(s.Hand) != 1 {
				t.Errorf("%s after %s: Hand len = %d, want 1 (one mid-turn draw)",
					sk.c.Name(), mw.Name(), len(s.Hand))
			}
		}
	}
}

// TestSunKiss_SynergyDoesNotFireOnUnrelatedAttacks: only a Moon Wish printing should trigger
// the synergy. A different attack (Arcanic Spike) earlier in CardsPlayed must not fire it.
// Sentinel for the name-prefix scan: if the predicate ever loosens to "any attack with 'wish'
// in the name" or similar, this catches it.
func TestSunKiss_SynergyDoesNotFireOnUnrelatedAttacks(t *testing.T) {
	notMoonWish := testutils.GenericAttackPitch(0, 0, 1)
	s := sim.NewTurnState([]sim.Card{testutils.GenericAttack(0, 0)}, nil)
	s.CardsPlayed = []sim.Card{notMoonWish}
	self := &sim.CardState{Card: SunKissRed{}}
	SunKissRed{}.Play(s, self)
	got := s.Value
	if got != 3 {
		t.Errorf("Play() = %d, want 3 (printed heal only)", got)
	}
	if self.GrantedGoAgain {
		t.Error("synergy fired on unrelated attack")
	}
	if len(s.Hand) != 0 {
		t.Errorf("synergy drew a card on unrelated attack (Hand len = %d, want 0)", len(s.Hand))
	}
}

// TestSunKiss_SynergyHandlesEmptyDeck: when the deck has been milled before Sun Kiss
// resolves, the synergy still grants go-again but the draw silently no-ops (DrawOne contract).
// Guards against a future regression that panics on Deck[0] read with no top.
func TestSunKiss_SynergyHandlesEmptyDeck(t *testing.T) {
	s := &sim.TurnState{
		CardsPlayed: []sim.Card{MoonWishRed{}},
		// Deck intentionally nil.
	}
	self := &sim.CardState{Card: SunKissRed{}}
	SunKissRed{}.Play(s, self)
	got := s.Value
	if got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
	if !self.GrantedGoAgain {
		t.Error("GrantedGoAgain = false; synergy should still grant go again on empty deck")
	}
	if len(s.Hand) != 0 {
		t.Errorf("Hand len = %d, want 0 on empty deck", len(s.Hand))
	}
}
