package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Play credits printed Attack with no aura in play (no bonus).
func TestYintiYanti_PlayNoAuraNoBonus(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{YintiYantiRed{}, 3},
		{YintiYantiYellow{}, 2},
		{YintiYantiBlue{}, 1},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		self := &sim.CardState{Card: tc.c}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if s.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d (printed attack, no aura bonus)", tc.c.Name(), s.Value, tc.want)
		}
	}
}

// Tests that Play credits printed Attack + 1 when an aura is in play.
func TestYintiYanti_PlayWithAuraGetsBonus(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{YintiYantiRed{}, 4},
		{YintiYantiYellow{}, 3},
		{YintiYantiBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{Auras: []sim.Aura{sim.NewRunechantAura(1)}}
		self := &sim.CardState{Card: tc.c}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if s.Value != tc.want {
			t.Errorf("%s with aura: Value = %d, want %d (printed +1 aura bonus)", tc.c.Name(), s.Value, tc.want)
		}
	}
}

// Tests that Block leaves BonusDefense untouched when no aura is in play.
func TestYintiYanti_BlockNoAuraNoBonus(t *testing.T) {
	for _, c := range []sim.Card{YintiYantiRed{}, YintiYantiYellow{}, YintiYantiBlue{}} {
		s := sim.TurnState{}
		self := &sim.CardState{Card: c}
		c.(sim.Blocker).Block(&s, s.Logger(), self)
		if self.BonusDefense != 0 {
			t.Errorf("%s: BonusDefense = %d, want 0 (no aura)", c.Name(), self.BonusDefense)
		}
	}
}

// Tests that Block bumps BonusDefense by 1 when an aura is in play.
func TestYintiYanti_BlockWithAuraGetsBonus(t *testing.T) {
	for _, c := range []sim.Card{YintiYantiRed{}, YintiYantiYellow{}, YintiYantiBlue{}} {
		s := sim.TurnState{Auras: []sim.Aura{sim.NewRunechantAura(1)}}
		self := &sim.CardState{Card: c}
		c.(sim.Blocker).Block(&s, s.Logger(), self)
		if self.BonusDefense != 1 {
			t.Errorf("%s with aura: BonusDefense = %d, want 1", c.Name(), self.BonusDefense)
		}
	}
}
