package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestSunKiss_SoloIsHealOnly: with no Moon Wish in CardsPlayed the printed health-gain
// returns alone — no draw, no go-again grant. Pins the "synergy is opt-in, not unconditional"
// shape.
func TestSunKiss_SoloIsHealOnly(t *testing.T) {
	cases := []struct {
		c    card.Card
		heal int
	}{
		{cards.SunKissRed{}, 3},
		{cards.SunKissYellow{}, 2},
		{cards.SunKissBlue{}, 1},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.GenericAttack(0, 0)}).Build()}
		self := &card.CardState{Card: tc.c}
		ge.ResolveChainStep(ge.Logger(), self)
		got := ge.Value()
		if got != tc.heal {
			t.Errorf("%s: solo Play() = %d, want %d", tc.c.Name(), got, tc.heal)
		}
		if self.GrantedGoAgain {
			t.Errorf("%s: solo grant set GrantedGoAgain", tc.c.Name())
		}
		if h := ge.Hand(); len(h) != 0 {
			t.Errorf("%s: solo grant drew a card (got %d in Hand)", tc.c.Name(), len(h))
		}
	}
}

// Tests that a prior Moon Wish printing makes Sun Kiss credit the heal, draw a card, and
// grant self go again — across all (Sun Kiss × Moon Wish) variant combinations.
func TestSunKiss_SynergyFiresOnPriorMoonWish(t *testing.T) {
	moonWishVariants := []card.Card{cards.MoonWishRed{}, cards.MoonWishYellow{}, cards.MoonWishBlue{}}
	for _, mw := range moonWishVariants {
		for _, sk := range []struct {
			c    card.Card
			heal int
		}{
			{cards.SunKissRed{}, 3},
			{cards.SunKissYellow{}, 2},
			{cards.SunKissBlue{}, 1},
		} {
			ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.GenericAttack(0, 0)}).Build()}
			ge.SetCardsPlayed([]card.Card{mw})
			self := &card.CardState{Card: sk.c}
			ge.ResolveChainStep(ge.Logger(), self)
			got := ge.Value()
			if got != sk.heal {
				t.Errorf("%s after %s: Play() = %d, want %d (synergy still credits printed heal)",
					sk.c.Name(), mw.Name(), got, sk.heal)
			}
			if !self.GrantedGoAgain {
				t.Errorf("%s after %s: GrantedGoAgain = false, want true", sk.c.Name(), mw.Name())
			}
			if h := ge.Hand(); len(h) != 1 {
				t.Errorf("%s after %s: Hand len = %d, want 1 (one mid-turn draw)",
					sk.c.Name(), mw.Name(), len(h))
			}
		}
	}
}

// Tests that the Sun Kiss synergy only fires on a Moon Wish printing, not any prior attack.
func TestSunKiss_SynergyDoesNotFireOnUnrelatedAttacks(t *testing.T) {
	notMoonWish := testutils.GenericAttackPitch(0, 0, 1)
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.GenericAttack(0, 0)}).Build()}
	ge.SetCardsPlayed([]card.Card{notMoonWish})
	self := &card.CardState{Card: cards.SunKissRed{}}
	ge.ResolveChainStep(ge.Logger(), self)
	got := ge.Value()
	if got != 3 {
		t.Errorf("Play() = %d, want 3 (printed heal only)", got)
	}
	if self.GrantedGoAgain {
		t.Error("synergy fired on unrelated attack")
	}
	if h := ge.Hand(); len(h) != 0 {
		t.Errorf("synergy drew a card on unrelated attack (Hand len = %d, want 0)", len(h))
	}
}

// TestSunKiss_SynergyHandlesEmptyDeck: when the deck has been milled before Sun Kiss
// resolves, the synergy still grants go-again but the draw silently no-ops (DrawOne contract).
// Guards against a future regression that panics on Deck[0] read with no top.
func TestSunKiss_SynergyHandlesEmptyDeck(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsPlayed([]card.Card{cards.MoonWishRed{}}).Build()}
	self := &card.CardState{Card: cards.SunKissRed{}}
	ge.ResolveChainStep(ge.Logger(), self)
	got := ge.Value()
	if got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
	if !self.GrantedGoAgain {
		t.Error("GrantedGoAgain = false; synergy should still grant go again on empty deck")
	}
	if h := ge.Hand(); len(h) != 0 {
		t.Errorf("Hand len = %d, want 0 on empty deck", len(h))
	}
}
