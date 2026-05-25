package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
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

// Tests that a "next attack action card +4" buff does NOT land on Tip-Off in mode 1 —
// mode 1's type-line is Generic Instant, so the buff scanning CardsRemaining for an
// attack action skips past Tip-Off. Hand is Tip-Off Blue + a fake red Go-Again
// non-attack action whose Play grants +4 to the next attack action. The only feasible
// path that scores is "pitch the fake red for 1{r}, play Tip-Off Blue mode 0 for 3{p}" =
// 3. Playing the fake red (queueing the +4) and then Tip-Off mode 1 deals 0 (mode 1 is
// Instant, so the buff doesn't apply and it credits no attack damage); mode 0 isn't
// fundable in that branch because the fake red was played rather than pitched. If the
// engine were treating mode 1 as an attack action regardless of the per-mode dispatch,
// the +4 would land and Value would jump to 4.
func TestTipOff_InstantModeDoesNotConsumeAttackActionBuff(t *testing.T) {
	buff := testutils.FakeRedAction().
		WithName("FakeBuffPlus4").
		WithGoAgain().
		WithPlay(func(ge card.GameEngine, _ card.Logger, _ *card.CardState) {
			cards.GrantNextCardBonusAttack(ge, 4, card.IsAttackAction)
		})

	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.TipOffBlue{}, buff}

	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	if summary.Value != 3 {
		t.Errorf("Value = %d, want 3 (mode 0 path: pitch buff, play Tip-Off Blue for 3{p}). A jump to 4 means the +4 buff is leaking onto mode 1 despite its Generic Instant type-line.\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}
