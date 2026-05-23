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

// Tests that a non-attack pitch funding Aether Slash activates the +1 arcane rider.
func TestPitchAttribution_AetherSlashSingleNonAttackPitchFiresRider(t *testing.T) {
	hand := []card.Card{cards.AetherSlashRed{}, cards.MaleficIncantationBlue{}}
	d := deck.New(heroes.Viserai, nil, fillerDeck())
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), hand)
	if summary.Value != 5 {
		t.Fatalf("Value = %d, want 5 (Aether Slash 4 + rider 1)", summary.Value)
	}
}

// Tests that an attack-typed pitch funding Aether Slash does not activate the rider.
func TestPitchAttribution_AetherSlashAttackPitchDoesNotFireRider(t *testing.T) {
	hand := []card.Card{cards.AetherSlashRed{}, testutils.FakeYellowAttack()}
	d := deck.New(heroes.Viserai, nil, fillerDeck())
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), hand)
	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Aether Slash base power, no rider)", summary.Value)
	}
}

// Tests that Deathly Duet fires both riders when funded by one attack and one non-attack action.
func TestPitchAttribution_DeathlyDuetBothRidersFireFromMixedFunding(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, fillerDeck())
	hand := []card.Card{
		cards.DeathlyDuetRed{},
		cards.AetherSlashRed{},
		cards.MaleficIncantationBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), hand)
	if got := summary.Value; got != 8 {
		t.Fatalf("Value = %d, want 8 (Deathly Duet 4 + attack rider 2 + 2 runechants)", got)
	}
}

// Tests that a single pitch paying for multiple Aether Slashes activates the bonus on each.
func TestPitchAttribution_OneNonAttackPitchFundsMultipleAetherSlashes(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, fillerDeck())

	withNonAttack := []card.Card{
		cards.MauvrionSkiesRed{},
		cards.AetherSlashRed{}, cards.AetherSlashRed{},
		cards.MaleficIncantationBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), withNonAttack)
	if got := summary.Value; got != 15 {
		t.Errorf("non-attack pitch: Value = %d, want 15", got)
	}

	withAttack := []card.Card{
		cards.MauvrionSkiesRed{},
		cards.AetherSlashRed{}, cards.AetherSlashRed{},
		testutils.FakeBlueAttack(),
	}
	summary = sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), withAttack)
	if got := summary.Value; got != 13 {
		t.Errorf("attack pitch: Value = %d, want 13", got)
	}
}

// fillerDeck is a non-empty filler for EvalOneTurnForTesting calls that need a real deck
// behind the supplied hand — e.g. a card that draws mid-chain (Snatch) needs something
// to pull, a card that recycles to deck bottom (Relentless Pursuit) needs enough cards
// above it that the end-of-turn refill doesn't pull it back into hand, and an optimal
// line involving an extra draw needs a real card to draw or the engine sees no card
// advantage. Tests that don't read turn-2 state or care about mid-chain draws should
// pass nil directly.
func fillerDeck() []deck.Card {
	return []deck.Card{
		testutils.FakeBlueAttack(), testutils.FakeBlueAttack(),
		testutils.FakeBlueAttack(), testutils.FakeBlueAttack(),
		testutils.FakeBlueAttack(), testutils.FakeBlueAttack(),
		testutils.FakeBlueAttack(), testutils.FakeBlueAttack(),
	}
}
