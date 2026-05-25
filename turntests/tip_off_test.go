package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Tip-Off's Instant-discard mode wins the partition when a same-chain
// marked-defender rider can read the resulting mark. Hand is Tip-Off Blue + Outed Red:
// mode 0 path is "pitch Outed for 1{r}, play Tip-Off Blue mode 0 for 3{p}" = 3 total;
// mode 1 path is "Tip-Off Blue mode 1 (0 damage, marks, go-again) then Outed Red
// (3 + 1 marked-defender = 4)" = 4 total. The chain runner's per-attack mark clear is
// gated on positive EffectiveAttack, so mode 1's 0-power swing leaves the mark intact
// for Outed to read.
func TestTipOff_InstantModeMarksOpponentForSameChainOutedBonus(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.TipOffBlue{}, cards.OutedRed{}}

	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (Tip-Off mode 1 marks + Outed 3{p} + 1 marked-defender)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if !bestLineHasRole(summary.BestLine, cards.TipOffBlue{}, card.Attack) {
		t.Errorf("BestLine missing Tip-Off as Attack: %s", formatBestLine(summary.BestLine))
	}
	if !bestLineHasRole(summary.BestLine, cards.OutedRed{}, card.Attack) {
		t.Errorf("BestLine missing Outed as Attack: %s", formatBestLine(summary.BestLine))
	}
}
