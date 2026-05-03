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

// Pins the on-hit-go-again chain extension through a weapon swing. Hand: a pure-pitch Blue,
// Critical Strike Red, Razor Reflex Blue, Nimblism Blue. Weapon: Reaping Blade.
//
// Expected line:
//   - Pitch the Blue (3 resource)
//   - Play Nimblism Blue (cost 0, GoAgain, +1{p} to next cost-≤1 attack action)
//   - Play Critical Strike Red (cost 1, base 5, +1{p} from Nimblism → 6 incoming)
//   - Play Razor Reflex Blue (cost 1, AR mode 1: +1{p} → CS at 7, on-hit go-again grants AP)
//   - Swing Reaping Blade (cost 1, power 3) using the granted AP
//
// Expected Value: 7 (CS post-buff) + 3 (Reaping) = 10.
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

// Same shape as the test above, but the Snatch + 2x Razor Reflex variant exercises two
// CONSECUTIVE attack reactions stacking on a single attack target — instead of a pre-attack
// buff (Nimblism) plus one AR. Hand: a pure-pitch Blue, Snatch Blue, Razor Reflex Blue,
// Razor Reflex Blue. Weapon: Reaping Blade.
//
// Expected line:
//   - Pitch the Blue (3 resource)
//   - Play Snatch Blue (cost 0, base 2, OnHit DrawOne)
//   - Play Razor Reflex Blue #1 (AR mode 1: +1{p} → Snatch at 3)
//   - Play Razor Reflex Blue #2 (AR mode 1: +1{p} → Snatch at 4, lands in 1/4/7 window;
//     on-hit go-again grants AP for the weapon swing)
//   - Swing Reaping Blade (cost 1, power 3) using the granted AP
//
// Expected Value: 4 (Snatch post-buff: 2 + 1 + 1) + 3 (Reaping) = 7.
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

// formatBestLine renders the chosen role assignment in chain order so a failing test can
// show what line the optimizer actually picked.
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
