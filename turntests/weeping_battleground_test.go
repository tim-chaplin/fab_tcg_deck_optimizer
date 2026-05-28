package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// zeroDefenseAura is an aura-typed card used to populate the prior-turn graveyard so
// tests can isolate Weeping Battleground's +1 arcane banish rider. Defense / Pitch are
// zero so it wouldn't contribute to a defense total if drawn into hand.
type zeroDefenseAura struct{}

func (zeroDefenseAura) ID() ids.CardID                                     { return ids.InvalidCard }
func (zeroDefenseAura) Name() string                                       { return "zeroDefenseAura" }
func (zeroDefenseAura) DisplayName() string                                { return "zeroDefenseAura" }
func (zeroDefenseAura) Cost() int                                          { return 0 }
func (zeroDefenseAura) Pitch() int                                         { return 0 }
func (zeroDefenseAura) Attack() int                                        { return 0 }
func (zeroDefenseAura) Defense() int                                       { return 0 }
func (zeroDefenseAura) Types(card.GameEngine) card.TypeSet                 { return card.NewTypeSet(card.TypeAura) }
func (zeroDefenseAura) GoAgain(card.GameEngine) bool                       { return false }
func (zeroDefenseAura) Play(card.GameEngine, card.Logger, *card.CardState) {}

// Tests that Weeping Battleground banishes a pre-existing graveyard aura for +1 arcane
// while also defending.
func TestBest_WeepingBattlegroundBanishesAuraFromGraveyard(t *testing.T) {
	h := []card.Card{cards.WeepingBattlegroundRed{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	state := gameengine.GameStateBuilder().
		SetIncomingPhysicalDamage(4).
		SetGraveyard([]card.Card{zeroDefenseAura{}}).
		Build()
	summary := sim.EvalOneTurnForTesting(d, state, h)
	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (3 block + 1 arcane from banish). Roles=[%s]",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
}

// TestBest_WeepingBattlegroundFizzlesWithoutAura: hand is just Weeping Battleground — no
// aura anywhere, so the banish rider fizzles. WB still blocks 3 of the 4 incoming. Value = 3.
func TestBest_WeepingBattlegroundFizzlesWithoutAura(t *testing.T) {
	h := []card.Card{cards.WeepingBattlegroundRed{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(4).Build(), h)
	if summary.Value != 3 {
		t.Errorf("Value = %d, want 3 (3 block only; banish fizzles). Roles=[%s]",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests that an aura banished by Weeping Battleground ends up only in the banished
// zone, not duplicated across both zones.
func TestBest_WeepingBattlegroundBanishedAuraOnlyInBanishedZone(t *testing.T) {
	h := []card.Card{cards.WeepingBattlegroundRed{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	state := gameengine.GameStateBuilder().
		SetIncomingPhysicalDamage(4).
		SetGraveyard([]card.Card{zeroDefenseAura{}}).
		Build()
	summary := sim.EvalOneTurnForTesting(d, state, h)

	post := summary.State.Deck().Size() + len(summary.State.Graveyard()) + len(summary.State.Hand()) + len(summary.State.Banished())
	if summary.State.Arsenal() != nil {
		post++
	}
	if post != 2 {
		t.Errorf("total card count = %d, want 2 (WB + zeroDefenseAura). deck=%d grav=%v hand=%v banished=%v arsenal=%v",
			post, summary.State.Deck().Size(), summary.State.Graveyard(), summary.State.Hand(), summary.State.Banished(), summary.State.Arsenal())
	}

	// Specifically: the banished aura must be in Banished, not Graveyard.
	for _, c := range summary.State.Graveyard() {
		if _, isAura := c.(zeroDefenseAura); isAura {
			t.Errorf("zeroDefenseAura found in graveyard %v — it was banished by Weeping Battleground and should only appear in banished zone %v",
				summary.State.Graveyard(), summary.State.Banished())
		}
	}
	foundInBanished := false
	for _, c := range summary.State.Banished() {
		if _, isAura := c.(zeroDefenseAura); isAura {
			foundInBanished = true
			break
		}
	}
	if !foundInBanished {
		t.Errorf("zeroDefenseAura missing from banished %v — Weeping Battleground's banish rider didn't take effect", summary.State.Banished())
	}
}
