// Runerager Swarm — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "If you've played or created an aura this turn, this gets go again."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (RuneragerSwarmRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	runeragerSwarmPlay(ge, l, self)
}

func (RuneragerSwarmYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	runeragerSwarmPlay(ge, l, self)
}

func (RuneragerSwarmBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	runeragerSwarmPlay(ge, l, self)
}
func runeragerSwarmPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.AuraCreated() {
		self.GrantedGoAgain = true
	}
}
