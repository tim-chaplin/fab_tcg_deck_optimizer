package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
)

// TestArcanicSpike_ArcaneDamageNotDealtReturnsBaseAttack covers the unsatisfied branch: when
// TurnState.ArcaneDamageDealt is false the +2{p} rider doesn't fire and Play returns only the
// printed attack.
func TestArcanicSpike_ArcaneDamageNotDealtReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ArcanicSpikeRed{}, 5},
		{cards.ArcanicSpikeYellow{}, 4},
		{cards.ArcanicSpikeBlue{}, 3},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, ArcaneDamageDealt=false)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestArcanicSpike_ArcaneDamageDealtTriggersBonus exercises the satisfied path: when
// ArcaneDamageDealt is set (an earlier attack fired a Runechant, or a direct-arcane card flipped
// the flag) the +2{p} rider activates and Play returns attack + 2.
func TestArcanicSpike_ArcaneDamageDealtTriggersBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ArcanicSpikeRed{}, 5 + 2},
		{cards.ArcanicSpikeYellow{}, 4 + 2},
		{cards.ArcanicSpikeBlue{}, 3 + 2},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetArcaneDamageDealt(true).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (attack + arcane bonus)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that a Runechant already in play fires when Arcanic Spike is played — a Runechant
// resolves before the attack that triggers it — so it self-destroys and flips
// ArcaneDamageDealt, turning on Arcanic Spike's "dealt arcane damage this turn" +2{p}
// rider. Companion to the Malefic Incantation test, which covers the create-side ordering.
func TestArcanicSpike_RunechantFiresBeforeAttackAndTriggersRider(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		AddAura(token.NewRunechant(1)).
		SetIncomingDamage(0).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{testutils.BluePitch{}, cards.ArcanicSpikeRed{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if summary.Value != 7 {
		t.Fatalf("Value = %d, want 7 (Arcanic Spike 5 + 2 arcane rider — the prior Runechant flipped ArcaneDamageDealt)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if got := summary.State.RunechantCount(); got != 0 {
		t.Fatalf("end-of-turn RunechantCount = %d, want 0 (the Runechant fired and self-destroyed when Arcanic Spike was played)", got)
	}
}
