// Reduce to Runechant — Runeblade Defense Reaction. Printed cost 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Reduce to Runechant costs {r} less to play for each Runechant you control. Create a
// Runechant token."
//
// Cost returns max(0, printed - s.Runechants()) at play time (card.VariableCost bounds [0, 1]).
// Play creates one Runechant, crediting +1 at creation. Defense-reaction state is reset
// between reactions so the token itself doesn't carry into next turn's carryover — only its
// damage credit lands.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

const reduceToRunechantPrintedCost = 1

func reduceToRunechantCost(s card.GameEngine) int {
	eff := reduceToRunechantPrintedCost - s.Runechants()
	if eff < 0 {
		return 0
	}
	return eff
}

func (ReduceToRunechantRed) Cost(s card.GameEngine) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantRed) MinCost() int               { return 0 }
func (ReduceToRunechantRed) MaxCost() int               { return reduceToRunechantPrintedCost }

func (ReduceToRunechantRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (ReduceToRunechantYellow) Cost(s card.GameEngine) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantYellow) MinCost() int               { return 0 }
func (ReduceToRunechantYellow) MaxCost() int               { return reduceToRunechantPrintedCost }

func (ReduceToRunechantYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (ReduceToRunechantBlue) Cost(s card.GameEngine) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantBlue) MinCost() int               { return 0 }
func (ReduceToRunechantBlue) MaxCost() int               { return reduceToRunechantPrintedCost }

func (ReduceToRunechantBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
