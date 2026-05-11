// Reduce to Runechant — Runeblade Defense Reaction. Printed cost 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Reduce to Runechant costs {r} less to play for each Runechant you control. Create a
// Runechant token."
//
// Cost returns max(0, printed - s.Runechants()) at play time (sim.VariableCost bounds [0, 1]).
// Play creates one Runechant, crediting +1 at creation. Defense-reaction state is reset
// between reactions so the token itself doesn't carry into next turn's carryover — only its
// damage credit lands.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

const reduceToRunechantPrintedCost = 1

func reduceToRunechantCost(s sim.GameEngine) int {
	eff := reduceToRunechantPrintedCost - s.Runechants()
	if eff < 0 {
		return 0
	}
	return eff
}

func (ReduceToRunechantRed) Cost(s sim.GameEngine) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantRed) MinCost() int              { return 0 }
func (ReduceToRunechantRed) MaxCost() int              { return reduceToRunechantPrintedCost }

func (ReduceToRunechantRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (ReduceToRunechantYellow) Cost(s sim.GameEngine) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantYellow) MinCost() int              { return 0 }
func (ReduceToRunechantYellow) MaxCost() int              { return reduceToRunechantPrintedCost }

func (ReduceToRunechantYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (ReduceToRunechantBlue) Cost(s sim.GameEngine) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantBlue) MinCost() int              { return 0 }
func (ReduceToRunechantBlue) MaxCost() int              { return reduceToRunechantPrintedCost }

func (ReduceToRunechantBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
