package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests that Captain's Call picks the go-again mode when a follow-up attack can extend the
// chain into more total damage than the +2{p} buff alone.
func TestModal_CaptainsCallPicksGoAgainOverBuffWhenChainExtends(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.CaptainsCallRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (mode 1 grants go-again so both Snatches chain)", got)
	}
}

// Tests that Captain's Call picks the +2{p} mode when no follow-up attack can use a granted
// go-again.
func TestModal_CaptainsCallPicksBuffWhenChainCantExtend(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.CaptainsCallRed{},
		cards.SnatchRed{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 6 {
		t.Fatalf("Value = %d, want 6 (mode 0 +2{p} since no second attack to extend into)", got)
	}
}

// Tests that Razor Reflex's mode-0 +N{p} buff lands on a sword weapon target.
func TestModal_RazorReflexMode0BuffsSwordWeapon(t *testing.T) {
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.NebulaBlade{}}, nil)
	hand := []card.Card{
		cards.RazorReflexRed{},
		cards.ToughenUpBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (NebulaBlade 1 + Razor Reflex mode 0 +3 + runechant 1)", got)
	}
}

// Tests that Razor Reflex mode 1's +N{p} buff plus on-hit go-again rider both land on a
// cost-≤1 attack action: the buffed Snatch hits 7 power (in the 1/4/7 likely-hit window),
// the eager on-hit go-again grants 1 AP, and a second Snatch chains for full damage.
func TestModal_RazorReflexMode1BuffAndOnHitGoAgainExtendChain(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.RazorReflexRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
		testutils.FakeBlueAttack(),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 11 {
		t.Fatalf("Value = %d, want 11 (Snatch1 4 + Razor Reflex +3 + Snatch2 4 via on-hit go-again)", got)
	}
}

// Tests that Pummel's mode-1 +N{p} buff and on-hit hero-discard rider both land on a cost-≥2
// attack action target.
func TestModal_PummelMode1BuffsAndDiscardsOnHit(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.PummelBlue{},
		cards.AdrenalineRushBlue{},
		testutils.FakeBlueAttack(),
		testutils.FakeBlueAttack(),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 7 {
		t.Fatalf("Value = %d, want 7 (AdrenalineRush 2 + Pummel +2 + on-hit discard 3)", got)
	}
}

// Tests that Pummel's mode-0 +N{p} buff lands on a Club weapon target.
func TestModal_PummelMode0BuffsClubWeapon(t *testing.T) {
	d := deck.New(heroes.Viserai, []deck.Weapon{testutils.ClubWeapon{}}, nil)
	hand := []card.Card{
		cards.PummelRed{},
		cards.ToughenUpBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	got := summary.Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (Club 1 + Pummel mode 0 +4)", got)
	}
}
