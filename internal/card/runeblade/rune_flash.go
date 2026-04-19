// Rune Flash — Runeblade Action - Attack. Cost 3, Defense 3. Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Rune Flash costs {r} less to play for each Runechant you control."
//
// Cost() returns the printed cost; DiscountPerRunechant lets the chain apply
// max(0, Cost() - Runechants) at play time.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runeFlashTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

const runeFlashPrintedCost = 3

type RuneFlashRed struct{}

func (RuneFlashRed) ID() card.ID                 { return card.RuneFlashRed }
func (RuneFlashRed) Name() string               { return "Rune Flash (Red)" }
func (RuneFlashRed) Cost() int                  { return runeFlashPrintedCost }
func (RuneFlashRed) PrintedCost() int           { return runeFlashPrintedCost }
func (RuneFlashRed) Pitch() int                 { return 1 }
func (RuneFlashRed) Attack() int                { return 4 }
func (RuneFlashRed) Defense() int               { return 3 }
func (RuneFlashRed) Types() card.TypeSet        { return runeFlashTypes }
func (RuneFlashRed) GoAgain() bool              { return true }
func (c RuneFlashRed) Play(*card.TurnState) int { return c.Attack() }

type RuneFlashYellow struct{}

func (RuneFlashYellow) ID() card.ID                 { return card.RuneFlashYellow }
func (RuneFlashYellow) Name() string               { return "Rune Flash (Yellow)" }
func (RuneFlashYellow) Cost() int                  { return runeFlashPrintedCost }
func (RuneFlashYellow) PrintedCost() int           { return runeFlashPrintedCost }
func (RuneFlashYellow) Pitch() int                 { return 2 }
func (RuneFlashYellow) Attack() int                { return 3 }
func (RuneFlashYellow) Defense() int               { return 3 }
func (RuneFlashYellow) Types() card.TypeSet        { return runeFlashTypes }
func (RuneFlashYellow) GoAgain() bool              { return true }
func (c RuneFlashYellow) Play(*card.TurnState) int { return c.Attack() }

type RuneFlashBlue struct{}

func (RuneFlashBlue) ID() card.ID                 { return card.RuneFlashBlue }
func (RuneFlashBlue) Name() string               { return "Rune Flash (Blue)" }
func (RuneFlashBlue) Cost() int                  { return runeFlashPrintedCost }
func (RuneFlashBlue) PrintedCost() int           { return runeFlashPrintedCost }
func (RuneFlashBlue) Pitch() int                 { return 3 }
func (RuneFlashBlue) Attack() int                { return 2 }
func (RuneFlashBlue) Defense() int               { return 3 }
func (RuneFlashBlue) Types() card.TypeSet        { return runeFlashTypes }
func (RuneFlashBlue) GoAgain() bool              { return true }
func (c RuneFlashBlue) Play(*card.TurnState) int { return c.Attack() }
