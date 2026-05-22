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

// recompenseDeck is five fillers — enough to deal a full hand next turn. This turn draws
// no cards, so deck contents don't affect the line under test.
func recompenseDeck() *deck.Deck {
	return deck.New(heroes.Viserai{}, nil, []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	})
}

// Tests that Talisman of Recompense, in play at start of turn, turns a 1-resource red
// pitch into three resources — enough to fund a 3-cost attack the hand can't otherwise
// afford. The optimal line pitches the red card and plays the attack for 4.
func TestTalismanOfRecompense_RedPitchFundsThreeCostAttack(t *testing.T) {
	hand := []card.Card{testutils.RedPitch{}, testutils.GenericAttack(3, 4)}
	state := gameengine.GameStateBuilder().
		CreateItemFromCard(cards.TalismanOfRecompenseYellow{}).
		Build()

	summary := sim.EvalOneTurnForTesting(recompenseDeck(), state, hand)

	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Recompense funds the 3-cost attack off a red pitch)", summary.Value)
	}
	if n := len(summary.State.Items()); n != 0 {
		t.Fatalf("Items() = %d after turn, want 0 (Recompense is destroyed when it fires)", n)
	}
}

// Control: with no Talisman in play the same hand can't fund the 3-cost attack — a red
// pitch yields only one resource — so the line can't reach 4.
func TestTalismanOfRecompense_NoTalismanCannotFundAttack(t *testing.T) {
	hand := []card.Card{testutils.RedPitch{}, testutils.GenericAttack(3, 4)}

	summary := sim.EvalOneTurnForTesting(recompenseDeck(), gameengine.GameStateBuilder().Build(), hand)

	if summary.Value != 0 {
		t.Fatalf("Value = %d, want 0 (a red pitch funds only 1 resource — the 3-cost attack can't be played)", summary.Value)
	}
}
