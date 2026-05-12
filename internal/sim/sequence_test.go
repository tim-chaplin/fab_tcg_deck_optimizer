package sim_test

import (
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func TestBest_ViseraiMaleficShrillCombo(t *testing.T) {
	// Hero = Viserai. Best line: pitch the Blue Malefic, then play both Red Maleficas and the
	// Red Shrill. Each Malefic Play credits 0; Shrill (an attack action) fires both Malefics'
	// AttackAction triggers for +2 runes (OncePerTurn each, both at Count=3 so neither hits
	// zero). Plus Viserai's runechant on Shrill, plus Shrill's 4+3 aura-bonus = 11. Future
	// turns will keep ticking each Malefic for two more runes apiece, but those don't show
	// up in this turn's Value.
	h := []card.Card{
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationRed{},
		cards.MaleficIncantationRed{},
		cards.ShrillOfSkullformRed{},
	}
	got := Best(nil, h, Matchup{IncomingDamage: 4}, nil, Prior{Hero: heroes.Viserai{}})
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
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(weapons, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: heroes.Viserai{}})
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
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.ShrillOfSkullformRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(weapons, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: heroes.Viserai{}})
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
		cards.HocusPocusBlue{},
		cards.OathOfTheArknightRed{},
		cards.MaleficIncantationRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(weapons, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: heroes.Viserai{}})
	if got.Value != 8 {
		t.Fatalf("want value 8, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_RunicReapingPrefersAttackPitch(t *testing.T) {
	// Pitching the Blue Hocus Pocus (attack-typed, pitch 3) pays for Runic Reaping + Shrill
	// AND satisfies Runic Reaping's pitched-attack rider, granting +1 to Shrill via
	// BonusAttack. Runic Reaping's "if this hits, create N Runechants" clause is appended
	// to Shrill's OnHit (same shape as Mauvrion Skies) and fires post-buff:
	// target.EffectiveAttack = printed 4 + BonusAttack 1 = 5, which falls OUT of the {1,4,7}
	// s.LikelyToHit window, so the runechant rider drops. The only damage
	// on Runic Reaping's slot is Viserai's runechant for the prior non-attack action.
	// Blue Malefic (1 arcane + 1 Viserai runechant = 2) → Runic Reaping (0 own damage + 1
	// Viserai runechant = 1) → Shrill (4 base + 3 aura-created bonus + 1 BonusAttack = 8).
	// Value = 2 + 1 + 8 = 11.
	h := []card.Card{
		cards.HocusPocusBlue{},
		cards.MaleficIncantationBlue{},
		cards.RunicReapingRed{},
		cards.ShrillOfSkullformRed{},
	}
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: heroes.Viserai{}})
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
		cards.HocusPocusBlue{},
		cards.MaleficIncantationBlue{},
		cards.MauvrionSkiesRed{},
		cards.ShrillOfSkullformRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(weapons, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: heroes.Viserai{}})
	if got.Value != 16 {
		t.Fatalf("want value 16, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests the resolution order: Viserai creates a Runechant on Drowning Dire's play →
// DD gains Dominate → Mauvrion's OnAttack fires against the now-likely-to-hit attack.
func TestBest_ViseraiMauvrionPredictsDrowningDireDominate(t *testing.T) {
	h := []card.Card{
		cards.MauvrionSkiesRed{},
		cards.DrowningDireRed{},
		testutils.YellowAttack{},
	}
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: heroes.Viserai{}})
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
		cards.MauvrionSkiesRed{},
		cards.RuneragerSwarmRed{},
		cards.ShrillOfSkullformRed{},
		weapons.ReapingBlade{}.Ability(),
	}
	ctx := NewSequenceContextForTest(heroes.Viserai{}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.PlaySequence(order); legal {
		t.Fatalf("ordering %v should be illegal (Shrill has no go-again and Mauvrion granted Runerager instead)",
			testutils.CardNamesSim(order))
	}
}

func TestBest_ViseraiMauvrionChainsShrillIntoRuneragerIntoWeapon(t *testing.T) {
	// Pitch Blue Hocus → Mauvrion → Shrill → Runerager → Reaping Blade. Value = 3+7+3+3+2 +
	// Viserai runechants = 18.
	h := []card.Card{
		cards.HocusPocusBlue{},
		cards.MauvrionSkiesRed{},
		cards.RuneragerSwarmRed{},
		cards.ShrillOfSkullformRed{},
	}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(weapons, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: heroes.Viserai{}})
	if got.Value != 18 {
		t.Fatalf("want value 18, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that state.Value equals the summed Play returns (no double-counting or drops).
func TestBest_StateValueMatchesSummedReturns(t *testing.T) {
	h := []card.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, Prior{Hero: testutils.Hero{Intel: 4}})
	if got.Value != 7 {
		t.Errorf("Value = %d, want 7 (Blue 1 + Red 3 + Red 3 chain off one Blue pitch). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that bestSequence allocates fresh *CardState wrappers per permutation so a grant
// applied in one permutation can't leak into a later permutation's checks.
func TestBestSequence_CardStateGrantsDontLeakAcrossPermutations(t *testing.T) {
	var sawLeak bool
	attackers := []card.Card{testutils.GrantAll{}, testutils.GrantSpy{Saw: &sawLeak}, testutils.GrantAll{}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 1_000_000, 0, len(attackers))
	_, _, _ = ctx.BestSequence(attackers)
	if sawLeak {
		t.Fatalf("CardState wrapper state leaked across permutations: GrantSpy saw a pre-existing GrantedGoAgain when playing first")
	}
}

// Tests that a non-Go-again attack followed by a non-Instant card rejects the chain — the
// AP pool drains to 0 on the first card and the second can't pay its 1 AP cost.
func TestPlaySequence_NonGoAgainStopsChain(t *testing.T) {
	order := []card.Card{testutils.NoGoAgainAttackStub{}, testutils.NoGoAgainAttackStub{}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.PlaySequence(order); legal {
		t.Fatalf("ordering %v should be illegal (no Go again grant after card 0)", testutils.CardNamesSim(order))
	}
}

// Tests that an Instant follow-up after a non-Go-again card resolves legally — Instants cost
// 0 AP so the empty pool isn't a barrier.
func TestPlaySequence_InstantBypassesAPRequirement(t *testing.T) {
	order := []card.Card{testutils.NoGoAgainAttackStub{}, testutils.InstantStub{}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 1_000_000, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("ordering %v should be legal (Instant costs 0 AP)", testutils.CardNamesSim(order))
	}
	if dmg != 1 {
		t.Errorf("dmg = %d, want 1 (NoGoAgainAttack 1 + Instant 0)", dmg)
	}
}

// Tests that a chain of two Instants opens with neither card paying AP — the AP pool stays
// at 1 the whole way, and a non-Instant follow-up still works (which would fail if Instants
// had silently consumed AP).
func TestPlaySequence_InstantsDontConsumeAP(t *testing.T) {
	order := []card.Card{testutils.InstantStub{}, testutils.InstantStub{}, testutils.NoGoAgainAttackStub{}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 1_000_000, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("ordering %v should be legal (Instants cost 0 AP)", testutils.CardNamesSim(order))
	}
	if dmg != 1 {
		t.Errorf("dmg = %d, want 1 (two Instants at 0 + NoGoAgainAttack 1)", dmg)
	}
}
