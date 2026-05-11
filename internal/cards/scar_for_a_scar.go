// Scar for a Scar — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ScarForAScarRed) GoAgain() bool { return sim.HeroWantsLowerHealth() }
func (c ScarForAScarRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
}

func (ScarForAScarYellow) GoAgain() bool { return sim.HeroWantsLowerHealth() }
func (c ScarForAScarYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
}

func (ScarForAScarBlue) GoAgain() bool { return sim.HeroWantsLowerHealth() }
func (c ScarForAScarBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
}
