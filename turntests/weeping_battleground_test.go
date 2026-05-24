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
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
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
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 3 {
		t.Errorf("Value = %d, want 3 (3 block only; banish fizzles). Roles=[%s]",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests that a defender a DR banished during the defense phase doesn't ALSO appear in
// the post-defense graveyard. Pre-fix, runDefense rebuilt chainGraveyard from the original
// defenders list unconditionally, so the banished card ended up in BOTH the banished zone
// AND the graveyard — double-counted across the per-turn zone accounting, and visible to
// subsequent BanishFromGraveyard / RecycleFromGraveyard scans as a phantom copy.
func TestBest_WeepingBattlegroundBanishedAuraOnlyInBanishedZone(t *testing.T) {
	h := []card.Card{cards.WeepingBattlegroundRed{}, zeroDefenseAura{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)

	// Conservation: pre = 2 cards in hand, 0 elsewhere. Post must equal 2.
	post := summary.State.Deck().Size() + len(summary.State.Graveyard()) + len(summary.State.Hand()) + len(summary.State.Banished())
	if summary.State.Arsenal() != nil {
		post++
	}
	if post != 2 {
		t.Errorf("total card count = %d, want 2 (DR banish moved zeroDefenseAura grav -> banished; no card should spawn/vanish). deck=%d grav=%v hand=%v banished=%v arsenal=%v",
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
