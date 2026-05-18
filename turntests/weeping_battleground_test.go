package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// zeroDefenseAura is an aura-typed card that blocks for nothing — used to park an aura in
// the graveyard via the plain-block seeding path without adding to the defense total, so
// tests can isolate Weeping Battleground's +1 arcane banish rider.
type zeroDefenseAura struct{}

func (zeroDefenseAura) ID() ids.CardID                                     { return ids.InvalidCard }
func (zeroDefenseAura) Name() string                                       { return "zeroDefenseAura" }
func (zeroDefenseAura) DisplayName() string                                { return "zeroDefenseAura" }
func (zeroDefenseAura) Cost(card.GameEngine) int                           { return 0 }
func (zeroDefenseAura) Pitch() int                                         { return 0 }
func (zeroDefenseAura) Attack() int                                        { return 0 }
func (zeroDefenseAura) Defense() int                                       { return 0 }
func (zeroDefenseAura) Types(card.GameEngine) card.TypeSet                 { return card.NewTypeSet(card.TypeAura) }
func (zeroDefenseAura) GoAgain(card.GameEngine) bool                       { return false }
func (zeroDefenseAura) Play(card.GameEngine, card.Logger, *card.CardState) {}

// Tests that Weeping Battleground banishes a same-turn-blocked aura from the graveyard
// for 1 arcane while also defending.
func TestBest_WeepingBattlegroundBanishesAuraFromGraveyard(t *testing.T) {
	h := []card.Card{cards.WeepingBattlegroundRed{}, zeroDefenseAura{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
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
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if summary.Value != 3 {
		t.Errorf("Value = %d, want 3 (3 block only; banish fizzles). Roles=[%s]",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
}
