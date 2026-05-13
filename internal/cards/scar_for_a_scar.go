// Scar for a Scar — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// scarForAScarGoAgain returns true when the active hero opts into LowerHealthWanter. nil-g
// reads as false (the printed default).
func scarForAScarGoAgain(ge card.GameEngine) bool {
	if ge == nil {
		return false
	}
	return ge.HeroWantsLowerHealth()
}

func (ScarForAScarRed) GoAgain(ge card.GameEngine) bool { return scarForAScarGoAgain(ge) }
func (c ScarForAScarRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (ScarForAScarYellow) GoAgain(ge card.GameEngine) bool { return scarForAScarGoAgain(ge) }
func (c ScarForAScarYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (ScarForAScarBlue) GoAgain(ge card.GameEngine) bool { return scarForAScarGoAgain(ge) }
func (c ScarForAScarBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
