// Reduce to Runechant — Runeblade Defense Reaction. Printed cost 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Reduce to Runechant costs {r} less to play for each Runechant you control. Create a
// Runechant token."
//
// The created Runechant credits at creation; since defense-reaction state is reset between
// reactions, only its damage credit lands — the token itself doesn't carry forward.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

const reduceToRunechantPrintedCost = 1

// reduceToRunechantEffectiveCost is max(0, printed - RunechantCount).
func reduceToRunechantEffectiveCost(ge card.GameEngine) int {
	eff := reduceToRunechantPrintedCost - ge.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (ReduceToRunechantRed) EffectiveCost(ge card.GameEngine) int {
	return reduceToRunechantEffectiveCost(ge)
}
func (ReduceToRunechantRed) MinCost() int { return 0 }

func (ReduceToRunechantRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (ReduceToRunechantYellow) EffectiveCost(ge card.GameEngine) int {
	return reduceToRunechantEffectiveCost(ge)
}
func (ReduceToRunechantYellow) MinCost() int { return 0 }

func (ReduceToRunechantYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}

func (ReduceToRunechantBlue) EffectiveCost(ge card.GameEngine) int {
	return reduceToRunechantEffectiveCost(ge)
}
func (ReduceToRunechantBlue) MinCost() int { return 0 }

func (ReduceToRunechantBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
