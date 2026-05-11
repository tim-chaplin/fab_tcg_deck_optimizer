package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// TestCardState_EffectiveGoAgain: printed GoAgain OR a mid-chain grant qualifies the card
// for Go again. Neither printed nor granted → false.
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
		base := NewFakeCard(tc.name)
		if tc.printed {
			base = base.WithGoAgain()
		}
		p := &card.CardState{Card: base, GrantedGoAgain: tc.granted}
		if got := p.EffectiveGoAgain(); got != tc.want {
			t.Errorf("%s: EffectiveGoAgain() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestCardState_EffectiveDominate: the Dominator marker OR a mid-chain grant (a "gains
// dominate" rider flipping self.GrantedDominate) each qualifies the attack as dominating.
func TestCardState_EffectiveDominate(t *testing.T) {
	plain := NewFakeCard("plain")
	dominator := DominatingFakeCard{FakeCard: NewFakeCard("printed")}

	cases := []struct {
		name    string
		card    card.Card
		granted bool
		want    bool
	}{
		{"neither", plain, false, false},
		{"printed only", dominator, false, true},
		{"granted only", plain, true, true},
		{"both", dominator, true, true},
	}
	for _, tc := range cases {
		p := &card.CardState{Card: tc.card, GrantedDominate: tc.granted}
		if got := p.EffectiveDominate(); got != tc.want {
			t.Errorf("%s: EffectiveDominate() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestHasDominate_MatchesMarker: the free helper is the static printed-keyword check;
// type assertion to Dominator decides.
func TestHasDominate_MatchesMarker(t *testing.T) {
	if card.HasDominate(NewFakeCard("plain")) {
		t.Error("HasDominate(plain) = true, want false")
	}
	if !card.HasDominate(DominatingFakeCard{}) {
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
		p := &card.CardState{
			Card:        NewFakeCard(tc.name).WithAttack(tc.printed),
			BonusAttack: tc.bonusAttack,
		}
		if got := p.EffectiveAttack(); got != tc.want {
			t.Errorf("%s: EffectiveAttack() = %d, want %d", tc.name, got, tc.want)
		}
	}
}
