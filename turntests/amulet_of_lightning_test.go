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

// Tests that without a Lightning fuse the ability stays gated off, so the item isn't destroyed.
func TestAmuletOfLightning_NoFuseAbilityGatedOff(t *testing.T) {
	hero := testutils.Hero{Intel: 4}
	a1 := testutils.FakeRedAttack().
		WithPower(3).
		WithName("Attack1")
	a2 := testutils.FakeRedAttack().
		WithPower(3).
		WithName("Attack2")
	state := gameengine.GameStateBuilder().
		SetHero(hero).
		CreateItemFromCard(cards.AmuletOfLightningBlue{}).
		Build()
	d := deck.New(hero, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, state, []card.Card{a1, a2})

	if got := len(summary.State.Items()); got != 1 {
		t.Errorf("Items() = %d, want 1 — the ability is gated off without a fuse, so the item persists\n"+
			"BestLine: %s", got, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests that a Lightning fuse unlocks the ability — it grants go again (enabling a second swing)
// and self-destructs the item.
func TestAmuletOfLightning_FusedGrantsGoAgainAndSelfDestructs(t *testing.T) {
	hero := testutils.Hero{Intel: 4}
	fuse := testutils.FakeRedInstant().
		WithName("Fuse").
		WithPlay(func(ge card.GameEngine, _ card.Logger, _ *card.CardState) { ge.SetLightningFused(true) })
	a1 := testutils.FakeRedAttack().
		WithPower(3).
		WithName("Attack1")
	a2 := testutils.FakeRedAttack().
		WithPower(3).
		WithName("Attack2")
	state := gameengine.GameStateBuilder().
		SetHero(hero).
		CreateItemFromCard(cards.AmuletOfLightningBlue{}).
		Build()
	d := deck.New(hero, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, state, []card.Card{fuse, a1, a2})

	if got := len(summary.State.Items()); got != 0 {
		t.Errorf("Items() = %d after the ability resolved, want 0 — it should self-destruct\nBestLine: %s",
			got, sim.FormatBestLine(summary.BestLine))
	}
	if got := summary.Value; got < 6 {
		t.Errorf("Value = %d, want >= 6 (two swings enabled by the granted go again)\nBestLine: %s",
			got, sim.FormatBestLine(summary.BestLine))
	}
}
