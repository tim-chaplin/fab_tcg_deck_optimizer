// Amplify the Arknight — Runeblade Action - Attack. Printed cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 6, Yellow 5, Blue 4.
// Text: "Amplify the Arknight costs {r} less to play for each Runechant you control."
//
// Variable cost: Cost reads s.Runechants() to return max(0, printed - Runechants).
// Standard card.VariableCost wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

const amplifyTheArknightPrintedCost = 3

func amplifyTheArknightCost(s card.GameEngine) int {
	eff := amplifyTheArknightPrintedCost - s.Runechants()
	if eff < 0 {
		return 0
	}
	return eff
}

func (AmplifyTheArknightRed) Cost(s card.GameEngine) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightRed) MinCost() int               { return 0 }
func (AmplifyTheArknightRed) MaxCost() int               { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}

func (AmplifyTheArknightYellow) Cost(s card.GameEngine) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightYellow) MinCost() int               { return 0 }
func (AmplifyTheArknightYellow) MaxCost() int               { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}

func (AmplifyTheArknightBlue) Cost(s card.GameEngine) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightBlue) MinCost() int               { return 0 }
func (AmplifyTheArknightBlue) MaxCost() int               { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}
