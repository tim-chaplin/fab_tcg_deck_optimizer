package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Tit for Tat's go-again lets a follow-up attack swing and leaves the hero untapped.
func TestTitForTat_GoAgainLetsAttackFollow(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.TitForTatBlue{}, testutils.FakeRedAttack().WithPower(2)}

	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)

	if got := summary.Value; got != 2 {
		t.Fatalf("Value = %d, want 2 (Tit for Tat's go-again lets the 2-power attack swing)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
	if summary.State.HeroTapped() {
		t.Errorf("HeroTapped = true after Tit for Tat played; want false (Play calls UntapHero)")
	}
}
