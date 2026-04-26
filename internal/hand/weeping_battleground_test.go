package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
)

// zeroDefenseAura is an aura-typed card that blocks for nothing — used to park an aura in
// the graveyard via the plain-block seeding path without adding to the defense total, so
// tests can isolate Weeping Battleground's +1 arcane banish rider.
type zeroDefenseAura struct{}

func (zeroDefenseAura) ID() card.ID                               { return card.Invalid }
func (zeroDefenseAura) Name() string                              { return "zeroDefenseAura" }
func (zeroDefenseAura) Cost(*card.TurnState) int                  { return 0 }
func (zeroDefenseAura) Pitch() int                                { return 0 }
func (zeroDefenseAura) Attack() int                               { return 0 }
func (zeroDefenseAura) Defense() int                              { return 0 }
func (zeroDefenseAura) Types() card.TypeSet                       { return card.NewTypeSet(card.TypeAura) }
func (zeroDefenseAura) GoAgain() bool                             { return false }
func (zeroDefenseAura) Play(*card.TurnState, *card.CardState) {}
// TestBest_WeepingBattlegroundBanishesAuraFromGraveyard: hand is Weeping Battleground + an
// aura filler. The filler plain-blocks (0 defense, but lands in the graveyard via the
// defense-phase seeding), WB plays as DR, banishes the filler for 1 arcane, and blocks 3 of
// the 4 incoming. Value = 3 prevented + 1 arcane = 4.
func TestBest_WeepingBattlegroundBanishesAuraFromGraveyard(t *testing.T) {
	h := []card.Card{runeblade.WeepingBattlegroundRed{}, zeroDefenseAura{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 4 {
		t.Errorf("Value = %d, want 4 (3 block + 1 arcane from banish). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_WeepingBattlegroundFizzlesWithoutAura: hand is just Weeping Battleground — no
// aura anywhere, so the banish rider fizzles. WB still blocks 3 of the 4 incoming. Value = 3.
func TestBest_WeepingBattlegroundFizzlesWithoutAura(t *testing.T) {
	h := []card.Card{runeblade.WeepingBattlegroundRed{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 3 {
		t.Errorf("Value = %d, want 3 (3 block only; banish fizzles). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}
