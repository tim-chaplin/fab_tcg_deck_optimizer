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
)

// melodySurvived reports whether any aura in auras carries c's display name.
func melodySurvived(auras []gameengine.Aura, c card.Card) bool {
	for _, a := range auras {
		if a.CardName() == c.DisplayName() {
			return true
		}
	}
	return false
}

// Tests that each Enchanting Melody variant flips AuraCreated on Play and credits no
// direct damage.
func TestEnchantingMelody_SetsAuraCreated(t *testing.T) {
	cases := []card.Card{cards.EnchantingMelodyRed{}, cards.EnchantingMelodyYellow{}, cards.EnchantingMelodyBlue{}}
	for _, c := range cases {
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}

// Tests that Enchanting Melody's DamageAboutToBeTaken absorb fires before DamageTaken,
// so a sibling DamageTaken-subscribed aura (Arcane Cussing) survives.
func TestEnchantingMelody_AbsorbsBeforeDamageTakenSoCussingSurvives(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.EnchantingMelodyRed{}).
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		SetIncomingDamage(4).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, nil)

	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (Melody absorbs all 4 physical damage)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if melodySurvived(summary.State.Auras(), cards.EnchantingMelodyRed{}) {
		t.Errorf("Enchanting Melody still in auras after absorb; should have self-destructed")
	}
	cussingSurvived := false
	for _, a := range summary.State.Auras() {
		if a.CardName() == "Arcane Cussing [R]" {
			cussingSurvived = true
		}
	}
	if !cussingSurvived {
		t.Errorf("Arcane Cussing destroyed; expected it to survive because Melody absorbed all damage before DamageTaken fired\nFinal auras: %v", summary.State.Auras())
	}
}

// Tests that Enchanting Melody self-destructs on the EOT clause when no non-attack
// action was played this turn — no damage to absorb, no value credited, aura gone.
func TestEnchantingMelody_DestroysAtEndOfTurnWhenNoNonAttackActionPlayed(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.EnchantingMelodyRed{}).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, nil)

	if summary.Value != 0 {
		t.Errorf("Value = %d, want 0 (no damage to prevent)", summary.Value)
	}
	if melodySurvived(summary.State.Auras(), cards.EnchantingMelodyRed{}) {
		t.Errorf("Enchanting Melody still in auras after EOT; should self-destruct when no non-attack action played\nFinal auras: %v", summary.State.Auras())
	}
}

// Tests that the damage-absorb clause destroys Melody after preventing even 1 damage,
// regardless of the EOT clause being satisfied — confirms the destroy is part of the
// "instead destroy and prevent" rider, not the EOT "unless non-attack action" clause.
// A non-attack action played this turn satisfies the EOT clause (Melody would survive
// EOT), so any destruction observed must come from the absorb fire.
func TestEnchantingMelody_AbsorbDestroysEvenWhenNonAttackActionSatisfiesEOTClause(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.EnchantingMelodyRed{}).
		SetIncomingDamage(1).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, []card.Card{testutils.FakeRedAction()})

	if summary.Value != 1 {
		t.Errorf("Value = %d, want 1 (Melody prevented 1)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if !summary.State.NonAttackActionPlayed() {
		t.Fatalf("NonAttackActionPlayed = false; setup expected the FakeRedAction to be played so the EOT clause would NOT destroy Melody\nBestLine: %s",
			formatBestLine(summary.BestLine))
	}
	if melodySurvived(summary.State.Auras(), cards.EnchantingMelodyRed{}) {
		t.Errorf("Enchanting Melody still in auras; absorb-and-destroy clause should fire on any prevention amount")
	}
}

// Tests that Enchanting Melody absorbs arcane damage when only arcane is incoming —
// the DamageAboutToBeTaken fire on the arcane side picks the arcane branch and EM eats
// up to its printed N. Red N=4, incoming arcane 5 → 4 prevented, 1 leaks; EM destroyed.
func TestEnchantingMelody_AbsorbsArcaneDamage(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.EnchantingMelodyRed{}).
		SetArcaneIncomingDamage(5).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, nil)

	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (Red EM absorbs 4 arcane of 5 incoming)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if melodySurvived(summary.State.Auras(), cards.EnchantingMelodyRed{}) {
		t.Errorf("Enchanting Melody still in auras; should self-destruct after absorbing arcane")
	}
}

// Tests the "OR not AND" rule: a turn with both physical and arcane damage incoming
// fires DamageAboutToBeTaken twice (physical then arcane). EM destroys on the first fire
// (physical), so it cannot also eat the subsequent arcane fire. Setup: 3 physical and 3
// arcane incoming, Red EM (N=4). Physical fires first: EM absorbs 3, credits 3 Value,
// destroys. Arcane fire then finds no EM → 3 arcane leaks (no further Value credit). The
// damage cap means total Value is exactly 3, not 6.
func TestEnchantingMelody_CannotMixPhysicalAndArcaneSameTurn(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.EnchantingMelodyRed{}).
		SetIncomingDamage(3).
		SetArcaneIncomingDamage(3).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, nil)

	if summary.Value != 3 {
		t.Errorf("Value = %d, want 3 (EM absorbs only the physical fire; arcane fire finds no EM)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if melodySurvived(summary.State.Auras(), cards.EnchantingMelodyRed{}) {
		t.Errorf("Enchanting Melody still in auras; should self-destruct after absorbing physical")
	}
}

// Tests that each Enchanting Melody variant caps prevention at its printed N
// (Red 4 / Yellow 3 / Blue 2) — incoming=N+1 leaves 1 damage leaking through.
func TestEnchantingMelody_PreventionCapsScaleByVariant(t *testing.T) {
	cases := []struct {
		c       card.Card
		prevent int
	}{
		{cards.EnchantingMelodyRed{}, 4},
		{cards.EnchantingMelodyYellow{}, 3},
		{cards.EnchantingMelodyBlue{}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.c.DisplayName(), func(t *testing.T) {
			d := deck.New(heroes.Viserai, nil, nil)
			initial := gameengine.GameStateBuilder().
				CreateAuraFromCard(tc.c).
				SetIncomingDamage(tc.prevent + 1).
				Build()

			summary := sim.EvalOneTurnForTesting(d, initial, nil)

			if summary.Value != tc.prevent {
				t.Errorf("Value = %d, want %d (cap at %d prevention)\nBestLine: %s",
					summary.Value, tc.prevent, tc.prevent, formatBestLine(summary.BestLine))
			}
			if melodySurvived(summary.State.Auras(), tc.c) {
				t.Errorf("%s still in auras; should self-destruct after absorbing", tc.c.DisplayName())
			}
		})
	}
}
