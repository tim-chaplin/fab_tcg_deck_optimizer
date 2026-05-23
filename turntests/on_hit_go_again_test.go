package turntests

import (
	"strings"
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

// Pins on-hit-go-again chain extension through a weapon swing: Nimblism + Razor Reflex
// buff Critical Strike to 7, RR's on-hit go-again grants AP for Reaping Blade (3).
func TestOnHitGoAgain_RazorReflexExtendsToWeaponSwing(t *testing.T) {
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, fillerDeck())
	hand := []card.Card{
		testutils.BluePitch{},
		cards.CriticalStrikeRed{},
		cards.RazorReflexBlue{},
		cards.NimblismBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 10 {
		t.Fatalf("Value = %d, want 10 (CS buffed to 7 by Nimblism+RR + Reaping 3 via on-hit go-again)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}

// Pins two consecutive attack reactions stacking on a single attack: 2x Razor Reflex
// buff Snatch to 4, the second RR's on-hit go-again grants AP for Reaping Blade (3).
func TestOnHitGoAgain_TwoConsecutiveARsExtendToWeaponSwing(t *testing.T) {
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.ReapingBlade{}}, fillerDeck())
	hand := []card.Card{
		testutils.BluePitch{},
		cards.SnatchBlue{},
		cards.RazorReflexBlue{},
		cards.RazorReflexBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 7 {
		t.Fatalf("Value = %d, want 7 (Snatch buffed to 4 by 2x RR + Reaping 3 via on-hit go-again)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}

// formatBestLine renders the chosen role assignment in chain order for failure messages.
func formatBestLine(line []card.CardAssignment) string {
	var parts []string
	for _, a := range line {
		parts = append(parts, a.Card.DisplayName()+":"+roleName(a.Role))
	}
	return strings.Join(parts, ", ")
}

func roleName(r card.Role) string {
	switch r {
	case card.Pitch:
		return "Pitch"
	case card.Attack:
		return "Attack"
	case card.Defend:
		return "Defend"
	case card.Held:
		return "Held"
	case card.Arsenal:
		return "Arsenal"
	}
	return "Unknown"
}
