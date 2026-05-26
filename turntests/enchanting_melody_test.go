package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

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
	for _, a := range summary.State.Auras() {
		if a.CardName() == "Enchanting Melody [R]" {
			t.Errorf("Enchanting Melody still in auras after absorb; should have self-destructed")
		}
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
