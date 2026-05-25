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

// Tests that Eirina's Prayer's DefensiveInstant routing prevents arcane damage scaled by
// the revealed deck-top pitch. Hand = [Eirina Red, FakeBlueResource]; matchup has
// IncomingDamage=2 (so the partition allows Defend role) plus ArcaneIncomingDamage=5.
// Deck top is pitch 2, so X = 6 - 2 = 4. Eirina prevents min(4, 5) = 4 arcane and
// credits 4 Value; FakeBlue funds the 1{r} cost. Eirina's printed Defense is 0, so the
// chain-step DR delta contributes nothing — the 4 comes entirely from the arcane
// prevention rider.
func TestEirinasPrayer_PreventsArcaneDamageScaledByDeckTopPitch(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, []deck.Card{testutils.FakeYellowResource()})
	hand := []card.Card{cards.EirinasPrayerRed{}, testutils.FakeBlueResource()}
	initial := gameengine.GameStateBuilder().
		SetIncomingDamage(2).
		SetArcaneIncomingDamage(5).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, hand)

	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (Eirina prevents X=6-2=4 arcane; deck-top pitch 2)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if !bestLineHasRole(summary.BestLine, cards.EirinasPrayerRed{}, card.Defend) {
		t.Errorf("BestLine missing Eirina as Defend: %s", formatBestLine(summary.BestLine))
	}
}

// Tests that Eirina's prevention is capped at the remaining ArcaneIncomingDamage — when
// X (= 6 - deck-top pitch) exceeds incoming arcane, the credited Value is the incoming
// total, not the formula amount. Deck top is pitch 1 → X = 5; incoming arcane = 2 →
// prevented = 2.
func TestEirinasPrayer_PreventionCapsAtRemainingArcane(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, []deck.Card{testutils.FakeRedResource()})
	hand := []card.Card{cards.EirinasPrayerRed{}, testutils.FakeBlueResource()}
	initial := gameengine.GameStateBuilder().
		SetIncomingDamage(2).
		SetArcaneIncomingDamage(2).
		Build()

	summary := sim.EvalOneTurnForTesting(d, initial, hand)

	if summary.Value != 2 {
		t.Errorf("Value = %d, want 2 (X=5 but only 2 arcane incoming)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}
