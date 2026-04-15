// Drawn to the Dark Dimension — Runeblade Action - Attack. Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Drawn to the Dark Dimension costs {r} less to play for each Runechant you control.
// Draw a card."
// Simplification: assume enough Runechants to fully discount, so the effective cost is 0. The
// "Draw a card" rider is modelled as a flat +3 damage — assume the drawn card goes to arsenal
// and is played on a future turn for roughly one card's worth of value.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var drawnToTheDarkDimensionTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type DrawnToTheDarkDimensionRed struct{}

func (DrawnToTheDarkDimensionRed) Name() string               { return "Drawn to the Dark Dimension (Red)" }
func (DrawnToTheDarkDimensionRed) Cost() int                  { return 0 }
func (DrawnToTheDarkDimensionRed) Pitch() int                 { return 1 }
func (DrawnToTheDarkDimensionRed) Attack() int                { return 3 }
func (DrawnToTheDarkDimensionRed) Defense() int               { return 3 }
func (DrawnToTheDarkDimensionRed) Types() map[string]bool     { return drawnToTheDarkDimensionTypes }
func (DrawnToTheDarkDimensionRed) GoAgain() bool              { return false }
func (c DrawnToTheDarkDimensionRed) Play(*card.TurnState) int { return c.Attack() + 3 }

type DrawnToTheDarkDimensionYellow struct{}

func (DrawnToTheDarkDimensionYellow) Name() string               { return "Drawn to the Dark Dimension (Yellow)" }
func (DrawnToTheDarkDimensionYellow) Cost() int                  { return 0 }
func (DrawnToTheDarkDimensionYellow) Pitch() int                 { return 2 }
func (DrawnToTheDarkDimensionYellow) Attack() int                { return 2 }
func (DrawnToTheDarkDimensionYellow) Defense() int               { return 3 }
func (DrawnToTheDarkDimensionYellow) Types() map[string]bool     { return drawnToTheDarkDimensionTypes }
func (DrawnToTheDarkDimensionYellow) GoAgain() bool              { return false }
func (c DrawnToTheDarkDimensionYellow) Play(*card.TurnState) int { return c.Attack() + 3 }

type DrawnToTheDarkDimensionBlue struct{}

func (DrawnToTheDarkDimensionBlue) Name() string               { return "Drawn to the Dark Dimension (Blue)" }
func (DrawnToTheDarkDimensionBlue) Cost() int                  { return 0 }
func (DrawnToTheDarkDimensionBlue) Pitch() int                 { return 3 }
func (DrawnToTheDarkDimensionBlue) Attack() int                { return 1 }
func (DrawnToTheDarkDimensionBlue) Defense() int               { return 3 }
func (DrawnToTheDarkDimensionBlue) Types() map[string]bool     { return drawnToTheDarkDimensionTypes }
func (DrawnToTheDarkDimensionBlue) GoAgain() bool              { return false }
func (c DrawnToTheDarkDimensionBlue) Play(*card.TurnState) int { return c.Attack() + 3 }
