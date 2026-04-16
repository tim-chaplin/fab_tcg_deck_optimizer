// Vexing Malice — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 2.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 2 arcane damage to target hero."
//
// Simplification: arcane damage counts as regular damage. Play returns power + 2.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var vexingMaliceTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type VexingMaliceRed struct{}

func (VexingMaliceRed) Name() string               { return "Vexing Malice (Red)" }
func (VexingMaliceRed) Cost() int                  { return 1 }
func (VexingMaliceRed) Pitch() int                 { return 1 }
func (VexingMaliceRed) Attack() int                { return 3 }
func (VexingMaliceRed) Defense() int               { return 3 }
func (VexingMaliceRed) Types() card.TypeSet        { return vexingMaliceTypes }
func (VexingMaliceRed) GoAgain() bool              { return false }
func (c VexingMaliceRed) Play(*card.TurnState) int { return c.Attack() + 2 }

type VexingMaliceYellow struct{}

func (VexingMaliceYellow) Name() string               { return "Vexing Malice (Yellow)" }
func (VexingMaliceYellow) Cost() int                  { return 1 }
func (VexingMaliceYellow) Pitch() int                 { return 2 }
func (VexingMaliceYellow) Attack() int                { return 2 }
func (VexingMaliceYellow) Defense() int               { return 3 }
func (VexingMaliceYellow) Types() card.TypeSet        { return vexingMaliceTypes }
func (VexingMaliceYellow) GoAgain() bool              { return false }
func (c VexingMaliceYellow) Play(*card.TurnState) int { return c.Attack() + 2 }

type VexingMaliceBlue struct{}

func (VexingMaliceBlue) Name() string               { return "Vexing Malice (Blue)" }
func (VexingMaliceBlue) Cost() int                  { return 1 }
func (VexingMaliceBlue) Pitch() int                 { return 3 }
func (VexingMaliceBlue) Attack() int                { return 1 }
func (VexingMaliceBlue) Defense() int               { return 3 }
func (VexingMaliceBlue) Types() card.TypeSet        { return vexingMaliceTypes }
func (VexingMaliceBlue) GoAgain() bool              { return false }
func (c VexingMaliceBlue) Play(*card.TurnState) int { return c.Attack() + 2 }
