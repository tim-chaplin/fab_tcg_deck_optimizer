// Amplify the Arknight — Runeblade Action - Attack. Printed cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 6, Yellow 5, Blue 4.
// Text: "Amplify the Arknight costs {r} less to play for each Runechant you control."
//
// Cost returns max(0, printed - RunechantCount).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

const amplifyTheArknightPrintedCost = 3

func amplifyTheArknightCost(ge card.GameEngine) int {
	eff := amplifyTheArknightPrintedCost - ge.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (AmplifyTheArknightRed) Cost(ge card.GameEngine) int { return amplifyTheArknightCost(ge) }
func (AmplifyTheArknightRed) MinCost() int                { return 0 }
func (AmplifyTheArknightRed) MaxCost() int                { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (AmplifyTheArknightYellow) Cost(ge card.GameEngine) int { return amplifyTheArknightCost(ge) }
func (AmplifyTheArknightYellow) MinCost() int                { return 0 }
func (AmplifyTheArknightYellow) MaxCost() int                { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (AmplifyTheArknightBlue) Cost(ge card.GameEngine) int { return amplifyTheArknightCost(ge) }
func (AmplifyTheArknightBlue) MinCost() int                { return 0 }
func (AmplifyTheArknightBlue) MaxCost() int                { return amplifyTheArknightPrintedCost }

func (AmplifyTheArknightBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
