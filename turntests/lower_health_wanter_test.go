package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// fakeLowHeroOn implements gameengine.LowerHealthWanter — used to exercise the "hero opts
// in" branch.
type fakeLowHeroOn struct{}

func (fakeLowHeroOn) ID() ids.HeroID                                                       { return ids.InvalidHero }
func (fakeLowHeroOn) Name() string                                                         { return "fakeLowHeroOn" }
func (fakeLowHeroOn) Intelligence() int                                                    { return 4 }
func (fakeLowHeroOn) Types() card.TypeSet                                                  { return 0 }
func (fakeLowHeroOn) Class() card.CardType                                                 { return 0 }
func (fakeLowHeroOn) Opt(cards []card.Card) (top, bottom []card.Card)                      { return cards, nil }
func (fakeLowHeroOn) WantsLowerHealth()                                                    {}
func (fakeLowHeroOn) TriggerType() triggertype.Type                                        { return 0 }
func (fakeLowHeroOn) OncePerTurn() bool                                                    { return false }
func (fakeLowHeroOn) FiredThisTurn() bool                                                  { return false }
func (fakeLowHeroOn) SetFiredThisTurn(bool)                                                {}
func (fakeLowHeroOn) Matches(card.TypeSet) bool                                            { return false }
func (fakeLowHeroOn) Fire(card.GameEngine, card.Logger, *card.CardState, triggertype.Type) {}

// fakeLowHeroOff does NOT implement gameengine.LowerHealthWanter — the default branch.
type fakeLowHeroOff struct{}

func (fakeLowHeroOff) ID() ids.HeroID                                                       { return ids.InvalidHero }
func (fakeLowHeroOff) Name() string                                                         { return "fakeLowHeroOff" }
func (fakeLowHeroOff) Intelligence() int                                                    { return 4 }
func (fakeLowHeroOff) Types() card.TypeSet                                                  { return 0 }
func (fakeLowHeroOff) Class() card.CardType                                                 { return 0 }
func (fakeLowHeroOff) Opt(cards []card.Card) (top, bottom []card.Card)                      { return cards, nil }
func (fakeLowHeroOff) TriggerType() triggertype.Type                                        { return 0 }
func (fakeLowHeroOff) OncePerTurn() bool                                                    { return false }
func (fakeLowHeroOff) FiredThisTurn() bool                                                  { return false }
func (fakeLowHeroOff) SetFiredThisTurn(bool)                                                {}
func (fakeLowHeroOff) Matches(card.TypeSet) bool                                            { return false }
func (fakeLowHeroOff) Fire(card.GameEngine, card.Logger, *card.CardState, triggertype.Type) {}

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
		sOff := engineWithHero(fakeLowHeroOff{})
		sOff.ResolveAttackStep(sOff.Logger(), &card.CardState{Card: tc.card})
		if got := sOff.Value(); got != tc.wantOff {
			t.Errorf("%s: Play() off = %d, want %d (hero does not opt in)", tc.name, got, tc.wantOff)
		}
		sOn := engineWithHero(fakeLowHeroOn{})
		sOn.ResolveAttackStep(sOn.Logger(), &card.CardState{Card: tc.card})
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
	gOff := engineWithHero(fakeLowHeroOff{})
	for _, c := range cards {
		if c.GoAgain(gOff) {
			t.Errorf("%s: GoAgain() = true with hero off, want false", c.Name())
		}
	}
	gOn := engineWithHero(fakeLowHeroOn{})
	for _, c := range cards {
		if !c.GoAgain(gOn) {
			t.Errorf("%s: GoAgain() = false with hero on, want true", c.Name())
		}
	}
}

// TestLowerHealthWanter_NilHeroIsOff guards the startup / unset-hero case.
func TestLowerHealthWanter_NilHeroIsOff(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.AdrenalineRushRed{}})
	if got := ge.Value(); got != 4 {
		t.Errorf("AdrenalineRushRed nil-hero Play() = %d, want 4", got)
	}
	if (cards.ScarForAScarRed{}).GoAgain(ge) {
		t.Errorf("ScarForAScarRed nil-hero GoAgain() = true, want false")
	}
}

// Tests that Pound for Pound flips pc.GrantedDominate iff the current hero opts into
// LowerHealthWanter.
func TestLowerHealthWanter_PoundForPoundDominateGrant(t *testing.T) {
	cards := []card.Card{cards.PoundForPoundRed{}, cards.PoundForPoundYellow{}, cards.PoundForPoundBlue{}}

	for _, c := range cards {
		pc := &card.CardState{Card: c}
		s := engineWithHero(fakeLowHeroOff{})
		s.ResolveAttackStep(s.Logger(), pc)
		if pc.GrantedDominate {
			t.Errorf("%s: GrantedDominate = true with hero off, want false", c.Name())
		}
		if pc.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = true with hero off, want false", c.Name())
		}
	}
	for _, c := range cards {
		pc := &card.CardState{Card: c}
		s := engineWithHero(fakeLowHeroOn{})
		s.ResolveAttackStep(s.Logger(), pc)
		if !pc.GrantedDominate {
			t.Errorf("%s: GrantedDominate = false with hero on, want true", c.Name())
		}
		if !pc.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = false with hero on, want true", c.Name())
		}
	}
}
