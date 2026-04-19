// Drawn to the Dark Dimension — Runeblade Action - Attack. Printed cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Drawn to the Dark Dimension costs {r} less to play for each Runechant you control.
// Draw a card."
//
// Variable cost: Cost reads s.Runechants to return max(0, printed - Runechants) at play time.
// Implements card.VariableCost so the solver can pre-screen with MinCost / MaxCost bounds.
//
// Simplification: the "Draw a card" rider is modelled as a flat +3 damage — assume the drawn
// card goes to arsenal and is played on a future turn for roughly one card's worth of value.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var drawnToTheDarkDimensionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

const drawnToTheDarkDimensionPrintedCost = 2

func drawnToTheDarkDimensionCost(s *card.TurnState) int {
	eff := drawnToTheDarkDimensionPrintedCost - s.Runechants
	if eff < 0 {
		return 0
	}
	return eff
}

type DrawnToTheDarkDimensionRed struct{}

func (DrawnToTheDarkDimensionRed) ID() card.ID                 { return card.DrawnToTheDarkDimensionRed }
func (DrawnToTheDarkDimensionRed) Name() string                { return "Drawn to the Dark Dimension (Red)" }
func (DrawnToTheDarkDimensionRed) Cost(s *card.TurnState) int  { return drawnToTheDarkDimensionCost(s) }
func (DrawnToTheDarkDimensionRed) MinCost() int                { return 0 }
func (DrawnToTheDarkDimensionRed) MaxCost() int                { return drawnToTheDarkDimensionPrintedCost }
func (DrawnToTheDarkDimensionRed) Pitch() int                  { return 1 }
func (DrawnToTheDarkDimensionRed) Attack() int                 { return 3 }
func (DrawnToTheDarkDimensionRed) Defense() int                { return 3 }
func (DrawnToTheDarkDimensionRed) Types() card.TypeSet         { return drawnToTheDarkDimensionTypes }
func (DrawnToTheDarkDimensionRed) GoAgain() bool               { return false }
func (c DrawnToTheDarkDimensionRed) Play(*card.TurnState) int  { return c.Attack() + card.DrawValue }

type DrawnToTheDarkDimensionYellow struct{}

func (DrawnToTheDarkDimensionYellow) ID() card.ID                { return card.DrawnToTheDarkDimensionYellow }
func (DrawnToTheDarkDimensionYellow) Name() string               { return "Drawn to the Dark Dimension (Yellow)" }
func (DrawnToTheDarkDimensionYellow) Cost(s *card.TurnState) int { return drawnToTheDarkDimensionCost(s) }
func (DrawnToTheDarkDimensionYellow) MinCost() int               { return 0 }
func (DrawnToTheDarkDimensionYellow) MaxCost() int               { return drawnToTheDarkDimensionPrintedCost }
func (DrawnToTheDarkDimensionYellow) Pitch() int                 { return 2 }
func (DrawnToTheDarkDimensionYellow) Attack() int                { return 2 }
func (DrawnToTheDarkDimensionYellow) Defense() int               { return 3 }
func (DrawnToTheDarkDimensionYellow) Types() card.TypeSet        { return drawnToTheDarkDimensionTypes }
func (DrawnToTheDarkDimensionYellow) GoAgain() bool              { return false }
func (c DrawnToTheDarkDimensionYellow) Play(*card.TurnState) int { return c.Attack() + card.DrawValue }

type DrawnToTheDarkDimensionBlue struct{}

func (DrawnToTheDarkDimensionBlue) ID() card.ID                { return card.DrawnToTheDarkDimensionBlue }
func (DrawnToTheDarkDimensionBlue) Name() string               { return "Drawn to the Dark Dimension (Blue)" }
func (DrawnToTheDarkDimensionBlue) Cost(s *card.TurnState) int { return drawnToTheDarkDimensionCost(s) }
func (DrawnToTheDarkDimensionBlue) MinCost() int               { return 0 }
func (DrawnToTheDarkDimensionBlue) MaxCost() int               { return drawnToTheDarkDimensionPrintedCost }
func (DrawnToTheDarkDimensionBlue) Pitch() int                 { return 3 }
func (DrawnToTheDarkDimensionBlue) Attack() int                { return 1 }
func (DrawnToTheDarkDimensionBlue) Defense() int               { return 3 }
func (DrawnToTheDarkDimensionBlue) Types() card.TypeSet        { return drawnToTheDarkDimensionTypes }
func (DrawnToTheDarkDimensionBlue) GoAgain() bool              { return false }
func (c DrawnToTheDarkDimensionBlue) Play(*card.TurnState) int { return c.Attack() + card.DrawValue }
