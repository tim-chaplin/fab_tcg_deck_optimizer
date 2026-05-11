// Amplify the Arknight — Runeblade Action - Attack. Printed cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 6, Yellow 5, Blue 4.
// Text: "Amplify the Arknight costs {r} less to play for each Runechant you control."
//
// Variable cost: Cost reads s.Runechants() to return max(0, printed - Runechants).
// Standard sim.VariableCost wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

const amplifyTheArknightPrintedCost = 3

func amplifyTheArknightCost(s *sim.TurnState) int {
	eff := amplifyTheArknightPrintedCost - s.Runechants()
	if eff < 0 {
		return 0
	}
	return eff
}

func (AmplifyTheArknightRed) Cost(s *sim.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightRed) MinCost() int              { return 0 }
func (AmplifyTheArknightRed) MaxCost() int              { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (AmplifyTheArknightYellow) Cost(s *sim.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightYellow) MinCost() int              { return 0 }
func (AmplifyTheArknightYellow) MaxCost() int              { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (AmplifyTheArknightBlue) Cost(s *sim.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightBlue) MinCost() int              { return 0 }
func (AmplifyTheArknightBlue) MaxCost() int              { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
