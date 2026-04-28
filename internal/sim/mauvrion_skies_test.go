package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// End-to-end coverage for Mauvrion Skies's EphemeralAttackTrigger. Each test plays a full
// turn via Best and asserts the resulting damage picks up (or drops) the rider's +N
// Runechant contribution. Running through the sim pins Mauvrion's two grants — the go-again
// look-ahead and the trigger-based Runechant rider — together, so a regression in either
// path shows up here.

// TestBest_MauvrionAloneFizzlesWithoutDamage: with Mauvrion the only attacker, no matching
// target ever resolves and the ephemeral trigger fizzles silently. Best picks the highest-
// value line, which is to simply play Mauvrion for zero damage.
func TestBest_MauvrionAloneFizzlesWithoutDamage(t *testing.T) {
	h := []Card{cards.MauvrionSkiesRed{}}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
	if got.Value != 0 {
		t.Fatalf("want value 0 (trigger has no target), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_MauvrionBladeOnlyFizzles: a Runeblade weapon swing after Mauvrion isn't an
// attack action card, so Mauvrion's Matches predicate rejects it and the trigger fizzles.
// Total damage is just the weapon's own contribution.
func TestBest_MauvrionBladeOnlyFizzles(t *testing.T) {
	h := []Card{cards.MauvrionSkiesRed{}, testutils.YellowAttack{}}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(StubHero, weapons, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → play Mauvrion (cost 0, go again) → Blade swing (cost 1,
	// 3 damage). Mauvrion's trigger doesn't match the weapon, so no Runechants.
	if got.Value != 3 {
		t.Fatalf("want value 3 (weapon swing only, Mauvrion fizzles), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_MauvrionNonRunebladeAttackFizzles: a Generic attack action card after Mauvrion
// is still an attack action but isn't a Runeblade attack, so Mauvrion's Matches predicate
// rejects it and the trigger fizzles.
func TestBest_MauvrionNonRunebladeAttackFizzles(t *testing.T) {
	h := []Card{cards.MauvrionSkiesRed{}, testutils.RedAttack{}, testutils.YellowAttack{}}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → play Mauvrion (cost 0, go again) → play fake RedAttack
	// (cost 1, 3 damage, go again). The Generic attack action doesn't qualify for
	// Mauvrion's rider, so Runechants never fire.
	if got.Value != 3 {
		t.Fatalf("want value 3 (Generic attack only), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_MauvrionLikelyHitRunebladeAttackCreditsRider: with a Runeblade attack action
// whose printed power lands in the likely-to-hit set, Mauvrion's trigger fires on the
// target's resolution and credits its +3 Runechants. Shrill Red has power 4 (in 1/4/7) and
// no printed go-again — Mauvrion's look-ahead grant is what lets the chain reach Shrill
// legally.
func TestBest_MauvrionLikelyHitRunebladeAttackCreditsRider(t *testing.T) {
	h := []Card{
		cards.MauvrionSkiesRed{},
		cards.ShrillOfSkullformRed{},
		testutils.YellowAttack{},
	}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → Mauvrion (cost 0, go again, grants go-again to Shrill +
	// registers trigger) → Shrill (cost 2, power 4). No aura exists when Shrill's Play
	// runs (StubHero has no trigger and Mauvrion's trigger hasn't fired yet), so Shrill's
	// own +3 "aura played" bonus stays off. After Shrill, Mauvrion's ephemeral fires:
	// LikelyToHit(4) is true, so it creates 3 Runechants (+3 damage, credited to Mauvrion).
	// Total: 4 (Shrill) + 3 (Mauvrion's rider) = 7.
	if got.Value != 7 {
		t.Fatalf("want value 7 (Shrill 4 + Mauvrion rider 3), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_MauvrionBlockableRunebladeAttackDropsRider: when the matching target's printed
// power is blockable, the trigger's Handler runs LikelyToHit and returns 0 — the rider
// doesn't fire, no Runechants are created, and the go-again grant still lands.
func TestBest_MauvrionBlockableRunebladeAttackDropsRider(t *testing.T) {
	h := []Card{
		cards.MauvrionSkiesRed{},
		cards.ShrillOfSkullformBlue{},
		testutils.YellowAttack{},
	}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
	// Pitch YellowAttack (2 res) → Mauvrion (cost 0) → Shrill Blue (cost 2, power 2).
	// Trigger fires on Shrill's resolution but LikelyToHit(2) is false, so Mauvrion's
	// Runechants don't land. Shrill's own +3 aura bonus also stays off (no auras).
	// Total: 2 (Shrill Blue).
	if got.Value != 2 {
		t.Fatalf("want value 2 (Shrill Blue only, Mauvrion rider drops), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}
