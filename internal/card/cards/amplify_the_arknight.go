// Amplify the Arknight — Runeblade Action - Attack. Printed cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 6, Yellow 5, Blue 4.
// Text: "Amplify the Arknight costs {r} less to play for each Runechant you control."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

const amplifyTheArknightPrintedCost = 3

// amplifyTheArknightEffectiveCost is max(0, printed - RunechantCount).
func amplifyTheArknightEffectiveCost(ge card.GameEngine) int {
	eff := amplifyTheArknightPrintedCost - ge.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (AmplifyTheArknightRed) EffectiveCost(ge card.GameEngine) int {
	return amplifyTheArknightEffectiveCost(ge)
}
func (AmplifyTheArknightRed) MinCost() int { return 0 }

func (AmplifyTheArknightRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (AmplifyTheArknightYellow) EffectiveCost(ge card.GameEngine) int {
	return amplifyTheArknightEffectiveCost(ge)
}
func (AmplifyTheArknightYellow) MinCost() int { return 0 }

func (AmplifyTheArknightYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (AmplifyTheArknightBlue) EffectiveCost(ge card.GameEngine) int {
	return amplifyTheArknightEffectiveCost(ge)
}
func (AmplifyTheArknightBlue) MinCost() int { return 0 }

func (AmplifyTheArknightBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
