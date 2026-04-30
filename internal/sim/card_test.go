package sim_test

import (
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestNotImplementedMarker pins the type-assertion contract: a plain Card does NOT satisfy the
// NotImplemented interface, and a Card whose type carries a NotImplemented() method does.
// That's the exact check the deck legal-pool filter performs when deciding whether to skip a
// card in random generation or mutation pools.
func TestNotImplementedMarker(t *testing.T) {
	var plain Card = testutils.NewStubCard("plain")
	if _, ok := plain.(NotImplemented); ok {
		t.Error("plain stub satisfied NotImplemented — the marker must be opt-in, not implicit")
	}
	var tagged Card = testutils.NotImplementedStubCard{StubCard: testutils.NewStubCard("tagged")}
	if _, ok := tagged.(NotImplemented); !ok {
		t.Error("tagged stub failed NotImplemented assertion — defining NotImplemented() must opt in")
	}
}

// TestUnplayableMarker mirrors TestNotImplementedMarker for the Unplayable interface — opt-in,
// not implicit, and orthogonal to NotImplemented (a plain stub satisfies neither, an
// Unplayable stub satisfies Unplayable but not NotImplemented).
func TestUnplayableMarker(t *testing.T) {
	var plain Card = testutils.NewStubCard("plain")
	if _, ok := plain.(Unplayable); ok {
		t.Error("plain stub satisfied Unplayable — the marker must be opt-in, not implicit")
	}
	var tagged Card = testutils.UnplayableStubCard{StubCard: testutils.NewStubCard("tagged")}
	if _, ok := tagged.(Unplayable); !ok {
		t.Error("tagged stub failed Unplayable assertion — defining Unplayable() must opt in")
	}
	if _, ok := tagged.(NotImplemented); ok {
		t.Error("Unplayable stub also satisfied NotImplemented — the markers must be orthogonal")
	}
}

// TestIsExcludedFromPool_BothMarkers pins the centralised pool-exclusion check: NotImplemented
// OR Unplayable trips the filter, plain stubs don't.
func TestIsExcludedFromPool_BothMarkers(t *testing.T) {
	cases := []struct {
		name string
		card Card
		want bool
	}{
		{"plain", testutils.NewStubCard("plain"), false},
		{"NotImplemented", testutils.NotImplementedStubCard{StubCard: testutils.NewStubCard("ni")}, true},
		{"Unplayable", testutils.UnplayableStubCard{StubCard: testutils.NewStubCard("up")}, true},
	}
	for _, tc := range cases {
		if got := IsExcludedFromPool(tc.card); got != tc.want {
			t.Errorf("%s: IsExcludedFromPool = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestCardState_EffectiveGoAgain: printed GoAgain OR a mid-chain grant (Mauvrion Skies et al)
// each qualifies the card for Go again. Neither printed nor granted → false.
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
		base := testutils.NewStubCard(tc.name)
		if tc.printed {
			base = base.WithGoAgain()
		}
		p := &CardState{Card: base, GrantedGoAgain: tc.granted}
		if got := p.EffectiveGoAgain(); got != tc.want {
			t.Errorf("%s: EffectiveGoAgain() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestCardState_EffectiveDominate: the Dominator marker OR a mid-chain grant (a "gains
// dominate" rider flipping self.GrantedDominate) each qualifies the attack as dominating.
func TestCardState_EffectiveDominate(t *testing.T) {
	plain := testutils.NewStubCard("plain")
	dominator := testutils.DominatingStubCard{StubCard: testutils.NewStubCard("printed")}

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
		p := &CardState{Card: tc.card, GrantedDominate: tc.granted}
		if got := p.EffectiveDominate(); got != tc.want {
			t.Errorf("%s: EffectiveDominate() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestHasDominate_MatchesMarker: the free helper is the static printed-keyword check;
// type assertion to Dominator decides.
func TestHasDominate_MatchesMarker(t *testing.T) {
	if HasDominate(testutils.NewStubCard("plain")) {
		t.Error("HasDominate(plain) = true, want false")
	}
	if !HasDominate(testutils.DominatingStubCard{}) {
		t.Error("HasDominate(dominator) = false, want true")
	}
}

// TestCardState_EffectiveAttack: printed Attack plus any granted BonusAttack from a prior
// card's "next attack +N{p}" rider. Default BonusAttack of 0 leaves EffectiveAttack equal to
// the printed power. Negative bonuses (defender-side -N{p} debuffs like Drag Down) clamp at
// 0 — an attack's power can't be reduced below 0 in FaB.
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
			Card:        testutils.NewStubCard(tc.name).WithAttack(tc.printed),
			BonusAttack: tc.bonusAttack,
		}
		if got := p.EffectiveAttack(); got != tc.want {
			t.Errorf("%s: EffectiveAttack() = %d, want %d", tc.name, got, tc.want)
		}
	}
}
