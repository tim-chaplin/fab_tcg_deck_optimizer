package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// stubLowHeroOn implements sim.LowerHealthWanter — used to exercise the "hero opts in" branch.
type stubLowHeroOn struct{}

func (stubLowHeroOn) ID() ids.HeroID                            { return ids.InvalidHero }
func (stubLowHeroOn) Name() string                              { return "stubLowHeroOn" }
func (stubLowHeroOn) Health() int                               { return 20 }
func (stubLowHeroOn) Intelligence() int                         { return 4 }
func (stubLowHeroOn) Types() card.TypeSet                       { return 0 }
func (stubLowHeroOn) OnCardPlayed(sim.Card, *sim.TurnState) int { return 0 }
func (stubLowHeroOn) WantsLowerHealth()                         {}

// stubLowHeroOff does NOT implement sim.LowerHealthWanter — the default-hero branch.
type stubLowHeroOff struct{}

func (stubLowHeroOff) ID() ids.HeroID                            { return ids.InvalidHero }
func (stubLowHeroOff) Name() string                              { return "stubLowHeroOff" }
func (stubLowHeroOff) Health() int                               { return 20 }
func (stubLowHeroOff) Intelligence() int                         { return 4 }
func (stubLowHeroOff) Types() card.TypeSet                       { return 0 }
func (stubLowHeroOff) OnCardPlayed(sim.Card, *sim.TurnState) int { return 0 }

// TestLowerHealthWanter_DamageRiders checks the +3{p} / +1{p} / +1{h} damage riders fire iff the
// current hero opts into sim.LowerHealthWanter.
func TestLowerHealthWanter_DamageRiders(t *testing.T) {
	cases := []struct {
		name    string
		card    sim.Card
		wantOff int
		wantOn  int
	}{
		{"AdrenalineRushRed +3p", AdrenalineRushRed{}, 4, 4 + 3},
		{"AdrenalineRushYellow +3p", AdrenalineRushYellow{}, 3, 3 + 3},
		{"AdrenalineRushBlue +3p", AdrenalineRushBlue{}, 2, 2 + 3},
		{"WoundedBullRed +1p", WoundedBullRed{}, 7, 7 + 1},
		{"WoundedBullYellow +1p", WoundedBullYellow{}, 6, 6 + 1},
		{"WoundedBullBlue +1p", WoundedBullBlue{}, 5, 5 + 1},
		{"FyendalsFightingSpiritRed +1h", FyendalsFightingSpiritRed{}, 7, 7 + 1},
		{"FyendalsFightingSpiritYellow +1h", FyendalsFightingSpiritYellow{}, 6, 6 + 1},
		{"FyendalsFightingSpiritBlue +1h", FyendalsFightingSpiritBlue{}, 5, 5 + 1},
	}
	for _, tc := range cases {
		sim.CurrentHero = stubLowHeroOff{}
		var sOff sim.TurnState
		tc.card.Play(&sOff, &sim.CardState{Card: tc.card})
		if got := sOff.Value; got != tc.wantOff {
			t.Errorf("%s: Play() off = %d, want %d (hero does not opt in)", tc.name, got, tc.wantOff)
		}
		sim.CurrentHero = stubLowHeroOn{}
		var sOn sim.TurnState
		tc.card.Play(&sOn, &sim.CardState{Card: tc.card})
		if got := sOn.Value; got != tc.wantOn {
			t.Errorf("%s: Play() on = %d, want %d (hero opts in)", tc.name, got, tc.wantOn)
		}
	}
	sim.CurrentHero = nil
}

// TestLowerHealthWanter_GoAgainRiders checks the conditional go-again flips iff the current hero
// opts into sim.LowerHealthWanter.
func TestLowerHealthWanter_GoAgainRiders(t *testing.T) {
	cards := []sim.Card{
		ScarForAScarRed{}, ScarForAScarYellow{}, ScarForAScarBlue{},
		BlowForABlowRed{},
		LifeForALifeRed{}, LifeForALifeYellow{}, LifeForALifeBlue{},
	}
	sim.CurrentHero = stubLowHeroOff{}
	for _, c := range cards {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true with hero off, want false", c.Name())
		}
	}
	sim.CurrentHero = stubLowHeroOn{}
	for _, c := range cards {
		if !c.GoAgain() {
			t.Errorf("%s: GoAgain() = false with hero on, want true", c.Name())
		}
	}
	sim.CurrentHero = nil
}

// TestLowerHealthWanter_NilHeroIsOff guards the startup / unset-hero case: with no hero, the rider
// must not fire.
func TestLowerHealthWanter_NilHeroIsOff(t *testing.T) {
	sim.CurrentHero = nil
	var s sim.TurnState
	(AdrenalineRushRed{}).Play(&s, &sim.CardState{Card: AdrenalineRushRed{}})
	if got := s.Value; got != 4 {
		t.Errorf("AdrenalineRushRed nil-hero Play() = %d, want 4", got)
	}
	if (ScarForAScarRed{}).GoAgain() {
		t.Errorf("ScarForAScarRed nil-hero GoAgain() = true, want false")
	}
}

// TestLowerHealthWanter_PoundForPoundDominateGrant: Pound for Pound's conditional Dominate
// fires via self.GrantedDominate iff the current hero opts into LowerHealthWanter. Damage
// itself is unchanged — the grant feeds EffectiveDominate for downstream scanners / future
// on-hit riders.
func TestLowerHealthWanter_PoundForPoundDominateGrant(t *testing.T) {
	cards := []sim.Card{PoundForPoundRed{}, PoundForPoundYellow{}, PoundForPoundBlue{}}

	sim.CurrentHero = stubLowHeroOff{}
	for _, c := range cards {
		self := &sim.CardState{Card: c}
		c.Play(&sim.TurnState{}, self)
		if self.GrantedDominate {
			t.Errorf("%s: GrantedDominate = true with hero off, want false", c.Name())
		}
		if self.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = true with hero off, want false", c.Name())
		}
	}

	sim.CurrentHero = stubLowHeroOn{}
	for _, c := range cards {
		self := &sim.CardState{Card: c}
		c.Play(&sim.TurnState{}, self)
		if !self.GrantedDominate {
			t.Errorf("%s: GrantedDominate = false with hero on, want true", c.Name())
		}
		if !self.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = false with hero on, want true", c.Name())
		}
	}
	sim.CurrentHero = nil
}
