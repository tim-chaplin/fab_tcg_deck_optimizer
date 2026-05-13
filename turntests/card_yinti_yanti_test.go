package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Play credits printed Attack with no aura in play (no bonus).
func TestYintiYanti_PlayNoAuraNoBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.YintiYantiRed{}, 3},
		{cards.YintiYantiYellow{}, 2},
		{cards.YintiYantiBlue{}, 1},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		self := &card.CardState{Card: tc.c}
		ge.ResolveChainStep(ge.Logger(), self)
		if ge.Value() != tc.want {
			t.Errorf("%s: Value = %d, want %d (printed attack, no aura bonus)", tc.c.Name(), ge.Value(), tc.want)
		}
	}
}

// Tests that Play credits printed Attack + 1 when an aura is in play.
func TestYintiYanti_PlayWithAuraGetsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.YintiYantiRed{}, 4},
		{cards.YintiYantiYellow{}, 3},
		{cards.YintiYantiBlue{}, 2},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.CreateAura(cards.NewRunechant(1))
		self := &card.CardState{Card: tc.c}
		ge.ResolveChainStep(ge.Logger(), self)
		if ge.Value() != tc.want {
			t.Errorf("%s with aura: Value = %d, want %d (printed +1 aura bonus)", tc.c.Name(), ge.Value(), tc.want)
		}
	}
}

// Tests that Block leaves BonusDefense untouched when no aura is in play.
func TestYintiYanti_BlockNoAuraNoBonus(t *testing.T) {
	for _, c := range []card.Card{cards.YintiYantiRed{}, cards.YintiYantiYellow{}, cards.YintiYantiBlue{}} {
		ge := gameengine.New()
		self := &card.CardState{Card: c}
		c.(card.Blocker).Block(ge, ge.Logger(), self)
		if self.BonusDefense != 0 {
			t.Errorf("%s: BonusDefense = %d, want 0 (no aura)", c.Name(), self.BonusDefense)
		}
	}
}

// Tests that Block bumps BonusDefense by 1 when an aura is in play.
func TestYintiYanti_BlockWithAuraGetsBonus(t *testing.T) {
	for _, c := range []card.Card{cards.YintiYantiRed{}, cards.YintiYantiYellow{}, cards.YintiYantiBlue{}} {
		ge := gameengine.New()
		ge.CreateAura(cards.NewRunechant(1))
		self := &card.CardState{Card: c}
		c.(card.Blocker).Block(ge, ge.Logger(), self)
		if self.BonusDefense != 1 {
			t.Errorf("%s with aura: BonusDefense = %d, want 1", c.Name(), self.BonusDefense)
		}
	}
}
