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

// Tests that a Rifting hit lets the next non-attack action card play without an action point.
func TestRifting_HitLetsNextNonAttackActionPlayAsInstant(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	prior := gameengine.GameStateBuilder().SetIncomingDamage(0).Build()
	// Rifting (Blue, power 4) hits; BluePitch funds its cost 2; NonAttack is the non-attack
	// action that, with no AP left after Rifting, can only be played via the instant grant.
	hand := []card.Card{cards.RiftingBlue{}, testutils.FakeRedAction(), testutils.FakeBlueResource()}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	played := false
	for _, a := range summary.BestLine {
		if a.Card.Name() == "FakeRedAction" && a.Role == card.Attack {
			played = true
		}
	}
	if !played {
		t.Errorf("FakeRedAction not played — Rifting's instant grant didn't let it into the chain\nBestLine: %s",
			formatBestLine(summary.BestLine))
	}
}
