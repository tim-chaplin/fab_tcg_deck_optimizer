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

// melodySurvived reports whether any aura in auras carries Enchanting Melody [R]'s name.
func melodySurvived(auras []gameengine.Aura) bool {
	for _, a := range auras {
		if a.CardName() == "Enchanting Melody [R]" {
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
	if melodySurvived(summary.State.Auras()) {
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
	if melodySurvived(summary.State.Auras()) {
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
	if melodySurvived(summary.State.Auras()) {
		t.Errorf("Enchanting Melody still in auras; absorb-and-destroy clause should fire on any prevention amount")
	}
}

// Tests that Melody caps prevention at 4 damage — incoming=5 leaves 1 leaking through.
func TestEnchantingMelody_CapsPreventionAtFour(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.EnchantingMelodyRed{}).
		SetIncomingDamage(5).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, nil)

	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (cap at 4 prevention)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if melodySurvived(summary.State.Auras()) {
		t.Errorf("Enchanting Melody still in auras; should self-destruct after absorbing")
	}
}
