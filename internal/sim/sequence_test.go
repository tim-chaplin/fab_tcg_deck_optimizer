package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

func TestBest_ViseraiMaleficShrillCombo(t *testing.T) {
	// Hero = Viserai. Best line: pitch the Blue Malefic, then play both Red Maleficas and the
	// Red Shrill. Each Malefic Play credits 0; Shrill (an attack action) fires both Malefics'
	// AttackAction triggers for +2 runes (OncePerTurn each, both at Count=3 so neither hits
	// zero). Plus Viserai's runechant on Shrill, plus Shrill's 4+3 aura-bonus = 11. Future
	// turns will keep ticking each Malefic for two more runes apiece, but those don't show
	// up in this turn's Value.
	h := []Card{
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationRed{},
		cards.MaleficIncantationRed{},
		cards.ShrillOfSkullformRed{},
	}
	got := Best(heroes.Viserai{}, nil, h, 4, nil, 0, nil)
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
	h := []Card{
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(heroes.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 5 {
		t.Fatalf("want value 5, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiReapingBladeMaleficsPlusShrill(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), play 2 Blue Malefics (2 dmg + 1 Runechant), then Red Shrill
	// (cost 2, 4+3 aura bonus + 1 Runechant = 8). Reaping Blade stays holstered — Shrill has no
	// Go again, so nothing can follow it. Value = 2 + 1 + 8 = 11.
	h := []Card{
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.ShrillOfSkullformRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(heroes.Viserai{}, weapons, h, 0, nil, 0, nil)
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
	h := []Card{
		cards.HocusPocusBlue{},
		cards.OathOfTheArknightRed{},
		cards.MaleficIncantationRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(heroes.Viserai{}, weapons, h, 0, nil, 0, nil)
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
	h := []Card{
		cards.HocusPocusBlue{},
		cards.MaleficIncantationBlue{},
		cards.RunicReapingRed{},
		cards.ShrillOfSkullformRed{},
	}
	got := Best(heroes.Viserai{}, nil, h, 0, nil, 0, nil)
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
	h := []Card{
		cards.HocusPocusBlue{},
		cards.MaleficIncantationBlue{},
		cards.MauvrionSkiesRed{},
		cards.ShrillOfSkullformRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(heroes.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 16 {
		t.Fatalf("want value 16, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_ViseraiMauvrionPredictsDrowningDireDominate pins the full resolution stack each
// card's grants need to settle in the right order:
//
//	(1) target plays → hero ability fires (Viserai creates a Runechant → AuraCreated=true)
//	(2) target's Dominate clause resolves (Drowning Dire sees the aura → gains Dominate)
//	(3) Mauvrion's OnAttack trigger fires against the fully-resolved target state and
//	    credits its Runechant rider iff the attack is now likely to hit.
//
// Requires hero.OnCardPlayed running before Play (so DD's aura check sees Viserai's
// Runechant), DD's conditional Dominate grant, and Mauvrion's ephemeral trigger reading
// target.EffectiveDominate() at fire time. If any of those regresses, this test drops to 6
// (no Mauvrion rider) or less.
//
// Line: pitch YellowAttack (2 res), play Mauvrion Red (cost 0, grants go-again + "if hits,
// create 3 Runechants" to the next Runeblade attack action), play Drowning Dire Red (cost 2,
// attack 5). Viserai creates a Runechant → Drowning Dire has Dominate → 5+ dominating
// attack lands → Mauvrion's rider fires for 3 Runechants. Value = 3 + 5 + 1 = 9.
func TestBest_ViseraiMauvrionPredictsDrowningDireDominate(t *testing.T) {
	h := []Card{
		cards.MauvrionSkiesRed{},
		cards.DrowningDireRed{},
		testutils.YellowAttack{},
	}
	got := Best(heroes.Viserai{}, nil, h, 0, nil, 0, nil)
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
	order := []Card{
		cards.MauvrionSkiesRed{},
		cards.RuneragerSwarmRed{},
		cards.ShrillOfSkullformRed{},
		weapons.ReapingBlade{},
	}
	ctx := NewSequenceContextForTest(heroes.Viserai{}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.PlaySequence(order); legal {
		t.Fatalf("ordering %v should be illegal (Shrill has no go-again and Mauvrion granted Runerager instead)",
			CardNames(order))
	}
}

func TestBest_ViseraiMauvrionChainsShrillIntoRuneragerIntoWeapon(t *testing.T) {
	// Pitch Blue Hocus → Mauvrion → Shrill → Runerager → Reaping Blade. Value = 3 + 7 + 3 + 3 + 2
	// Viserai runechants = 18.
	h := []Card{
		cards.HocusPocusBlue{},
		cards.MauvrionSkiesRed{},
		cards.RuneragerSwarmRed{},
		cards.ShrillOfSkullformRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(heroes.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 18 {
		t.Fatalf("want value 18, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_StateValueMatchesSummedReturns pins that state.Value equals the explicit-summation
// total a hand's Plays would produce. Hand: 2 Blues + 2 Reds vs no incoming damage. The optimal
// chain pitches one Blue (3 resource) and chains the other Blue + 2 Reds — total 1 + 3 + 3 = 7
// damage.
func TestBest_StateValueMatchesSummedReturns(t *testing.T) {
	h := []Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
	if got.Value != 7 {
		t.Errorf("Value = %d, want 7 (Blue 1 + Red 3 + Red 3 chain off one Blue pitch). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBestSequence_CardStateGrantsDontLeakAcrossPermutations pins the per-permutation reset
// contract: the permutation loop in bestSequence must allocate fresh *CardState wrappers per
// permutation so a grant applied by one permutation's Play() can't bleed into a later
// permutation's legality/effect checks.
//
// Setup: attackers = [GrantAll, GrantSpy, GrantAll]. The permute order emits GrantAll-first
// permutations before the first GrantSpy-first permutation. Each GrantAll-first permutation
// mutates GrantedGoAgain on the wrappers for the cards behind it. When GrantSpy later plays
// FIRST in its own permutation, its CardsRemaining contains the other two cards' wrappers —
// which must be fresh (GrantedGoAgain=false), since no card has played yet in this permutation.
// If the wrappers were reused across permutations the spy would see leaked grants and trip.
func TestBestSequence_CardStateGrantsDontLeakAcrossPermutations(t *testing.T) {
	var sawLeak bool
	attackers := []Card{GrantAll{}, GrantSpy{Saw: &sawLeak}, GrantAll{}}
	ctx := NewSequenceContextForTest(StubHero, nil, nil, 1_000_000, 0, len(attackers))
	_, _, _ = ctx.BestSequence(attackers)
	if sawLeak {
		t.Fatalf("CardState wrapper state leaked across permutations: GrantSpy saw a pre-existing GrantedGoAgain when playing first")
	}
}
