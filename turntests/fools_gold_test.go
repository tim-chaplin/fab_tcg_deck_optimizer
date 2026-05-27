package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that a sole-Fool's-Gold hand stays in hand at end of turn. Promoting a Resource
// card to arsenal would render it unplayable (resources can't be pitched from arsenal),
// so the post-hoc promotion's isArsenalEligible filter excludes IsResource cards.
// Pitching it would be illegal in FaB rules with nothing to pay for, so the optimiser
// must keep it Held.
func TestFoolsGold_SoleHandCardStaysInHand(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.FoolsGoldYellow{}}
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	foundInHand := false
	for _, c := range summary.State.Hand() {
		if c.ID() == (cards.FoolsGoldYellow{}).ID() {
			foundInHand = true
			break
		}
	}
	if !foundInHand {
		t.Errorf("Fool's Gold missing from end-of-turn hand %v (arsenal=%v) — a sole-Resource hand "+
			"has nothing legal to pay for, so the card must stay Held\nBestLine: %s",
			summary.State.Hand(), summary.State.Arsenal(), sim.FormatBestLine(summary.BestLine))
	}
}

// End-to-end: Trade In plays as an attack and tries to discard a Held card. When Fool's
// Gold is the held target, the discard fires Fool's Gold's OnDiscard hook and a Gold token
// lands. The optimiser picks the Held-Fool's-Gold + Trade-In-discards branch over the
// pitch-for-resources branch because the discard's "draw a card" sweetener plus the Gold
// stockpile outweighs the {r}{r} pitch on a 0-cost attack.
func TestFoolsGold_OnDiscardFiresThroughTradeIn(t *testing.T) {
	// Deck has spare cards Trade In's "draw a card" can land on (otherwise DrawOne fizzles).
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{
		testutils.FakeRedAttack().WithName("DeckTop"),
		testutils.FakeRedAttack().WithName("DeckMid"),
	})
	hand := []card.Card{cards.TradeInRed{}, cards.FoolsGoldYellow{}}
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	if got := summary.State.GoldCount(); got != 1 {
		t.Errorf("GoldCount = %d, want 1 (Trade In discards Fool's Gold → OnDiscard mints Gold)\nBestLine: %s",
			got, sim.FormatBestLine(summary.BestLine))
	}
}
