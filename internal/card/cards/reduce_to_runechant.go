// Reduce to Runechant — Runeblade Defense Reaction. Printed cost 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Reduce to Runechant costs {r} less to play for each Runechant you control. Create a
// Runechant token."
//
// Cost returns max(0, printed - RunechantCount). The created Runechant credits at creation;
// since defense-reaction state is reset between reactions, only its damage credit lands —
// the token itself doesn't carry forward.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

const reduceToRunechantPrintedCost = 1

func reduceToRunechantCost(ge card.GameEngine) int {
	eff := reduceToRunechantPrintedCost - ge.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (ReduceToRunechantRed) Cost(ge card.GameEngine) int { return reduceToRunechantCost(ge) }
func (ReduceToRunechantRed) MinCost() int                { return 0 }
func (ReduceToRunechantRed) MaxCost() int                { return reduceToRunechantPrintedCost }

func (ReduceToRunechantRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	createRunechantsAndLog(ge, l, self.Card.DisplayName(), 1)
}

func (ReduceToRunechantYellow) Cost(ge card.GameEngine) int { return reduceToRunechantCost(ge) }
func (ReduceToRunechantYellow) MinCost() int                { return 0 }
func (ReduceToRunechantYellow) MaxCost() int                { return reduceToRunechantPrintedCost }

func (ReduceToRunechantYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	createRunechantsAndLog(ge, l, self.Card.DisplayName(), 1)
}

func (ReduceToRunechantBlue) Cost(ge card.GameEngine) int { return reduceToRunechantCost(ge) }
func (ReduceToRunechantBlue) MinCost() int                { return 0 }
func (ReduceToRunechantBlue) MaxCost() int                { return reduceToRunechantPrintedCost }

func (ReduceToRunechantBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	createRunechantsAndLog(ge, l, self.Card.DisplayName(), 1)
}
