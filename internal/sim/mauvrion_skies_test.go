package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// End-to-end coverage for Mauvrion Skies's on-hit Runechant rider via Best. Pins both
// grants — the go-again look-ahead and the OnHit-registered Runechant rider — together.

// Tests that Mauvrion alone (no matching target) deals zero damage.
func TestBest_MauvrionAloneFizzlesWithoutDamage(t *testing.T) {
	h := []Card{cards.MauvrionSkiesRed{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, 0, nil, 0, nil)
	if got.Value != 0 {
		t.Fatalf("want value 0 (trigger has no target), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that a Runeblade weapon swing doesn't satisfy Mauvrion's predicate (attack action
// only).
func TestBest_MauvrionBladeOnlyFizzles(t *testing.T) {
	h := []Card{cards.MauvrionSkiesRed{}, testutils.YellowAttack{}}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(testutils.Hero{Intel: 4}, weapons, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → play Mauvrion (cost 0, go again) → Blade swing (cost 1,
	// 3 damage). Mauvrion's trigger doesn't match the weapon, so no Runechants.
	if got.Value != 3 {
		t.Fatalf("want value 3 (weapon swing only, Mauvrion fizzles), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that a Generic (non-Runeblade) attack action doesn't satisfy Mauvrion's predicate.
func TestBest_MauvrionNonRunebladeAttackFizzles(t *testing.T) {
	h := []Card{cards.MauvrionSkiesRed{}, testutils.RedAttack{}, testutils.YellowAttack{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → play Mauvrion (cost 0, go again) → play fake RedAttack
	// (cost 1, 3 damage, go again). The Generic attack action doesn't qualify for
	// Mauvrion's rider, so Runechants never fire.
	if got.Value != 3 {
		t.Fatalf("want value 3 (Generic attack only), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that a likely-hit Runeblade attack action picks up Mauvrion's +3 Runechant rider.
func TestBest_MauvrionLikelyHitRunebladeAttackCreditsRider(t *testing.T) {
	h := []Card{
		cards.MauvrionSkiesRed{},
		cards.ShrillOfSkullformRed{},
		testutils.YellowAttack{},
	}
	got := Best(testutils.Hero{Intel: 4}, nil, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → Mauvrion (cost 0, go again, grants go-again to Shrill +
	// appends OnHit) → Shrill (cost 2, power 4). No aura when Shrill's Play runs, so its
	// own +3 "aura played" bonus stays off. Shrill's OnHit fires: LikelyToHit(4) is true,
	// 3 Runechants created (+3, credited to Mauvrion). Total: 4 + 3 = 7.
	if got.Value != 7 {
		t.Fatalf("want value 7 (Shrill 4 + Mauvrion rider 3), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that a blockable Runeblade attack drops Mauvrion's Runechant rider but keeps the
// go-again grant.
func TestBest_MauvrionBlockableRunebladeAttackDropsRider(t *testing.T) {
	h := []Card{
		cards.MauvrionSkiesRed{},
		cards.ShrillOfSkullformBlue{},
		testutils.YellowAttack{},
	}
	got := Best(testutils.Hero{Intel: 4}, nil, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → Mauvrion (cost 0) → Shrill Blue (cost 2, power 2).
	// LikelyToHit(2) is false, so the OnHit doesn't fire. Total: 2.
	if got.Value != 2 {
		t.Fatalf("want value 2 (Shrill Blue only, Mauvrion rider drops), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}
