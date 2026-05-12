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
func scarForAScarGoAgain(g card.GameEngine) bool {
	if g == nil {
		return false
	}
	return g.HeroWantsLowerHealth()
}

func (ScarForAScarRed) GoAgain(g card.GameEngine) bool { return scarForAScarGoAgain(g) }
func (c ScarForAScarRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

func (ScarForAScarYellow) GoAgain(g card.GameEngine) bool { return scarForAScarGoAgain(g) }
func (c ScarForAScarYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

func (ScarForAScarBlue) GoAgain(g card.GameEngine) bool { return scarForAScarGoAgain(g) }
func (c ScarForAScarBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}
