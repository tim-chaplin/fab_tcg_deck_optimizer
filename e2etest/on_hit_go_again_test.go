package e2etest

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Pins on-hit-go-again chain extension through a weapon swing: Nimblism + Razor Reflex
// buff Critical Strike to 7, RR's on-hit go-again grants AP for Reaping Blade (3).
func TestOnHitGoAgain_RazorReflexExtendsToWeaponSwing(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, fillerDeck())
	hand := []sim.Card{
		testutils.BluePitch{},
		cards.CriticalStrikeRed{},
		cards.RazorReflexBlue{},
		cards.NimblismBlue{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand)
	if state.Value != 10 {
		t.Fatalf("Value = %d, want 10 (CS buffed to 7 by Nimblism+RR + Reaping 3 via on-hit go-again)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
}

// Pins two consecutive attack reactions stacking on a single attack: 2x Razor Reflex
// buff Snatch to 4, the second RR's on-hit go-again grants AP for Reaping Blade (3).
func TestOnHitGoAgain_TwoConsecutiveARsExtendToWeaponSwing(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, fillerDeck())
	hand := []sim.Card{
		testutils.BluePitch{},
		cards.SnatchBlue{},
		cards.RazorReflexBlue{},
		cards.RazorReflexBlue{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand)
	if state.Value != 7 {
		t.Fatalf("Value = %d, want 7 (Snatch buffed to 4 by 2x RR + Reaping 3 via on-hit go-again)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
}

// formatBestLine renders the chosen role assignment in chain order for failure messages.
func formatBestLine(line []sim.CardAssignment) string {
	var parts []string
	for _, a := range line {
		parts = append(parts, sim.DisplayName(a.Card)+":"+roleName(a.Role))
	}
	return strings.Join(parts, ", ")
}

func roleName(r sim.Role) string {
	switch r {
	case sim.Pitch:
		return "Pitch"
	case sim.Attack:
		return "Attack"
	case sim.Defend:
		return "Defend"
	case sim.Held:
		return "Held"
	case sim.Arsenal:
		return "Arsenal"
	}
	return "Unknown"
}
