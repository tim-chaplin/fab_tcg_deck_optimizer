package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests that Captain's Call picks the go-again mode when a follow-up attack can extend the
// chain into more total damage than the +2{p} buff alone.
func TestModal_CaptainsCallPicksGoAgainOverBuffWhenChainExtends(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.CaptainsCallRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
	}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 8 {
		t.Fatalf("Value = %d, want 8 (mode 1 grants go-again so both Snatches chain)", got)
	}
}

// Tests that Captain's Call picks the +2{p} mode when no follow-up attack can use a granted
// go-again.
func TestModal_CaptainsCallPicksBuffWhenChainCantExtend(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.CaptainsCallRed{},
		cards.SnatchRed{},
	}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 6 {
		t.Fatalf("Value = %d, want 6 (mode 0 +2{p} since no second attack to extend into)", got)
	}
}

// Tests that Razor Reflex's mode-0 +N{p} buff lands on a sword weapon target.
func TestModal_RazorReflexMode0BuffsSwordWeapon(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.NebulaBlade{}}, fillerDeck())
	hand := []sim.Card{
		cards.RazorReflexRed{},
		cards.ToughenUpBlue{},
	}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (NebulaBlade 1 + Razor Reflex mode 0 +3 + runechant 1)", got)
	}
}

// Tests that Razor Reflex's mode-1 +N{p} buff lands on a cost-≤1 attack action target.
func TestModal_RazorReflexMode1BuffsCostLowAttackAction(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.RazorReflexRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
	}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 7 {
		t.Fatalf("Value = %d, want 7 (Snatch 4 + Razor Reflex mode 1 +3, second Snatch pitches for cost)", got)
	}
}

// Tests that Pummel's mode-1 +N{p} buff and on-hit hero-discard rider both land on a cost-≥2
// attack action target.
func TestModal_PummelMode1BuffsAndDiscardsOnHit(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.PummelBlue{},
		cards.AdrenalineRushBlue{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 7 {
		t.Fatalf("Value = %d, want 7 (AdrenalineRush 2 + Pummel +2 + on-hit discard 3)", got)
	}
}
