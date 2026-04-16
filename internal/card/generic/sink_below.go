// Sink Below — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a card."
// Simplification: the hand-cycling rider is ignored (card quality/draw isn't modelled).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type SinkBelowRed struct{}

func (SinkBelowRed) Name() string             { return "Sink Below (Red)" }
func (SinkBelowRed) Cost() int                { return 0 }
func (SinkBelowRed) Pitch() int               { return 1 }
func (SinkBelowRed) Attack() int              { return 0 }
func (SinkBelowRed) Defense() int             { return 4 }
func (SinkBelowRed) Types() card.TypeSet      { return defenseReactionTypes }
func (SinkBelowRed) GoAgain() bool            { return false }
func (SinkBelowRed) Play(*card.TurnState) int { return 0 }

type SinkBelowYellow struct{}

func (SinkBelowYellow) Name() string             { return "Sink Below (Yellow)" }
func (SinkBelowYellow) Cost() int                { return 0 }
func (SinkBelowYellow) Pitch() int               { return 2 }
func (SinkBelowYellow) Attack() int              { return 0 }
func (SinkBelowYellow) Defense() int             { return 3 }
func (SinkBelowYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (SinkBelowYellow) GoAgain() bool            { return false }
func (SinkBelowYellow) Play(*card.TurnState) int { return 0 }

type SinkBelowBlue struct{}

func (SinkBelowBlue) Name() string             { return "Sink Below (Blue)" }
func (SinkBelowBlue) Cost() int                { return 0 }
func (SinkBelowBlue) Pitch() int               { return 3 }
func (SinkBelowBlue) Attack() int              { return 0 }
func (SinkBelowBlue) Defense() int             { return 2 }
func (SinkBelowBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (SinkBelowBlue) GoAgain() bool            { return false }
func (SinkBelowBlue) Play(*card.TurnState) int { return 0 }
