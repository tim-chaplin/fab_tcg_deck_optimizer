package card

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// fakeCard is a minimal Card with the fields these tests read. Inline (rather than from
// a shared testutils package) to keep internal/card free of test-only outside imports — the
// alternative would force an external package_test, which the project doesn't use.
type fakeCard struct {
	name    string
	attack  int
	goAgain bool
}

func (c fakeCard) ID() ids.CardID                    { return ids.InvalidCard }
func (c fakeCard) Name() string                      { return c.name }
func (c fakeCard) DisplayName() string               { return c.name }
func (fakeCard) Cost(GameEngine) int                 { return 0 }
func (fakeCard) Pitch() int                          { return 0 }
func (c fakeCard) Attack() int                       { return c.attack }
func (fakeCard) Defense() int                        { return 0 }
func (fakeCard) Types(GameEngine) TypeSet            { return 0 }
func (c fakeCard) GoAgain(GameEngine) bool           { return c.goAgain }
func (fakeCard) Play(GameEngine, Logger, *CardState) {}

// dominatingFake is a fakeCard with the Dominator marker — exercises the printed-
// Dominate branches of EffectiveDominate / HasDominate.
type dominatingFake struct {
	fakeCard
}

func (dominatingFake) Dominate() {}

// TestCardState_EffectiveGoAgain: printed GoAgain OR a mid-chain grant qualifies the card
// for Go again. Neither printed nor granted -> false.
func TestCardState_EffectiveGoAgain(t *testing.T) {
	cases := []struct {
		name    string
		printed bool
		granted bool
		want    bool
	}{
		{"neither", false, false, false},
		{"printed only", true, false, true},
		{"granted only", false, true, true},
		{"both", true, true, true},
	}
	for _, tc := range cases {
		p := &CardState{
			Card:    fakeCard{name: tc.name, goAgain: tc.printed},
			Ephemeral: Ephemeral{GrantedGoAgain: tc.granted},
		}
		if got := p.EffectiveGoAgain(nil); got != tc.want {
			t.Errorf("%s: EffectiveGoAgain() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestCardState_EffectiveDominate: the Dominator marker OR a mid-chain grant (a "gains
// dominate" rider flipping self.GrantedDominate) each qualifies the attack as dominating.
func TestCardState_EffectiveDominate(t *testing.T) {
	plain := fakeCard{name: "plain"}
	dominator := dominatingFake{fakeCard{name: "printed"}}

	cases := []struct {
		name    string
		card    Card
		granted bool
		want    bool
	}{
		{"neither", plain, false, false},
		{"printed only", dominator, false, true},
		{"granted only", plain, true, true},
		{"both", dominator, true, true},
	}
	for _, tc := range cases {
		p := &CardState{Card: tc.card, Ephemeral: Ephemeral{GrantedDominate: tc.granted}}
		if got := p.EffectiveDominate(); got != tc.want {
			t.Errorf("%s: EffectiveDominate() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestHasDominate_MatchesMarker: the free helper is the static printed-keyword check;
// type assertion to Dominator decides.
func TestHasDominate_MatchesMarker(t *testing.T) {
	if HasDominate(fakeCard{name: "plain"}) {
		t.Error("HasDominate(plain) = true, want false")
	}
	if !HasDominate(dominatingFake{}) {
		t.Error("HasDominate(dominator) = false, want true")
	}
}

// Tests EffectiveAttack: printed Attack + BonusAttack, clamped at 0 for negative bonuses.
func TestCardState_EffectiveAttack(t *testing.T) {
	cases := []struct {
		name        string
		printed     int
		bonusAttack int
		want        int
	}{
		{"no bonus", 4, 0, 4},
		{"granted +1 bumps 3 into the 1/4/7 window", 3, 1, 4},
		{"granted +3 stacks", 4, 3, 7},
		{"-2 on a 5-power attack", 5, -2, 3},
		{"-3 on a 3-power attack lands at exactly 0", 3, -3, 0},
		{"-2 on a 1-power attack clamps at 0 (can't go negative)", 1, -2, 0},
		{"large negative on a 4-power attack still clamps at 0", 4, -10, 0},
	}
	for _, tc := range cases {
		p := &CardState{
			Card:    fakeCard{name: tc.name, attack: tc.printed},
			Ephemeral: Ephemeral{BonusAttack: tc.bonusAttack},
		}
		if got := p.EffectiveAttack(); got != tc.want {
			t.Errorf("%s: EffectiveAttack() = %d, want %d", tc.name, got, tc.want)
		}
	}
}
