// Weeping Battleground — Runeblade Defense Reaction. Cost 0, Defense 3, Arcane 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "You may banish an aura from your graveyard. If you do, deal 1 arcane damage to target
// hero."
// Simplification: assume we always have an aura in the graveyard to banish, so the 1 arcane
// damage always triggers. Reported as Play()'s return so it counts toward dealt damage even if
// the printed Defense already covers all incoming damage.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var weepingBattlegroundTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

// weepingBattlegroundPlay deals the 1 arcane from the always-banished aura and marks
// ArcaneDamageDealt so later-this-turn triggers see the flag. Shared across all three printings.
func weepingBattlegroundPlay(s *card.TurnState) int {
	s.ArcaneDamageDealt = true
	return 1
}

type WeepingBattlegroundRed struct{}

func (WeepingBattlegroundRed) ID() card.ID                  { return card.WeepingBattlegroundRed }
func (WeepingBattlegroundRed) Name() string                 { return "Weeping Battleground (Red)" }
func (WeepingBattlegroundRed) Cost() int                    { return 0 }
func (WeepingBattlegroundRed) Pitch() int                   { return 1 }
func (WeepingBattlegroundRed) Attack() int                  { return 0 }
func (WeepingBattlegroundRed) Defense() int                 { return 3 }
func (WeepingBattlegroundRed) Types() card.TypeSet          { return weepingBattlegroundTypes }
func (WeepingBattlegroundRed) GoAgain() bool                { return false }
func (WeepingBattlegroundRed) Play(s *card.TurnState) int   { return weepingBattlegroundPlay(s) }

type WeepingBattlegroundYellow struct{}

func (WeepingBattlegroundYellow) ID() card.ID                  { return card.WeepingBattlegroundYellow }
func (WeepingBattlegroundYellow) Name() string                 { return "Weeping Battleground (Yellow)" }
func (WeepingBattlegroundYellow) Cost() int                    { return 0 }
func (WeepingBattlegroundYellow) Pitch() int                   { return 2 }
func (WeepingBattlegroundYellow) Attack() int                  { return 0 }
func (WeepingBattlegroundYellow) Defense() int                 { return 3 }
func (WeepingBattlegroundYellow) Types() card.TypeSet          { return weepingBattlegroundTypes }
func (WeepingBattlegroundYellow) GoAgain() bool                { return false }
func (WeepingBattlegroundYellow) Play(s *card.TurnState) int   { return weepingBattlegroundPlay(s) }

type WeepingBattlegroundBlue struct{}

func (WeepingBattlegroundBlue) ID() card.ID                  { return card.WeepingBattlegroundBlue }
func (WeepingBattlegroundBlue) Name() string                 { return "Weeping Battleground (Blue)" }
func (WeepingBattlegroundBlue) Cost() int                    { return 0 }
func (WeepingBattlegroundBlue) Pitch() int                   { return 3 }
func (WeepingBattlegroundBlue) Attack() int                  { return 0 }
func (WeepingBattlegroundBlue) Defense() int                 { return 3 }
func (WeepingBattlegroundBlue) Types() card.TypeSet          { return weepingBattlegroundTypes }
func (WeepingBattlegroundBlue) GoAgain() bool                { return false }
func (WeepingBattlegroundBlue) Play(s *card.TurnState) int   { return weepingBattlegroundPlay(s) }
