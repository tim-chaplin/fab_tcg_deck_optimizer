package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

func TestBest_ViseraiMaleficShrillCombo(t *testing.T) {
	// Hero = Viserai. Best line: pitch the Blue Malefic, then play both Red Maleficas and the
	// Red Shrill. Each Malefic Play credits 0; Shrill (an attack action) fires both Malefics'
	// AttackAction triggers for +2 runes (OncePerTurn each, both at Count=3 so neither hits
	// zero). Plus Viserai's runechant on Shrill, plus Shrill's 4+3 aura-bonus = 11. Future
	// turns will keep ticking each Malefic for two more runes apiece, but those don't show
	// up in this turn's Value.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationRed{},
		runeblade.MaleficIncantationRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 4, nil, 0, nil)
	if got.Value != 11 {
		t.Fatalf("want value 11, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiReapingBladeBlueMalefics(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), play the other 3 Blue Malefics (Viserai runechants on #2
	// and #3), then swing Reaping Blade (cost 1, 3 dmg). Malefic's AttackAction triggers
	// don't fire here — the only attack is the weapon swing, which isn't an attack ACTION
	// card. Value = 0 + 1 + 1 + 3 = 5. The 3 Malefic verse counters carry forward and pay
	// out one rune apiece on future turns when an attack action lands.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 5 {
		t.Fatalf("want value 5, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiReapingBladeMaleficsPlusShrill(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), play 2 Blue Malefics (2 dmg + 1 Runechant), then Red Shrill
	// (cost 2, 4+3 aura bonus + 1 Runechant = 8). Reaping Blade stays holstered — Shrill has no
	// Go again, so nothing can follow it. Value = 2 + 1 + 8 = 11.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 11 {
		t.Fatalf("want value 11, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiOathBlueHocusRedMalefic(t *testing.T) {
	// Pitch Blue Hocus Pocus (3 res). Play Red Malefic (0 dmg, registers AttackAction
	// trigger). Play Red Oath (+1 Runechant, peeks ahead and sees the Blade swing = +3
	// bonus, +1 Viserai Runechant from prior non-attack action = 5). Swing Reaping Blade
	// (cost 1, 3 dmg) — Malefic's trigger doesn't fire (weapon swings aren't attack ACTION
	// cards). Value = 5 + 3 = 8. Future turns will tick Malefic when an attack action lands.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.OathOfTheArknightRed{},
		runeblade.MaleficIncantationRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 8 {
		t.Fatalf("want value 8, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_RunicReapingPrefersAttackPitch(t *testing.T) {
	// Pitching the Blue Hocus Pocus (attack-typed, pitch 3) pays for Runic Reaping + Shrill
	// AND satisfies Runic Reaping's pitched-attack rider, granting +1 to Shrill via
	// BonusAttack. Runic Reaping's "if this hits, create N Runechants" clause is registered
	// as an EphemeralAttackTrigger (same shape as Mauvrion Skies) and fires after Shrill's
	// full resolution: target.EffectiveAttack = printed 4 + BonusAttack 1 = 5, which falls
	// OUT of the {1,4,7} LikelyToHit window, so the runechant rider drops. The only damage
	// on Runic Reaping's slot is Viserai's runechant for the prior non-attack action.
	// Blue Malefic (1 arcane + 1 Viserai runechant = 2) → Runic Reaping (0 own damage + 1
	// Viserai runechant = 1) → Shrill (4 base + 3 aura-created bonus + 1 BonusAttack = 8).
	// Value = 2 + 1 + 8 = 11.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.RunicReapingRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, nil)
	if got.Value != 11 {
		t.Fatalf("want value 11, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiMauvrionGrantsGoAgainToShrill(t *testing.T) {
	// Pitch Blue Hocus Pocus (3 res). Play Blue Malefic (1 arcane, go again). Play Red Mauvrion
	// Skies (0 cost, go again; grants go-again to the next Runeblade attack action card = Shrill,
	// and emits 3 runechants). Play Red Shrill (cost 2, 4 base + 3 aura-created bonus = 7; chains
	// thanks to Mauvrion's grant). Swing Reaping Blade (cost 1, 3 dmg). Viserai fires +1 on
	// Mauvrion (prior Malefic is a non-attack action) and +1 on Shrill (priors include non-attack
	// actions). Value = 1 + 3 + 7 + 3 + 2 = 16.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MauvrionSkiesRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 16 {
		t.Fatalf("want value 16, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_ViseraiMauvrionPredictsDrowningDireDominate pins the full resolution stack each
// card's grants need to settle in the right order:
//   (1) target plays → hero ability fires (Viserai creates a Runechant → AuraCreated=true)
//   (2) target's Dominate clause resolves (Drowning Dire sees the aura → gains Dominate)
//   (3) Mauvrion's OnAttack trigger fires against the fully-resolved target state and
//       credits its Runechant rider iff the attack is now likely to hit.
// Requires hero.OnCardPlayed running before card.Play (so DD's aura check sees Viserai's
// Runechant), DD's conditional Dominate grant, and Mauvrion's ephemeral trigger reading
// target.EffectiveDominate() at fire time. If any of those regresses, this test drops to 6
// (no Mauvrion rider) or less.
//
// Line: pitch YellowAttack (2 res), play Mauvrion Red (cost 0, grants go-again + "if hits,
// create 3 Runechants" to the next Runeblade attack action), play Drowning Dire Red (cost 2,
// attack 5). Viserai creates a Runechant → Drowning Dire has Dominate → 5+ dominating
// attack lands → Mauvrion's rider fires for 3 Runechants. Value = 3 + 5 + 1 = 9.
func TestBest_ViseraiMauvrionPredictsDrowningDireDominate(t *testing.T) {
	h := []card.Card{
		runeblade.MauvrionSkiesRed{},
		runeblade.DrowningDireRed{},
		fake.YellowAttack{},
	}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, nil)
	if got.Value != 9 {
		t.Fatalf("want value 9, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestIsLegalOrder_MauvrionCantSaveShrillWhenRuneragerIsAhead(t *testing.T) {
	// Mauvrion's grant lands on the first matching Runeblade attack action card in CardsRemaining.
	// In the ordering Mauvrion → Runerager → Shrill → weapon, Runerager is that first match, so
	// Shrill never gets the grant. Shrill has no printed go-again, so the Shrill → weapon chain
	// must break — isLegalOrder rejects the ordering.
	order := []card.Card{
		runeblade.MauvrionSkiesRed{},
		runeblade.RuneragerSwarmRed{},
		runeblade.ShrillOfSkullformRed{},
		weapon.ReapingBlade{},
	}
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.playSequence(order, nil, nil, nil); legal {
		t.Fatalf("ordering %v should be illegal (Shrill has no go-again and Mauvrion granted Runerager instead)",
			cardNames(order))
	}
}

func TestBest_ViseraiMauvrionChainsShrillIntoRuneragerIntoWeapon(t *testing.T) {
	// Pitch Blue Hocus → Mauvrion → Shrill → Runerager → Reaping Blade. Value = 3 + 7 + 3 + 3 + 2
	// Viserai runechants = 18.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MauvrionSkiesRed{},
		runeblade.RuneragerSwarmRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 18 {
		t.Fatalf("want value 18, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBestSequence_CardStateGrantsDontLeakAcrossPermutations pins the per-permutation reset
// contract: the permutation loop in bestSequence must allocate fresh *CardState wrappers per
// permutation so a grant applied by one permutation's Play() can't bleed into a later
// permutation's legality/effect checks.
//
// Setup: attackers = [grantAll, grantSpy, grantAll]. The permute order emits grantAll-first
// permutations before the first grantSpy-first permutation. Each grantAll-first permutation
// mutates GrantedGoAgain on the wrappers for the cards behind it. When grantSpy later plays
// FIRST in its own permutation, its CardsRemaining contains the other two cards' wrappers —
// which must be fresh (GrantedGoAgain=false), since no card has played yet in this permutation.
// If the wrappers were reused across permutations the spy would see leaked grants and trip.
func TestBestSequence_CardStateGrantsDontLeakAcrossPermutations(t *testing.T) {
	var sawLeak bool
	attackers := []card.Card{grantAll{}, grantSpy{saw: &sawLeak}, grantAll{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 1_000_000, 0, len(attackers))
	_, _, _ = ctx.bestSequence(attackers, nil, nil, nil, nil)
	if sawLeak {
		t.Fatalf("CardState wrapper state leaked across permutations: grantSpy saw a pre-existing GrantedGoAgain when playing first")
	}
}
