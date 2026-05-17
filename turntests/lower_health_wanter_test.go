package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
)

// stubLowHeroOn implements gameengine.LowerHealthWanter — used to exercise the "hero opts
// in" branch.
type stubLowHeroOn struct{}

func (stubLowHeroOn) ID() ids.HeroID                                           { return ids.InvalidHero }
func (stubLowHeroOn) Name() string                                             { return "stubLowHeroOn" }
func (stubLowHeroOn) Intelligence() int                                        { return 4 }
func (stubLowHeroOn) Types() card.TypeSet                                      { return 0 }
func (stubLowHeroOn) Class() card.CardType                                     { return 0 }
func (stubLowHeroOn) OnCardPlayed(card.Card, hero.GameEngine, hero.Logger) int { return 0 }
func (stubLowHeroOn) Opt(cards []card.Card) (top, bottom []card.Card)          { return cards, nil }
func (stubLowHeroOn) WantsLowerHealth()                                        {}

// stubLowHeroOff does NOT implement gameengine.LowerHealthWanter — the default branch.
type stubLowHeroOff struct{}

func (stubLowHeroOff) ID() ids.HeroID                                           { return ids.InvalidHero }
func (stubLowHeroOff) Name() string                                             { return "stubLowHeroOff" }
func (stubLowHeroOff) Intelligence() int                                        { return 4 }
func (stubLowHeroOff) Types() card.TypeSet                                      { return 0 }
func (stubLowHeroOff) Class() card.CardType                                     { return 0 }
func (stubLowHeroOff) OnCardPlayed(card.Card, hero.GameEngine, hero.Logger) int { return 0 }
func (stubLowHeroOff) Opt(cards []card.Card) (top, bottom []card.Card)          { return cards, nil }

// engineWithHero returns a fresh empty engine with hero installed.
func engineWithHero(h gameengine.Hero) *gameengine.GameEngine {
	ge := gameengine.New()
	ge.SetHero(h)
	return ge
}

// TestLowerHealthWanter_DamageRiders checks the +3{p} / +1{p} / +1{h} damage riders fire
// iff the current hero opts into gameengine.LowerHealthWanter.
func TestLowerHealthWanter_DamageRiders(t *testing.T) {
	cases := []struct {
		name    string
		card    card.Card
		wantOff int
		wantOn  int
	}{
		{"AdrenalineRushRed +3p", cards.AdrenalineRushRed{}, 4, 4 + 3},
		{"AdrenalineRushYellow +3p", cards.AdrenalineRushYellow{}, 3, 3 + 3},
		{"AdrenalineRushBlue +3p", cards.AdrenalineRushBlue{}, 2, 2 + 3},
		{"WoundedBullRed +1p", cards.WoundedBullRed{}, 7, 7 + 1},
		{"WoundedBullYellow +1p", cards.WoundedBullYellow{}, 6, 6 + 1},
		{"WoundedBullBlue +1p", cards.WoundedBullBlue{}, 5, 5 + 1},
		{"FyendalsFightingSpiritRed +1h", cards.FyendalsFightingSpiritRed{}, 7, 7 + 1},
		{"FyendalsFightingSpiritYellow +1h", cards.FyendalsFightingSpiritYellow{}, 6, 6 + 1},
		{"FyendalsFightingSpiritBlue +1h", cards.FyendalsFightingSpiritBlue{}, 5, 5 + 1},
	}
	for _, tc := range cases {
		sOff := engineWithHero(stubLowHeroOff{})
		sOff.ResolveChainStep(sOff.Logger(), &card.CardState{Card: tc.card})
		if got := sOff.Value(); got != tc.wantOff {
			t.Errorf("%s: Play() off = %d, want %d (hero does not opt in)", tc.name, got, tc.wantOff)
		}
		sOn := engineWithHero(stubLowHeroOn{})
		sOn.ResolveChainStep(sOn.Logger(), &card.CardState{Card: tc.card})
		if got := sOn.Value(); got != tc.wantOn {
			t.Errorf("%s: Play() on = %d, want %d (hero opts in)", tc.name, got, tc.wantOn)
		}
	}
}

// TestLowerHealthWanter_GoAgainRiders checks the conditional go-again flips iff the
// current hero opts into gameengine.LowerHealthWanter.
func TestLowerHealthWanter_GoAgainRiders(t *testing.T) {
	cards := []card.Card{
		cards.ScarForAScarRed{}, cards.ScarForAScarYellow{}, cards.ScarForAScarBlue{},
		cards.BlowForABlowRed{},
		cards.LifeForALifeRed{}, cards.LifeForALifeYellow{}, cards.LifeForALifeBlue{},
	}
	gOff := engineWithHero(stubLowHeroOff{})
	for _, c := range cards {
		if c.GoAgain(gOff) {
			t.Errorf("%s: GoAgain() = true with hero off, want false", c.Name())
		}
	}
	gOn := engineWithHero(stubLowHeroOn{})
	for _, c := range cards {
		if !c.GoAgain(gOn) {
			t.Errorf("%s: GoAgain() = false with hero on, want true", c.Name())
		}
	}
}

// TestLowerHealthWanter_NilHeroIsOff guards the startup / unset-hero case.
func TestLowerHealthWanter_NilHeroIsOff(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.AdrenalineRushRed{}})
	if got := ge.Value(); got != 4 {
		t.Errorf("AdrenalineRushRed nil-hero Play() = %d, want 4", got)
	}
	if (cards.ScarForAScarRed{}).GoAgain(ge) {
		t.Errorf("ScarForAScarRed nil-hero GoAgain() = true, want false")
	}
}

// Tests that Pound for Pound flips self.GrantedDominate iff the current hero opts into
// LowerHealthWanter.
func TestLowerHealthWanter_PoundForPoundDominateGrant(t *testing.T) {
	cards := []card.Card{cards.PoundForPoundRed{}, cards.PoundForPoundYellow{}, cards.PoundForPoundBlue{}}

	for _, c := range cards {
		self := &card.CardState{Card: c}
		s := engineWithHero(stubLowHeroOff{})
		s.ResolveChainStep(s.Logger(), self)
		if self.GrantedDominate {
			t.Errorf("%s: GrantedDominate = true with hero off, want false", c.Name())
		}
		if self.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = true with hero off, want false", c.Name())
		}
	}
	for _, c := range cards {
		self := &card.CardState{Card: c}
		s := engineWithHero(stubLowHeroOn{})
		s.ResolveChainStep(s.Logger(), self)
		if !self.GrantedDominate {
			t.Errorf("%s: GrantedDominate = false with hero on, want true", c.Name())
		}
		if !self.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = false with hero on, want true", c.Name())
		}
	}
}
