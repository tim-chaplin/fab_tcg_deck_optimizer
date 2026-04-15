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

var weepingBattlegroundTypes = map[string]bool{"Runeblade": true, "Defense Reaction": true}

type WeepingBattlegroundRed struct{}

func (WeepingBattlegroundRed) Name() string             { return "Weeping Battleground (Red)" }
func (WeepingBattlegroundRed) Cost() int                { return 0 }
func (WeepingBattlegroundRed) Pitch() int               { return 1 }
func (WeepingBattlegroundRed) Attack() int              { return 0 }
func (WeepingBattlegroundRed) Defense() int             { return 3 }
func (WeepingBattlegroundRed) Types() map[string]bool   { return weepingBattlegroundTypes }
func (WeepingBattlegroundRed) GoAgain() bool            { return false }
func (WeepingBattlegroundRed) Play(*card.TurnState) int { return 1 }

type WeepingBattlegroundYellow struct{}

func (WeepingBattlegroundYellow) Name() string             { return "Weeping Battleground (Yellow)" }
func (WeepingBattlegroundYellow) Cost() int                { return 0 }
func (WeepingBattlegroundYellow) Pitch() int               { return 2 }
func (WeepingBattlegroundYellow) Attack() int              { return 0 }
func (WeepingBattlegroundYellow) Defense() int             { return 3 }
func (WeepingBattlegroundYellow) Types() map[string]bool   { return weepingBattlegroundTypes }
func (WeepingBattlegroundYellow) GoAgain() bool            { return false }
func (WeepingBattlegroundYellow) Play(*card.TurnState) int { return 1 }

type WeepingBattlegroundBlue struct{}

func (WeepingBattlegroundBlue) Name() string             { return "Weeping Battleground (Blue)" }
func (WeepingBattlegroundBlue) Cost() int                { return 0 }
func (WeepingBattlegroundBlue) Pitch() int               { return 3 }
func (WeepingBattlegroundBlue) Attack() int              { return 0 }
func (WeepingBattlegroundBlue) Defense() int             { return 3 }
func (WeepingBattlegroundBlue) Types() map[string]bool   { return weepingBattlegroundTypes }
func (WeepingBattlegroundBlue) GoAgain() bool            { return false }
func (WeepingBattlegroundBlue) Play(*card.TurnState) int { return 1 }
