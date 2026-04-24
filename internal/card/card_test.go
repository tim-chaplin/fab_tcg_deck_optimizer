package card

import "testing"

// TestNotImplementedMarker pins the type-assertion contract: a plain Card does NOT satisfy the
// NotImplemented interface, and a Card whose type carries a NotImplemented() method does.
// That's the exact check the deck legal-pool filter performs when deciding whether to skip a
// card in random generation or mutation pools.
func TestNotImplementedMarker(t *testing.T) {
	var plain Card = stubCard{name: "plain"}
	if _, ok := plain.(NotImplemented); ok {
		t.Error("plain stub satisfied NotImplemented — the marker must be opt-in, not implicit")
	}
	var tagged Card = notImplementedStubCard{stubCard{name: "tagged"}}
	if _, ok := tagged.(NotImplemented); !ok {
		t.Error("tagged stub failed NotImplemented assertion — defining NotImplemented() must opt in")
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
		p := &CardState{Card: stubCard{name: tc.name, goAgain: tc.printed}, GrantedGoAgain: tc.granted}
		if got := p.EffectiveGoAgain(); got != tc.want {
			t.Errorf("%s: EffectiveGoAgain() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestCardState_EffectiveDominate: the Dominator marker OR a mid-chain grant (a "gains
// dominate" rider flipping self.GrantedDominate) each qualifies the attack as dominating.
func TestCardState_EffectiveDominate(t *testing.T) {
	plain := stubCard{name: "plain"}
	dominator := dominatingStubCard{stubCard: stubCard{name: "printed"}}

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
	if HasDominate(stubCard{name: "plain"}) {
		t.Error("HasDominate(plain) = true, want false")
	}
	if !HasDominate(dominatingStubCard{}) {
		t.Error("HasDominate(dominator) = false, want true")
	}
}
